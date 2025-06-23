-- +goose Up
-- Create the "chargeback" table

CREATE TABLE "chargeback" (
    "id" BIGSERIAL PRIMARY KEY,
    "reporting_source" chargeback_reporting_source NOT NULL,
    "fund" chargeback_fund NOT NULL,
    "business_line" chargeback_business_line NOT NULL,
    "region" SMALLINT NOT NULL,
    "location_system" VARCHAR(3),
    "program" VARCHAR(4) NOT NULL,
    "al_num" SMALLINT NOT NULL,
    "source_num" VARCHAR(30) NOT NULL,
    "agreement_num" VARCHAR(20),
    "title" VARCHAR(255),
    "alc" VARCHAR(8) NOT NULL,
    "customer_tas" VARCHAR(30) NOT NULL,
    "task_subtask" VARCHAR(20) NOT NULL,
    "class_id" VARCHAR(10),
    "customer_name" VARCHAR(255) NOT NULL,
    "org_code" VARCHAR(8) NOT NULL,
    "document_date" DATE NOT NULL,
    "accomp_date" DATE NOT NULL,
    "assigned_rebill_drn" VARCHAR(8),
    "chargeback_amount" NUMERIC(12, 2) NOT NULL,
    "statement" VARCHAR(8) NOT NULL,
    "bd_doc_num" VARCHAR(20) NOT NULL,
    "vendor" VARCHAR(8) NOT NULL,
    "articles_services" TEXT,
    "current_status" chargeback_status NOT NULL DEFAULT 'New',
    "issue_in_research_date" DATE,
    "reason_code" chargeback_reason_code,
    "action" chargeback_action,
    "alc_to_rebill" VARCHAR(50),
    "tas_to_rebill" VARCHAR(50),
    "line_of_accounting_rebill" TEXT,
    "special_instruction" VARCHAR(200),
    "new_ipac_document_ref" VARCHAR(30),
    "pfs_completion_date" DATE,
    "reconciliation_date" DATE,
    "chargeback_count" SMALLINT,
    "passed_to_psf" DATE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE
);

-- +goose Down
-- Drop the "chargeback" table

DROP TABLE IF EXISTS "chargeback";
