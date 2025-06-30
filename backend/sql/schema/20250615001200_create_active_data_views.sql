-- +goose Up
-- Create views for active chargebacks and nonipac records, joined with vendor info

CREATE OR REPLACE VIEW active_chargebacks_with_vendor_info AS
WITH chargeback_dates AS (
    -- This CTE gets all the necessary dates from status history
    SELECT
        cb.id,
        MIN(CASE WHEN sh.status = 'Passed to PFS' THEN sh.status_date::DATE END) AS passed_to_pfs_date,
        MIN(CASE WHEN sh.status = 'Completed by PFS' THEN sh.status_date::DATE END) AS completed_by_pfs_date
    FROM
        chargeback cb
    LEFT JOIN
        chargeback_status_merge csm ON cb.id = csm.chargeback_id
    LEFT JOIN
        status_history sh ON csm.status_history_id = sh.id
    GROUP BY
        cb.id
),
chargeback_with_age AS (
    -- This CTE calculates the initial days_old value
    SELECT
        cb.*,
        (
            CASE
                WHEN cb.accomp_date IS NOT NULL THEN (NOW()::DATE - cb.accomp_date::DATE)
                ELSE (NOW()::DATE - cb.document_date::DATE)
            END
        ) AS days_old
    FROM
        chargeback cb
)
SELECT
    cbwa.id,
    cbwa.fund,
    cbwa.business_line,
    cbwa.region,
    cbwa.location_system,
    cbwa.program,
    cbwa.al_num,
    cbwa.source_num,
    cbwa.agreement_num,
    cbwa.title,
    cbwa.alc,
    cbwa.customer_tas,
    cbwa.class_id,
    cbwa.customer_name,
    cbwa.org_code,
    cbwa.document_date,
    cbwa.accomp_date,
    cbwa.assigned_rebill_drn,
    cbwa.chargeback_amount,
    cbwa.statement,
    cbwa.bd_doc_num,
    cbwa.vendor,
    cbwa.articles_services,
    cbwa.current_status,
    cbwa.issue_in_research_date,
    cbwa.reason_code,
    cbwa.action,
    cbwa.alc_to_rebill,
    cbwa.tas_to_rebill,
    cbwa.line_of_accounting_rebill,
    cbwa.special_instruction,
    cbwa.new_ipac_document_ref,
    cbwa.pfs_completion_date,
    cbwa.reconciliation_date,
    cbwa.chargeback_count,
    cbwa.passed_to_psf,
    cbwa.created_at,
    cbwa.updated_at,
    ab.agency AS agency_id,
    ab.bureau_code AS bureau_code,
    cbwa.days_old,
    ABS(cbwa.chargeback_amount) AS abs_amount,
    CASE
        WHEN cbwa.days_old <= 90 THEN 'Less than 90 days'
        WHEN cbwa.days_old > 90 AND cbwa.days_old <= 365 THEN '91 - 365 days'
        WHEN cbwa.days_old > 365 AND cbwa.days_old <= 730 THEN 'One - Two Years Old'
        ELSE 'Over Two Years Old'
    END AS aging_category,
    (cd.passed_to_pfs_date - cbwa.created_at::DATE) AS days_open_to_pfs,
    (cd.completed_by_pfs_date - cd.passed_to_pfs_date) AS days_pfs_to_complete,
    (NOW()::DATE - cd.completed_by_pfs_date) AS days_complete
FROM
    chargeback_with_age cbwa
LEFT JOIN
    chargeback_dates cd ON cbwa.id = cd.id
JOIN
    "agency_bureau" ab ON cbwa.vendor = ab."vendor_code"
WHERE
    cbwa.is_active = TRUE;

CREATE OR REPLACE VIEW active_nonipac_with_vendor_info AS
WITH nonipac_with_age AS (
    -- This CTE calculates the initial days_old value
    SELECT
        ni.*,
        (
            CASE
                WHEN ni.document_date IS NOT NULL THEN (NOW()::DATE - ni.document_date::DATE)
                ELSE (NOW()::DATE - ni.collection_due_date::DATE)
            END
        ) AS days_old
    FROM
        nonipac ni
)
SELECT
    nwa.id,
    nwa.business_line,
    nwa.billed_total_amount,
    nwa.principle_amount,
    nwa.interest_amount,
    nwa.penalty_amount,
    nwa.administration_charges_amount,
    nwa.debit_outstanding_amount,
    nwa.credit_total_amount,
    nwa.credit_outstanding_amount,
    nwa.title,
    nwa.document_date,
    nwa.address_code,
    nwa.vendor,
    nwa.debt_appeal_forbearance,
    nwa.statement,
    nwa.current_status,
    nwa.document_number,
    nwa.vendor_code,
    nwa.collection_due_date,
    nwa.pfs_poc,
    nwa.gsa_poc,
    nwa.customer_poc,
    nwa.pfs_contacts,
    nwa.open_date,
    nwa.reconciled_date,
    nwa.created_at,
    nwa.updated_at,
    ab.agency AS agency_id,
    ab.bureau_code AS bureau_code,
    nwa.days_old,
        CASE
            WHEN nwa.days_old <= 90 THEN 'Less than 90 days'
            WHEN nwa.days_old > 90 AND nwa.days_old <= 365 THEN '91 - 365 days'
            WHEN nwa.days_old > 365 AND nwa.days_old <= 730 THEN 'One - Two Years Old'
            ELSE 'Over Two Years Old'
        END AS aging_category,
        ABS(nwa.billed_total_amount) AS abs_amount
    FROM
        nonipac_with_age nwa
    JOIN
        "agency_bureau" ab ON nwa.address_code = ab."vendor_code"
    WHERE
        nwa.is_active = TRUE;

-- +goose Down
-- Drop the views

DROP VIEW IF EXISTS active_nonipac_with_vendor_info;
DROP VIEW IF EXISTS active_chargebacks_with_vendor_info;


