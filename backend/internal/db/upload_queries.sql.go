// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: upload_queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUpload = `-- name: CreateUpload :one
INSERT INTO uploads (
    id,
    storage_key,
    filename,
    report_type,
    status,
    processed_by_user_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, storage_key, filename, report_type, status, uploaded_at, processed_at, error_details, processed_by_user_id, rows_upserted, rows_removed
`

type CreateUploadParams struct {
	ID                pgtype.UUID `json:"id"`
	StorageKey        string      `json:"storage_key"`
	Filename          string      `json:"filename"`
	ReportType        string      `json:"report_type"`
	Status            string      `json:"status"`
	ProcessedByUserID pgtype.UUID `json:"processed_by_user_id"`
}

// Create a record to track a new file upload
func (q *Queries) CreateUpload(ctx context.Context, arg CreateUploadParams) (Upload, error) {
	row := q.db.QueryRow(ctx, createUpload,
		arg.ID,
		arg.StorageKey,
		arg.Filename,
		arg.ReportType,
		arg.Status,
		arg.ProcessedByUserID,
	)
	var i Upload
	err := row.Scan(
		&i.ID,
		&i.StorageKey,
		&i.Filename,
		&i.ReportType,
		&i.Status,
		&i.UploadedAt,
		&i.ProcessedAt,
		&i.ErrorDetails,
		&i.ProcessedByUserID,
		&i.RowsUpserted,
		&i.RowsRemoved,
	)
	return i, err
}

const getUpload = `-- name: GetUpload :one
SELECT id, storage_key, filename, report_type, status, uploaded_at, processed_at, error_details, processed_by_user_id, rows_upserted, rows_removed FROM uploads
WHERE id = $1
`

// Retrieve a detailed summary for a specific upload
func (q *Queries) GetUpload(ctx context.Context, id pgtype.UUID) (Upload, error) {
	row := q.db.QueryRow(ctx, getUpload, id)
	var i Upload
	err := row.Scan(
		&i.ID,
		&i.StorageKey,
		&i.Filename,
		&i.ReportType,
		&i.Status,
		&i.UploadedAt,
		&i.ProcessedAt,
		&i.ErrorDetails,
		&i.ProcessedByUserID,
		&i.RowsUpserted,
		&i.RowsRemoved,
	)
	return i, err
}

const listUploads = `-- name: ListUploads :many
SELECT id, storage_key, filename, report_type, status, uploaded_at, processed_at, error_details, processed_by_user_id, rows_upserted, rows_removed FROM uploads
ORDER BY uploaded_at DESC
LIMIT $1
OFFSET $2
`

type ListUploadsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// Provides a paginated list of recent report uploads and their statuses
func (q *Queries) ListUploads(ctx context.Context, arg ListUploadsParams) ([]Upload, error) {
	rows, err := q.db.Query(ctx, listUploads, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Upload
	for rows.Next() {
		var i Upload
		if err := rows.Scan(
			&i.ID,
			&i.StorageKey,
			&i.Filename,
			&i.ReportType,
			&i.Status,
			&i.UploadedAt,
			&i.ProcessedAt,
			&i.ErrorDetails,
			&i.ProcessedByUserID,
			&i.RowsUpserted,
			&i.RowsRemoved,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUploadStatus = `-- name: UpdateUploadStatus :exec
UPDATE uploads
SET
    status = $2,
    error_details = $3,
    rows_upserted = $4,
    rows_removed = $5,
    processed_at = NOW()
WHERE id = $1
`

type UpdateUploadStatusParams struct {
	ID           pgtype.UUID `json:"id"`
	Status       string      `json:"status"`
	ErrorDetails pgtype.Text `json:"error_details"`
	RowsUpserted pgtype.Int4 `json:"rows_upserted"`
	RowsRemoved  pgtype.Int4 `json:"rows_removed"`
}

// Update the status of an upload record after processing is complete or has failed
func (q *Queries) UpdateUploadStatus(ctx context.Context, arg UpdateUploadStatusParams) error {
	_, err := q.db.Exec(ctx, updateUploadStatus,
		arg.ID,
		arg.Status,
		arg.ErrorDetails,
		arg.RowsUpserted,
		arg.RowsRemoved,
	)
	return err
}
