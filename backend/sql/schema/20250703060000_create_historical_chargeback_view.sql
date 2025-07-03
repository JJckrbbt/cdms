-- +goose Up
-- This view is for historical reporting and includes ALL chargebacks.
-- It is built using a subquery in the JOIN clause to ensure there is no ambiguity.
CREATE OR REPLACE VIEW historical_chargebacks_with_vendor_info AS
SELECT
    -- All columns from the 'chargeback' table
    cb.id,
    cb.is_active,
    cb.reporting_source,
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
    cb.reason_code,
    cb.action,
    cb.alc_to_rebill,
    cb.tas_to_rebill,
    cb.line_of_accounting_rebill,
    cb.special_instruction,
    cb.new_ipac_document_ref,
    cb.created_at,
    cb.updated_at,
    -- The calculated date columns from our subquery
    dates.issue_in_research_date,
    dates.passed_to_pfs_date,
    dates.pfs_completion_date,
    -- The columns from the 'agency_bureau' table
    ab.agency AS agency_id,
    ab.bureau_code,
    -- All other calculated metrics
    (CASE WHEN cb.accomp_date IS NOT NULL THEN (NOW()::DATE - cb.accomp_date::DATE) ELSE (NOW()::DATE - cb.document_date::DATE) END) AS days_old,
    ABS(cb.chargeback_amount) AS abs_amount,
    (dates.passed_to_pfs_date - cb.created_at::DATE) AS days_open_to_pfs,
    (dates.pfs_completion_date - dates.passed_to_pfs_date) AS days_pfs_to_complete,
    (NOW()::DATE - dates.pfs_completion_date) AS days_complete
FROM
    chargeback cb
JOIN
    "agency_bureau" ab ON cb.vendor = ab."vendor_code"
LEFT JOIN
    (
        SELECT
            csm.chargeback_id,
            MIN(CASE WHEN sh.status = 'In Research' THEN sh.status_date::DATE END) AS issue_in_research_date,
            MIN(CASE WHEN sh.status = 'Passed to PFS' THEN sh.status_date::DATE END) AS passed_to_pfs_date,
            MIN(CASE WHEN sh.status = 'Completed by PFS' THEN sh.status_date::DATE END) AS pfs_completion_date
        FROM
            chargeback_status_merge csm
        JOIN
            status_history sh ON csm.status_history_id = sh.id
        GROUP BY
            csm.chargeback_id
    ) AS dates ON cb.id = dates.chargeback_id;

-- +goose Down
DROP VIEW IF EXISTS historical_chargebacks_with_vendor_info;
