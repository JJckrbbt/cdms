-- +goose Up
-- This view is for historical reporting and includes ALL nonipac items.
CREATE OR REPLACE VIEW historical_nonipac_with_vendor_info AS
SELECT
    -- All columns from the 'nonipac' table
    ni.id,
    ni.reporting_source,
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
    ni.document_number,
    ni.vendor_code,
    ni.collection_due_date,
    ni.current_status,
    ni.pfs_poc,
    ni.gsa_poc,
    ni.customer_poc,
    ni.pfs_contacts,
    ni.open_date,
    ni.reconciled_date,
    ni.created_at,
    ni.updated_at,
    ni.is_active,
    -- The calculated date columns from our subquery for status changes
    dates.in_process_date,
    dates.referred_to_treasury_date,
    dates.closed_payment_received_date,
    dates.refund_date,
    dates.offset_date,
    dates.write_off_date,
    dates.bill_as_ipac_date,
    dates.bill_as_dod_date,
    dates.eis_issues_date,
    -- The columns from the 'agency_bureau' table
    ab.agency AS agency_id,
    ab.bureau_code,
    -- All other calculated metrics
    (NOW()::DATE - ni.document_date::DATE) AS days_old,
    ABS(ni.billed_total_amount) AS abs_amount,
    (dates.closed_payment_received_date - ni.open_date::DATE) AS days_to_close
FROM
    nonipac ni
JOIN
    "agency_bureau" ab ON ni.vendor_code = ab."vendor_code" -- Adjusted to use vendor_code from nonipac for join
LEFT JOIN
    (
        SELECT
            nsm.nonipac_id,
            MIN(CASE WHEN sh.status = 'In Process' THEN sh.status_date::DATE END) AS in_process_date,
            MIN(CASE WHEN sh.status = 'Referred to Treasury for Collections' THEN sh.status_date::DATE END) AS referred_to_treasury_date,
            MIN(CASE WHEN sh.status = 'Closed - Payment Received' THEN sh.status_date::DATE END) AS closed_payment_received_date,
            MIN(CASE WHEN sh.status = 'Refund' THEN sh.status_date::DATE END) AS refund_date,
            MIN(CASE WHEN sh.status = 'Offset' THEN sh.status_date::DATE END) AS offset_date,
            MIN(CASE WHEN sh.status = 'Write Off' THEN sh.status_date::DATE END) AS write_off_date,
            MIN(CASE WHEN sh.status = 'Bill as IPAC' THEN sh.status_date::DATE END) AS bill_as_ipac_date,
            MIN(CASE WHEN sh.status = 'Bill as DoD' THEN sh.status_date::DATE END) AS bill_as_dod_date,
            MIN(CASE WHEN sh.status = 'EIS Issues' THEN sh.status_date::DATE END) AS eis_issues_date
        FROM
            nonipac_status_merge nsm
        JOIN
            status_history sh ON nsm.status_history_id = sh.id
        GROUP BY
            nsm.nonipac_id
    ) AS dates ON ni.id = dates.nonipac_id;

-- +goose Down
DROP VIEW IF EXISTS historical_nonipac_with_vendor_info;
