package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jjckrbbt/cdms/backend/internal/cdms_data/importer"
	"github.com/jjckrbbt/cdms/backend/internal/cdms_data/processor"
	"github.com/jjckrbbt/cdms/backend/internal/db"
)

type UploadHandler struct {
	importer  *importer.Importer
	processor *processor.Processor
	queries   db.Querier
	logger    *slog.Logger
}

func NewUploadHandler(imp *importer.Importer, proc *processor.Processor, q db.Querier, appLogger *slog.Logger) *UploadHandler {
	return &UploadHandler{
		importer:  imp,
		processor: proc,
		queries:   q,
		logger:    appLogger.With("component", "cdms_api_handler"),
	}
}

var AllowedReportTypes = map[string]bool{
	"BC1300":            true,
	"BC1048":            true,
	"OUTSTANDING_BILLS": true,
	"VENDOR_CODE":       true,
}

func (h *UploadHandler) HandleGetUploads(c echo.Context) error {
	ctx := c.Request().Context()
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 25
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	params := db.ListUploadsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	uploads, err := h.queries.ListUploads(ctx, params)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to list uploads", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve upload history")
	}

	return c.JSON(http.StatusOK, uploads)
}

func (h *UploadHandler) HandleGetRemovedRows(c echo.Context) error {
	ctx := c.Request().Context()
	uploadIDStr := c.Param("id")
	uploadID, err := uuid.Parse(uploadIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid upload ID format")
	}

	pgUUID := pgtype.UUID{Bytes: uploadID, Valid: true}

	rows, err := h.queries.GetRemovedRowsByUploadID(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusOK, []db.RemovedRowsLog{})
		}
		h.logger.ErrorContext(ctx, "Failed to get removed rows for upload", uploadIDStr, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve removed rows")
	}

	return c.JSON(http.StatusOK, rows)
}

func (h *UploadHandler) HandleUpload(c echo.Context) error {
	requestID, _ := c.Get("requestID").(string)
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

	uploadRecord, err := h.importer.StoreFile(ctx, fileHeader, reportType)
	if err != nil {
		reqLogger.ErrorContext(ctx, "Failed to store uploaded file via Importer service", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to accept upload: %s", err.Error()))
	}

	response := map[string]string{
		"message":     fmt.Sprintf("File '%s' accepted for processing.", uploadRecord.Filename),
		"upload_id":   uploadRecord.ID.String(),
		"report_type": reportType,
	}
	if err := c.JSON(http.StatusAccepted, response); err != nil {
		reqLogger.ErrorContext(ctx, "Failed to write 202 JSON response", "error", err)
		return err
	}
	reqLogger.InfoContext(ctx, "File accepted and 202 response sent to client", "upload_id", uploadRecord.ID.String())

	go func() {
		procCtx, procCancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer procCancel()

		result := h.processor.ProcessFileFromCloudStorage(procCtx, uploadRecord.ID.String(), uploadRecord.StorageKey, uploadRecord.ReportType)

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

		if err := h.importer.UpdateUploadStatus(procCtx, uploadRecord.ID, result.Status, errorMsg, result.RowsUpserted, result.RowsRemoved); err != nil {
			reqLogger.ErrorContext(procCtx, "FATAL: Could not update final upload status in DB", "update_error", err)
		}
	}()

	return nil
}
