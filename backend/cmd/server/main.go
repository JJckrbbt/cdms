// cmd/api/main.go
package main

import (
	//"context" // Required for context.Background() in Echo's logger setup
	"fmt"
	"io"
	"log/slog" // For slog types
	"net/http" // Still needed for http.StatusX, etc.
	"os"       // For os.Exit, os.Stderr
	"time"

	"github.com/jjckrbbt/cdms/backend/internal/cdms_data/api"
	"github.com/jjckrbbt/cdms/backend/internal/cdms_data/importer"
	"github.com/jjckrbbt/cdms/backend/internal/cdms_data/processor"
	"github.com/jjckrbbt/cdms/backend/internal/config"
	"github.com/jjckrbbt/cdms/backend/internal/database"
	"github.com/jjckrbbt/cdms/backend/internal/db"
	"github.com/jjckrbbt/cdms/backend/internal/logger"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Load application configuration FIRST.
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize Sentry and then Sentry's handler
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.AppEnv,
		TracesSampleRate: 1.0,
		Debug:            true,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	defer sentry.Flush(2 * time.Second)

	// 3. Initialize the Logger.
	logger.InitLogger(cfg.AppEnv)
	appLogger := logger.L() // Get the configured logger instance

	appLogger.Info("Application starting up...", "environment", cfg.AppEnv)

	// 4. Connect to the Database.
	dbClient, err := database.ConnectDB(cfg.DatabaseURL, appLogger.With("component", "database_connector"))
	if err != nil {
		appLogger.Error("Failed to connect to database at startup", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := dbClient.Close(); err != nil {
			appLogger.Error("Error closing database connection", slog.Any("error", err))
		}
	}()
	appLogger.Info("Database connection established.")

	// 5. Initialize Core Application Components.
	realQuerier := db.New(dbClient.Pool)

	importerLogger := appLogger.With("service", "file_importer")
	fileImporter, err := importer.NewImporter(dbClient, importerLogger, cfg)
	if err != nil {
		appLogger.Error("Failed to initialize file importer service", slog.Any("error", err))
		os.Exit(1)
	}
	appLogger.Info("File Importer service initialized.")

	processorLogger := appLogger.With("service", "cdms_data_processor")
	cdmsProcessor, err := processor.NewProcessor(dbClient, processorLogger, cfg)
	if err != nil {
		appLogger.Error("Failed to initialize CDMS data processor service", slog.Any("error", err))
		os.Exit(1)
	}
	appLogger.Info("CDMS Data Processor service initialized.")

	// Initialize your HTTP API handlers.
	apiLogger := appLogger.With("service", "api_handlers")

	uploadHandler := api.NewUploadHandler(fileImporter, cdmsProcessor, realQuerier, apiLogger)
	chargebackHandler := api.NewChargebackHandler(realQuerier, apiLogger)
	delinquencyHandler := api.NewDelinquencyHandler(realQuerier, apiLogger)
	dashboardHandler := api.NewDashboardHandler(realQuerier, apiLogger)
	userHandler := api.NewUserHandler(realQuerier, apiLogger)

	appLogger.Info("API handlers initialized.")

	// 6. Initialize Echo.
	e := echo.New()

	// Configure Echo's logger to use our slog instance.
	e.Logger.SetOutput(io.Discard)
	e.Logger.SetLevel(0)   // Set to 0 to disable logging, we use slog
	e.Logger.SetHeader("") // Remove default header, slog adds better ones

	// 7. Register Middleware.
	// Recover middleware: Recovers from panics anywhere in the chain and handles the error.
	e.Use(middleware.Recover())
	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://34.8.206.198", "https://cdms-backend-414620627769.us-central1.run.app", "https://cdms.jjckrbbt.dev"}, // Replace with your React dev server URL
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Accept", "Authorization"},
		// Add AllowCredentials: true if you send cookies/credentials
	}))

	// --- Auth Middleware Setup ---
	authMiddleware, err := api.NewAuthMiddleware(cfg.Auth0Domain, cfg.Auth0Audience, realQuerier, appLogger)
	if err != nil {
		appLogger.Error("Failed to initialize auth middleware", slog.Any("error", err))
		os.Exit(1)
	}

	apiGroup := e.Group("/api")
	apiGroup.Use(authMiddleware.ValidateRequest)
	//-------------------------------
	// Request Logger Middleware (For consistent request logging)
	// This logs basic request info using our slog instance.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := uuid.New().String() // Generate/extract request ID
			c.Set("requestID", reqID)    // Store request ID in context for later access

			start := time.Now()

			if hub := sentryecho.GetHubFromContext(c); hub != nil {
				hub.Scope().SetTag("request_id", c.Get("requestID").(string))
			}

			err := next(c)
			stop := time.Now()

			status := c.Response().Status
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				}
			}

			// Log the request summary with context
			appLogger.InfoContext(c.Request().Context(), "HTTP Request",
				"request_id", reqID,
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", status,
				"latency_ms", stop.Sub(start).Milliseconds(),
				"user_agent", c.Request().UserAgent(),
				"ip", c.RealIP(),
			)
			return err
		}
	})

	e.Use(sentryecho.New(sentryecho.Options{
		Repanic: true,
	}))

	// 8. Register Routes.

	// Health check endpoint (simple GET)
	e.GET("/health", func(c echo.Context) error {
		// Log using a logger with request context
		reqLogger := appLogger.With("request_id", c.Get("requestID")) // Retrieve request ID from context
		reqLogger.InfoContext(c.Request().Context(), "Health check requested", "ip", c.RealIP())

		if err := dbClient.Ping(); err != nil {
			reqLogger.ErrorContext(c.Request().Context(), "Database ping failed during health check", slog.Any("error", err))

			sentry.CaptureException(err)

			return c.String(http.StatusInternalServerError, "DB Not Ready") // Return string response for error
		}
		return c.String(http.StatusOK, "OK") // Return string response for success
	})

	//Upload group
	apiGroup.POST("/upload/:reportType", uploadHandler.HandleUpload)

	//User Routes
	userRoutes := apiGroup.Group("/users")
	userRoutes.Use(userHandler.LoadUserContextMiddleware)
	apiGroup.GET("/me", userHandler.HandleGetMe)
	//Admin routes for user management
	adminUserRoutes := userRoutes.Group("/admin")
	adminUserRoutes.GET("", userHandler.HandleListUsers, api.RequirePermission("users:view_scoped"))
	adminUserRoutes.PATCH("/:id", userHandler.HandleUpdateUser, api.RequirePermission("users:edit"))
	adminUserRoutes.PATCH("/id/roles", userHandler.HandleUpdateUserRoles, api.RequirePermission("roles:assign_global"))
	adminUserRoutes.PATCH("/:id/business-lines", userHandler.HandleUpdateUserBusinessLines, api.RequirePermission("roles:assign_scoped"))
	//Upload Reporting Group
	uploadRoutes := apiGroup.Group("/uploads")
	uploadRoutes.GET("", uploadHandler.HandleGetUploads)
	uploadRoutes.GET("/removed_rows/:id", uploadHandler.HandleGetRemovedRows)

	//Chargeback group
	chargebackRoutes := apiGroup.Group("/chargebacks")
	chargebackRoutes.GET("", chargebackHandler.HandleGetChargebacks)
	chargebackRoutes.GET("/:id", chargebackHandler.HandleGetByID)
	chargebackRoutes.GET("/history/:id", chargebackHandler.HandleChargebackStatus)
	chargebackRoutes.POST("", chargebackHandler.HandleCreate)
	chargebackRoutes.PATCH("/:id", chargebackHandler.HandleUpdate)

	//Delinquency group
	delinquencyRoutes := apiGroup.Group("/delinquencies")
	delinquencyRoutes.GET("", delinquencyHandler.HandleGetDelinquencies)
	delinquencyRoutes.GET("/:id", delinquencyHandler.HandleGetByID)
	delinquencyRoutes.GET("/history/:id", delinquencyHandler.HandleDelinquencyStatus)
	delinquencyRoutes.POST("", delinquencyHandler.HandleCreate)
	delinquencyRoutes.PATCH("/:id", delinquencyHandler.HandleUpdate)

	//Dashbord group
	apiGroup.GET("/dashboard", dashboardHandler.HandleGetDashboardStats)

	e.GET("/foo", func(ctx echo.Context) error {
		// sentryecho handler will catch it just fine. Also, because we attached "someRandomTag"
		// in the middleware before, it will be sent through as well
		panic("y tho")
	})

	//for _, route := range e.Routes() {
	//	appLogger.Info("Registered Route", "method", route.Method, "path", route.Path)
	//}

	// 9. Start the HTTP server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf(":%s", port)

	appLogger.Info("HTTP Server starting on port", "port", port)

	// e.Start blocks until the server is shut down or an error occurs.
	if err := e.Start(address); err != nil && err != http.ErrServerClosed {
		// Only log fatal if it's not a graceful shutdown error.
		appLogger.Error("HTTP Server failed to start", slog.Any("error", err))
		os.Exit(1)
	}
	// This message would appear after a graceful shutdown.
	appLogger.Info("HTTP Server stopped gracefully.")
}
