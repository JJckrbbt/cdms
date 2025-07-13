-- +goose Up

-- An `uploads` table will track the status and metadata of each report file uploaded 
CREATE TABLE "uploads" (
    "id" UUID PRIMARY KEY, -- 
    "storage_key" TEXT NOT NULL, -- 
    "filename" TEXT NOT NULL, -- 
    "report_type" TEXT NOT NULL, -- 
    "status" TEXT NOT NULL, -- 
    "uploaded_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- 
    "processed_at" TIMESTAMPTZ, -- 
    "error_details" TEXT, -- 
    "processed_by_user_id" BIGINT NOT NULL REFERENCES "user" ("id"),
    "rows_upserted" INTEGER,
    "rows_removed" INTEGER
);

-- A new removed_rows_log table will store details about rows that were excluded during report processing 
CREATE TABLE "removed_rows_log" (
    "id" UUID PRIMARY KEY, -- 
    "upload_id" UUID NOT NULL REFERENCES "uploads" ("id") ON DELETE CASCADE, -- 
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- 
    "report_type" TEXT NOT NULL, -- 
    "original_row_data" JSONB NOT NULL, -- 
    "reason_for_removal" TEXT NOT NULL -- 
);


-- +goose Down
-- This section reverses the migration, removing the tables in the reverse order of creation.

DROP TABLE "removed_rows_log";
DROP TABLE "uploads";
