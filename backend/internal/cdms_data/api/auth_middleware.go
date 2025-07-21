package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/jackc/pgx/v5"
	"github.com/jjckrbbt/cdms/backend/internal/db"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	Email string `json:"https://cdms.jjckrbbt.dev/email"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

type AuthMiddleware struct {
	auth0Domain    string
	auth0Audience  string
	queries        db.Querier
	logger         *slog.Logger
	tokenValidator *validator.Validator
}

func NewAuthMiddleware(domain string, audience string, q db.Querier, log *slog.Logger) (*AuthMiddleware, error) {

	issuerURL, err := url.Parse("https://" + domain + "/")
	if err != nil {
		return nil, err
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &CustomClaims{}
		}),
	)
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{
		auth0Domain:    domain,
		auth0Audience:  audience,
		queries:        q,
		logger:         log.With("component", "auth_middleware"),
		tokenValidator: jwtValidator,
	}, nil
}

func (m *AuthMiddleware) ValidateRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization Header fomat")
		}
		tokenString := strings.TrimSpace(parts[1])

		claims, err := m.tokenValidator.ValidateToken(ctx, tokenString)
		if err != nil {
			m.logger.WarnContext(ctx, "JWT validation failed", "error", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		validatedClaims, ok := claims.(*validator.ValidatedClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Could not parse validated claims")
		}

		authProviderSubject := validatedClaims.RegisteredClaims.Subject
		if authProviderSubject == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Subject (sub) claim not found in token")
		}

		m.logger.InfoContext(ctx, "Attempting to find user with subject ", "subject", authProviderSubject)

		//--JIT Provisioning ----
		user, err := m.queries.GetUserByAuthProviderSubject(ctx, authProviderSubject)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				m.logger.ErrorContext(ctx, "Database error looking up user by subject", "error", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
			}

			m.logger.InfoContext(ctx, "New user detected, provisioning account", "subject", authProviderSubject)

			customClaims, ok := validatedClaims.CustomClaims.(*CustomClaims)
			if !ok || customClaims.Email == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Could not parse custom claims for new user")
			}

			user, err = m.queries.CreateUserFromAuthProvider(ctx, db.CreateUserFromAuthProviderParams{
				AuthProviderSubject: authProviderSubject,
				Email:               customClaims.Email,
				FirstName:           "New",
				LastName:            "User",
				Org:                 db.UserOrgGSA,
			})
			if err != nil {
				m.logger.ErrorContext(ctx, "Failed to create new user in database during JIT provisioning", "error", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to provision new user")
			}
		}
		//--- END JIT Provisioning ---

		c.Set("user", user)

		return next(c)
	}
}
