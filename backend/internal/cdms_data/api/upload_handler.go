// internal/cdms_data/api/upload_handler.go
package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/jjckrbbt/cdms/backend/internal/cdms_data/importer"
	"github.com/jjckrbbt/cdms/backend/internal/cdms_data/processor"
)

// UploadHandler handles file uploads for various report types.
type UploadHandler struct {
	importer  *importer.Importer
	processor *processor.Processor
	logger    *slog.Logger
}

// NewUploadHandler creates and returns a new UploadHandler.
func NewUploadHandler(imp *importer.Importer, proc *processor.Processor, appLogger *slog.Logger) *UploadHandler {
	return &UploadHandler{
		importer:  imp,
		processor: proc,
		logger:    appLogger.With("component", "cdms_api_handler"),
	}
}

// AllowedReportTypes defines the valid report types for uploads.
var AllowedReportTypes = map[string]bool{
	"BC1300":            true,
	"BC1048":            true,
	"OUTSTANDING_BILLS": true,
	"VENDOR_CODE":       true,
}

// HandleUpload orchestrates the entire asynchronous file upload and processing workflow.
func (h *UploadHandler) HandleUpload(c echo.Context) error {
	requestID, _ := c.Get("requestID").(string) // Assuming requestID is from middleware
	if requestID == "" {
		requestID = uuid.New().String()
	}
	reqLogger := h.logger.With(
		"request_id", requestID,
		"http_method", c.Request().Method,
		"path", c.Request().URL.Path,
	)

	reportType := strings.ToUpper(c.Param("reportType"))
	if !AllowedReportTypes[reportType] {
		reqLogger.WarnContext(c.Request().Context(), "Attempted upload with unsupported report type", "report_type", reportType)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unsupported report type: '%s'", reportType))
	}

	reqLogger.InfoContext(c.Request().Context(), "Received file upload request", "report_type", reportType)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	fileHeader, err := c.FormFile("report_file")
	if err != nil {
		reqLogger.ErrorContext(ctx, "Failed to get file from form", "error", err, "form_field_name", "report_file")
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request: No file uploaded or wrong field name ('report_file')")
	}

	reqLogger.InfoContext(ctx, "File received from client", "filename", fileHeader.Filename, "size_bytes", fileHeader.Size)

	// Step 1: Store the file and create the initial "UPLOADED" record in the database.
	uploadRecord, err := h.importer.StoreFile(ctx, fileHeader, reportType)
	if err != nil {
		reqLogger.ErrorContext(ctx, "Failed to store uploaded file via Importer service", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to accept upload: %s", err.Error()))
	}

	// Step 2: Immediately return a 202 Accepted response to the client.
	response := map[string]string{
		"message":     fmt.Sprintf("File '%s' accepted for processing.", uploadRecord.Filename),
		"upload_id":   uploadRecord.ID.String(), // UPDATED: Use .String() to convert UUID to string for JSON
		"report_type": reportType,
	}
	if err := c.JSON(http.StatusAccepted, response); err != nil {
		reqLogger.ErrorContext(ctx, "Failed to write 202 JSON response", "error", err)
		return err
	}
	reqLogger.InfoContext(ctx, "File accepted and 202 response sent to client", "upload_id", uploadRecord.ID.String())

	// Step 3: Trigger the asynchronous processing in a background goroutine.
	go func() {
		procCtx, procCancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer procCancel()

		// The processor does its work and returns a detailed result.
		result := h.processor.ProcessFileFromCloudStorage(procCtx, uploadRecord.ID.String(), uploadRecord.StorageKey, uploadRecord.ReportType)

		// REVISED LOGIC: Use the result to update the database record and log the final outcome.
		var errorMsg string
		if result.Error != nil {
			errorMsg = result.Error.Error()
			reqLogger.ErrorContext(procCtx, "Asynchronous report processing failed",
				"upload_id", uploadRecord.ID.String(),
				"final_status", result.Status,
				"error", result.Error)
		} else {
			reqLogger.InfoContext(procCtx, "Asynchronous report processing completed",
				"upload_id", uploadRecord.ID.String(),
				"final_status", result.Status,
				"rows_upserted", result.RowsUpserted,
				"rows_removed", result.RowsRemoved)
		}

		// Persist the final status to the database.
		if err := h.importer.UpdateUploadStatus(procCtx, uploadRecord.ID, result.Status, errorMsg, result.RowsUpserted, result.RowsRemoved); err != nil {
			// This is a critical failure, as the system could not record the outcome.
			reqLogger.ErrorContext(procCtx, "FATAL: Could not update final upload status in DB", "update_error", err)
		}
	}()

	return nil // Return nil because the response has already been sent.
}
