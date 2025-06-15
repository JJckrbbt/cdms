-- +goose Up
-- Create the "staged_data_audit" table for logging problematic/audited rows from staging

CREATE TABLE "staged_data_audit" (
    "audit_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "ingestion_log_id" UUID NOT NULL, -- Foreign key will be added in a later migration
    "target_table" VARCHAR(50) NOT NULL, -- The main table this data was intended for (e.g., 'chargeback', 'nonIpac')
    "original_row_identifier" TEXT, -- An identifier from the original source (e.g., line number, unique ID)
    "staged_data_jsonb" JSONB NOT NULL, -- The full raw data of the problematic row from the staging table
    "audit_timestamp" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "audit_status" VARCHAR(50) NOT NULL, -- e.g., 'VALIDATION_FAILED', 'SKIPPED', 'DUPLICATE', 'PARSED_OK'
    "audit_message" TEXT -- Detailed message about the audit status or error
);

-- Add index for efficient lookup by ingestion log ID
CREATE INDEX idx_staged_data_audit_log_id ON "staged_data_audit" (ingestion_log_id);
-- Add index for lookup by target table and audit status
CREATE INDEX idx_staged_data_audit_target_status ON "staged_data_audit" (target_table, audit_status);
-- Add index for chronological ordering
CREATE INDEX idx_staged_data_audit_timestamp ON "staged_data_audit" (audit_timestamp DESC);

-- +goose Down
-- Drop the "staged_data_audit" table

DROP TABLE IF EXISTS "staged_data_audit";
