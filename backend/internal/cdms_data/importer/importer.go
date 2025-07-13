package importer

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jjckrbbt/cdms/backend/internal/cdms_data/model"
	"github.com/jjckrbbt/cdms/backend/internal/config"
	"github.com/jjckrbbt/cdms/backend/internal/database"
	"github.com/jjckrbbt/cdms/backend/internal/db"
	"google.golang.org/api/option"
)

type Importer struct {
	dbClient  *database.DBClient
	logger    *slog.Logger
	cfg       *config.Config
	gcsClient *storage.Client
	gcsBucket string
}

func NewImporter(dbClient *database.DBClient, logger *slog.Logger, cfg *config.Config) (*Importer, error) {
	ctx := context.Background()

	var clientOptions []option.ClientOption
	if keyFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); keyFilePath != "" {
		clientOptions = append(clientOptions, option.WithCredentialsFile(keyFilePath))
		logger.Info("Importer: Using GOOGLE_APPLICATION_CREDENTIALS for GCS client initialization")
	} else {
		logger.Info("Importer: GOOGLE_APPLICATION_CREDENTIALS not set, relying on default GCP authentication.")
	}

	gcsClient, err := storage.NewClient(ctx, clientOptions...)
	if err != nil {
		logger.Error("Importer: Failed to create GCS client", "error", err)
		return nil, fmt.Errorf("failed to create GCS client for importer: %w", err)
	}

	return &Importer{
		dbClient:  dbClient,
		logger:    logger.With("component", "file_importer"),
		cfg:       cfg,
		gcsClient: gcsClient,
		gcsBucket: cfg.GCSBucketName,
	}, nil
}

func (i *Importer) StoreFile(ctx context.Context, fileHeader *multipart.FileHeader, reportType string) (*model.Upload, error) {
	uploadID := uuid.New()
	gcsObjectKey := fmt.Sprintf("raw-reports/%s/%s-%s", reportType, uploadID.String(), fileHeader.Filename)

	file, err := fileHeader.Open()
	if err != nil {
		i.logger.ErrorContext(ctx, "Importer: Failed to open uploaded file", "error", err, "upload_id", uploadID)
		return nil, fmt.Errorf("importer: failed to open uploaded file: %w", err)
	}
	defer file.Close()

	wc := i.gcsClient.Bucket(i.gcsBucket).Object(gcsObjectKey).NewWriter(ctx)
	wc.ContentType = fileHeader.Header.Get("Content-Type")

	if _, err := io.Copy(wc, file); err != nil {
		wc.Close()
		i.logger.ErrorContext(ctx, "Importer: Failed to write file content to GCS", "error", err, "object_key", gcsObjectKey)
		return nil, fmt.Errorf("importer: failed to write file content to GCS: %w", err)
	}
	if err := wc.Close(); err != nil {
		i.logger.ErrorContext(ctx, "Importer: Failed to close GCS writer", "error", err, "object_key", gcsObjectKey)
		return nil, fmt.Errorf("importer: failed to close GCS writer: %w", err)
	}

	queries := db.New(i.dbClient.Pool)

	params := db.CreateUploadParams{
		ID:                pgtype.UUID{Bytes: uploadID, Valid: true},
		StorageKey:        gcsObjectKey,
		Filename:          fileHeader.Filename,
		ReportType:        reportType,
		Status:            "UPLOADED",
		ProcessedByUserID: 1, // Assuming a default user ID for now; this should be replaced with actual user ID logic.,
	}

	createdUpload, err := queries.CreateUpload(ctx, params)
	if err != nil {
		i.logger.ErrorContext(ctx, "Importer: Failed to record upload in database", "error", err, "upload_id", uploadID)
		return nil, fmt.Errorf("importer: failed to record upload in DB: %w", err)
	}

	i.logger.InfoContext(ctx, "Importer: File stored in GCS and upload record created", "upload_id", uploadID, "gcs_key", gcsObjectKey)

	return &model.Upload{
		ID:         createdUpload.ID.Bytes,
		StorageKey: createdUpload.StorageKey,
		Filename:   createdUpload.Filename,
		ReportType: createdUpload.ReportType,
		Status:     createdUpload.Status,
		UploadedAt: createdUpload.UploadedAt.Time,
	}, nil
}

func (i *Importer) UpdateUploadStatus(ctx context.Context, uploadID uuid.UUID, status string, errorDetails string, rowsUpserted int64, rowsRemoved int) error {
	queries := db.New(i.dbClient.Pool)

	// sqlc will generate these new fields on the params struct.
	params := db.UpdateUploadStatusParams{
		ID:           pgtype.UUID{Bytes: uploadID, Valid: true},
		Status:       status,
		ErrorDetails: pgtype.Text{String: errorDetails, Valid: errorDetails != ""},
		RowsUpserted: pgtype.Int4{Int32: int32(rowsUpserted), Valid: true},
		RowsRemoved:  pgtype.Int4{Int32: int32(rowsRemoved), Valid: true},
	}

	err := queries.UpdateUploadStatus(ctx, params)
	if err != nil {
		i.logger.ErrorContext(ctx, "Importer: Failed to update upload status in database",
			"error", err,
			"upload_id", uploadID,
			"new_status", status,
		)
		return err
	}

	return nil
}
