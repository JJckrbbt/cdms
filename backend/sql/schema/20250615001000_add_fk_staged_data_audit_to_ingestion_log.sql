-- +goose Up
-- Add foreign key constraint from staged_data_audit to ingestion_log

ALTER TABLE "staged_data_audit" ADD CONSTRAINT fk_staged_data_audit_ingestion_log
FOREIGN KEY ("ingestion_log_id") REFERENCES "ingestion_log" ("id");

-- Add foreign key constraint for initiated_by_user_id in ingestion_log
-- This FK is optional if you sometimes won't have a user ID.
-- If you always have a user ID (even a system user ID), you can make this NOT NULL.
ALTER TABLE "ingestion_log" ADD CONSTRAINT fk_ingestion_log_initiated_by_user
FOREIGN KEY ("initiated_by_user_id") REFERENCES "user" ("id");


-- +goose Down
-- Drop foreign key constraints (in reverse order of creation if dependencies exist)

ALTER TABLE "ingestion_log" DROP CONSTRAINT IF EXISTS fk_ingestion_log_initiated_by_user;
ALTER TABLE "staged_data_audit" DROP CONSTRAINT IF EXISTS fk_staged_data_audit_ingestion_log;
