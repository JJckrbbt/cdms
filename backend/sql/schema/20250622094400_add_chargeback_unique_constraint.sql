-- +goose Up
-- This migration adds the business key unique constraint to the chargeback table.
-- This constraint is required for the INSERT ... ON CONFLICT (UPSERT) logic to work correctly,
-- as it tells PostgreSQL how to identify a duplicate record.
ALTER TABLE "chargeback"
ADD CONSTRAINT chargeback_business_key UNIQUE (bd_doc_num, al_num);

-- +goose Down
-- This migration removes the unique constraint from the chargeback table.
ALTER TABLE "chargeback"
DROP CONSTRAINT "chargeback_business_key";
