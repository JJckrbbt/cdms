-- +goose Up
-- Create the "ingestion_log" table for logging merge process metadata

CREATE TABLE "ingestion_log" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "process_name" VARCHAR(100) NOT NULL, -- e.g., 'Chargeback_Merge_Process', 'NonIpac_Import'
    "started_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "completed_at" TIMESTAMPTZ, -- Will be set upon completion
    "status" VARCHAR(50) NOT NULL, -- e.g., 'SUCCESS', 'FAILED', 'IN_PROGRESS'
    "rows_processed" INTEGER, -- Total rows attempted to process from source
    "rows_inserted" INTEGER,  -- Rows successfully inserted into main tables
    "rows_updated" INTEGER,   -- Rows successfully updated in main tables
    "rows_failed" INTEGER,    -- Rows that failed validation/merge
    "source_identifier" VARCHAR(255), -- e.g., original filename, S3 object key, API call ID
    "initiated_by_user_id" BIGINT -- FK to your "user" table (can be NULL for automated processes)
);

-- Add index for efficient lookup by process name and status
CREATE INDEX idx_ingestion_log_process_status ON "ingestion_log" (process_name, status);
-- Add index for chronological ordering
CREATE INDEX idx_ingestion_log_started_at ON "ingestion_log" (started_at DESC);

-- +goose Down
-- Drop the "ingestion_log" table

DROP TABLE IF EXISTS "ingestion_log";
