-- name: CreateChargeback :one
-- Inserts a new chargeback record,from a manual UI entry.
-- The 'reporting_source' is hardcoded to 'ApplicationCreated'.
INSERT INTO chargeback (
    reporting_source,
    fund,
    business_line,
    region,
    program,
    al_num,
    source_num,
    alc,
    customer_tas,
    task_subtask,
    customer_name,
    org_code,
    document_date,
    accomp_date,
    chargeback_amount,
    statement,
    bd_doc_num,
    vendor,
    -- Nullable fields
    location_system,
    agreement_num,
    title,
    class_id,
    assigned_rebill_drn,
    articles_services,
    current_status,
    reason_code,
    action
) VALUES (
    'ApplicationCreated',
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
)
RETURNING *;

-- name: CreateDelinquency :one
-- Inserts a new delinquency (nonipac) record, from a manual UI entry.
-- The 'reporting_source' is hardcoded to 'ApplicationCreated'.
INSERT INTO "nonipac" (
    reporting_source,
    business_line,
    billed_total_amount,
    principle_amount,
    interest_amount,
    penalty_amount,
    administration_charges_amount,
    debit_outstanding_amount,
    credit_total_amount,
    credit_outstanding_amount,
    document_date,
    address_code,
    vendor,
    debt_appeal_forbearance,
    statement,
    document_number,
    vendor_code,
    collection_due_date,
    open_date,
    -- Nullable fields
    title,
    current_status
) VALUES (
    'ApplicationCreated',
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
)
RETURNING *;
