-- +goose Up
-- Create the "agency_bureau" table

CREATE TABLE "agency_bureau" (
    "agency" VARCHAR(3) NOT NULL,
    "bureau_code" VARCHAR(2) NOT NULL,
    "vendor_code" VARCHAR(8) PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);

-- +goose Down
-- Drop the "agency_bureau" table

DROP TABLE IF EXISTS "agency_bureau";
