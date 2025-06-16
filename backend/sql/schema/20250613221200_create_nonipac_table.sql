-- +goose Up
-- Create the "nonIpac" table

CREATE TABLE "nonIpac" (
    "id" BIGSERIAL PRIMARY KEY,
    "reporting_source" nonipac_reporting_source NOT NULL,
    "business_line" chargeback_business_line NOT NULL,
    "billed_total_amount" NUMERIC(12, 2) NOT NULL,
    "principle_amount" NUMERIC(12, 2) NOT NULL,
    "interest_amount" NUMERIC(12, 2) NOT NULL,
    "penalty_amount" NUMERIC(12, 2) NOT NULL,
    "administration_charges_amount" NUMERIC(12, 2) NOT NULL,
    "debit_outstanding_amount" NUMERIC(12, 2) NOT NULL,
    "credit_total_amount" NUMERIC(12, 2) NOT NULL,
    "credit_outstanding_amount" NUMERIC(12, 2) NOT NULL,
    "title" VARCHAR(255),
    "document_date" DATE NOT NULL,
    "address_code" VARCHAR(8) NOT NULL,
    "vendor" VARCHAR(255) NOT NULL,
    "debt_appeal_forbearance" BOOLEAN NOT NULL,
    "statement" VARCHAR(8) NOT NULL,
    "document_number" VARCHAR(20) NOT NULL,
    "vendor_code" VARCHAR(8) NOT NULL,
    "collection_due_date" DATE NOT NULL,
    "current_status" nonipac_status,
    "pfs_poc" BIGINT, -- (FK to user.id)
    "gsa_poc" BIGINT, -- (FK to user.id)
    "customer_poc" BIGINT, -- (FK to customer_poc.id)
    "pfs_contacts" SMALLINT NOT NULL,
    "open_date" DATE NOT NULL,
    "reconciled_date" DATE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE
);

-- +goose Down
-- Drop the "nonIpac" table

DROP TABLE IF EXISTS "nonIpac";
