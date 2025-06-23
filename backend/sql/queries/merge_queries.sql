-- name: DeactivateChargebacksBySource :exec
-- Mark all existing chargebacks from a specific report source as inactive before an UPSERT 
UPDATE chargeback SET is_active = false WHERE reporting_source = $1;

-- name: UpsertChargebacks :execrows
-- Insert new records from the staging table, or update existing ones based on the business key 
-- The business key for chargebacks is BD Document Number + AL Number 
INSERT INTO chargeback (
    reporting_source, fund, business_line, region, location_system, program, al_num,
    source_num, agreement_num, title, alc, customer_tas, task_subtask, class_id,
    customer_name, org_code, document_date, accomp_date, assigned_rebill_drn,
    chargeback_amount, statement, bd_doc_num, vendor, articles_services, current_status,
    reason_code, action, is_active, updated_at
)
SELECT
    reporting_source, fund, business_line, region, location_system, program, al_num,
    source_num, agreement_num, title, alc, customer_tas, task_subtask, class_id,
    customer_name, org_code, document_date, accomp_date, assigned_rebill_drn,
    chargeback_amount, statement, bd_doc_num, vendor, articles_services, current_status,
    reason_code, action,
    true, -- Set is_active to TRUE for all records processed from the current report 
    NOW()
FROM temp_chargeback_staging
ON CONFLICT (bd_doc_num, al_num) DO UPDATE SET
    reporting_source = EXCLUDED.reporting_source,
    fund = EXCLUDED.fund,
    business_line = EXCLUDED.business_line,
    region = EXCLUDED.region,
    location_system = EXCLUDED.location_system,
    program = EXCLUDED.program,
    source_num = EXCLUDED.source_num,
    agreement_num = EXCLUDED.agreement_num,
    title = EXCLUDED.title,
    alc = EXCLUDED.alc,
    customer_tas = EXCLUDED.customer_tas,
    task_subtask = EXCLUDED.task_subtask,
    class_id = EXCLUDED.class_id,
    customer_name = EXCLUDED.customer_name,
    org_code = EXCLUDED.org_code,
    document_date = EXCLUDED.document_date,
    accomp_date = EXCLUDED.accomp_date,
    assigned_rebill_drn = EXCLUDED.assigned_rebill_drn,
    chargeback_amount = EXCLUDED.chargeback_amount,
    statement = EXCLUDED.statement,
    vendor = EXCLUDED.vendor,
    articles_services = EXCLUDED.articles_services,
    current_status = EXCLUDED.current_status,
    reason_code = EXCLUDED.reason_code,
    action = EXCLUDED.action,
    is_active = true,
    updated_at = NOW();

-- name: DeactivateNonIpacsBySource :exec
UPDATE "nonipac" SET is_active = false WHERE reporting_source = $1;

-- name: UpsertNonIpacs :execrows
-- The business key for delinquencies is Document Number 
INSERT INTO "nonipac" (
    reporting_source, business_line, billed_total_amount, principle_amount,
    interest_amount, penalty_amount, administration_charges_amount, debit_outstanding_amount,
    credit_total_amount, credit_outstanding_amount, title, document_date, address_code,
    vendor, debt_appeal_forbearance, statement, document_number, vendor_code,
    collection_due_date, open_date, is_active, updated_at
)
SELECT
    reporting_source, business_line, billed_total_amount, principle_amount,
    interest_amount, penalty_amount, administration_charges_amount, debit_outstanding_amount,
    credit_total_amount, credit_outstanding_amount, title, document_date, address_code,
    vendor, debt_appeal_forbearance, statement, document_number, vendor_code,
    collection_due_date, open_date,
    true, -- Set is_active to TRUE
    NOW()
FROM temp_nonipac_staging
ON CONFLICT (document_number) DO UPDATE SET
    reporting_source = EXCLUDED.reporting_source,
    business_line = EXCLUDED.business_line,
    billed_total_amount = EXCLUDED.billed_total_amount,
    principle_amount = EXCLUDED.principle_amount,
    interest_amount = EXCLUDED.interest_amount,
    penalty_amount = EXCLUDED.penalty_amount,
    administration_charges_amount = EXCLUDED.administration_charges_amount,
    debit_outstanding_amount = EXCLUDED.debit_outstanding_amount,
    credit_total_amount = EXCLUDED.credit_total_amount,
    credit_outstanding_amount = EXCLUDED.credit_outstanding_amount,
    title = EXCLUDED.title,
    document_date = EXCLUDED.document_date,
    address_code = EXCLUDED.address_code,
    vendor = EXCLUDED.vendor,
    debt_appeal_forbearance = EXCLUDED.debt_appeal_forbearance,
    statement = EXCLUDED.statement,
    vendor_code = EXCLUDED.vendor_code,
    collection_due_date = EXCLUDED.collection_due_date,
    open_date = EXCLUDED.open_date,
    is_active = true,
    updated_at = NOW();

-- name: UpsertAgencyBureaus :execrows
INSERT INTO agency_bureau (agency, bureau_code, vendor_code, updated_at)
SELECT agency, bureau_code, vendor_code, NOW()
FROM temp_agency_bureau_staging
ON CONFLICT (vendor_code) DO UPDATE SET
    agency = EXCLUDED.agency,
    bureau_code = EXCLUDED.bureau_code,
    updated_at = NOW();

-- name: GetChargebackSourcesByBDDocNums :many
-- For a given list of bd_doc_nums, fetch the full business key and reporting source
-- to check for cross-report conflicts in Go before an UPSERT.
SELECT bd_doc_num, al_num, reporting_source FROM chargeback
WHERE bd_doc_num = ANY($1::text[]);
