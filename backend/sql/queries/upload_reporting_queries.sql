-- name: GetRemovedRowsByUploadID :many
-- Fetches all rows removed by processing of a particular upload
-- most recent first
SELECT * FROM removed_rows_log
WHERE upload_id = $1
ORDER BY timestamp DESC;

