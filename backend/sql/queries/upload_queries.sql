-- name: CreateUpload :one
-- Create a record to track a new file upload 
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
RETURNING *;

-- name: UpdateUploadStatus :exec
-- Update the status of an upload record after processing is complete or has failed 
UPDATE uploads
SET
    status = $2,
    error_details = $3,
    rows_upserted = $4,
    rows_removed = $5,
    processed_at = NOW()
WHERE id = $1;

-- name: GetUpload :one
-- Retrieve a detailed summary for a specific upload 
SELECT * FROM uploads
WHERE id = $1;

-- name: ListUploads :many
-- Provides a paginated list of recent report uploads and their statuses 
SELECT * FROM uploads
ORDER BY uploaded_at DESC
LIMIT $1
OFFSET $2;
