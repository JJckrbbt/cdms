-- +goose Up
-- Add all custom check constraints

ALTER TABLE "chargeback" ADD CONSTRAINT chk_agreement_num_format
CHECK ("agreement_num" IS NULL OR "agreement_num" ~ '^[A-Za-z][0-9]{0,19}$');

ALTER TABLE "chargeback" ADD CONSTRAINT chk_alc_format
CHECK ("alc" ~ '^[0-9]{4,8}$');

-- +goose Down
-- Drop all custom check constraints

ALTER TABLE "chargeback" DROP CONSTRAINT IF EXISTS chk_alc_format;
ALTER TABLE "chargeback" DROP CONSTRAINT IF EXISTS chk_agreement_num_format;
