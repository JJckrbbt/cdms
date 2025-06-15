-- +goose Up
-- Create views for active chargebacks and nonIpac records, joined with vendor info

CREATE VIEW active_chargebacks_with_vendor_info AS
SELECT
    cb.id,
    cb.fund,
    cb.business_line,
    cb.region,
    cb.location_system,
    cb.program,
    cb.al_num,
    cb.source_num,
    cb.agreement_num,
    cb.title,
    cb.alc,
    cb.customer_tas,
    cb.task_subtask,
    cb.class_id,
    cb.customer_name,
    cb.org_code,
    cb.document_date,
    cb.accomp_date,
    cb.assigned_rebill_drn,
    cb.chargeback_amount,
    cb.statement,
    cb.bd_doc_num,
    cb.vendor,
    cb.articles_services,
    cb.current_status,
    cb.issue_in_research_date,
    cb.reason_code,
    cb.action,
    cb.alc_to_rebill,
    cb.tas_to_rebill,
    cb.line_of_accounting_rebill,
    cb.special_instruction,
    cb.new_ipac_document_ref,
    cb.pfs_completion_date,
    cb.reconciliation_date,
    cb.chargeback_count,
    cb.passed_to_psf,
    cb.created_at,
    cb.updated_at,
    ab.agency AS agency_id,
    ab.bureau_code AS bureau_code
FROM
    "chargeback" cb
JOIN
    "agency_bureau" ab ON cb.vendor = ab."vendor_code"
WHERE
    cb.is_active = TRUE;


CREATE VIEW active_nonipac_with_vendor_info AS
SELECT
    ni.id,
    ni.business_line,
    ni.billed_total_amount,
    ni.principle_amount,
    ni.interest_amount,
    ni.penalty_amount,
    ni.administration_charges_amount,
    ni.debit_outstanding_amount,
    ni.credit_total_amount,
    ni.credit_outstanding_amount,
    ni.title,
    ni.document_date,
    ni.address_code,
    ni.vendor,
    ni.debt_appeal_forbearance,
    ni.statement,
    ni.current_status,
    ni.document_number,
    ni.vendor_code,
    ni.collection_due_date,
    ni.pfs_poc,
    ni.gsa_poc,
    ni.customer_poc,
    ni.pfs_contacts,
    ni.open_date,
    ni.reconciled_date,
    ni.created_at,
    ni.updated_at,
    ab.agency AS agency_id,
    ab.bureau_code AS bureau_code
FROM
    "nonIpac" ni
JOIN
    "agency_bureau" ab ON ni.address_code = ab."vendor_code"
WHERE
    ni.is_active = TRUE;

-- +goose Down
-- Drop the views

DROP VIEW IF EXISTS active_nonipac_with_vendor_info;
DROP VIEW IF EXISTS active_chargebacks_with_vendor_info;
