-- +goose Up
-- Create the "agency_bureau" table

CREATE TABLE "agency_bureau" (
    "agency" VARCHAR(3) NOT NULL,
    "gsa_bureau_code" VARCHAR(2) NOT NULL,
    "vendor_code" VARCHAR(8) PRIMARY KEY
);

-- +goose Down
-- Drop the "agency_bureau" table

DROP TABLE IF EXISTS "agency_bureau";
