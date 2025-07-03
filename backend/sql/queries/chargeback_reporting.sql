-- name: GetChargebackStatusSummary :many
-- Gets the count, total value, and percentage of total value for each chargeback status for active items.
SELECT
    current_status,
    COUNT(*) AS status_count,
    SUM(abs_amount)::NUMERIC AS total_value,
    (SUM(abs_amount) * 100.0 / SUM(SUM(abs_amount)) OVER ())::NUMERIC(5, 2) AS percentage_of_total
FROM
    historical_chargebacks_with_vendor_info
WHERE
    current_status != 'Reconciled - Off Report' -- Exclude reconciled items from this summary
GROUP BY
    current_status;

-- name: GetNewChargebackStatsForWindow :one
-- Gets the count and total value of new chargebacks created within a specific date window.
SELECT
    COUNT(*) AS new_chargebacks_count,
    COALESCE(SUM(chargeback_amount), 0)::NUMERIC(12, 2) AS new_chargebacks_value
FROM
    chargeback
WHERE
    created_at BETWEEN $1 AND $2;


-- name: GetAverageDaysToPFSForWindow :one
-- Gets the average days from creation to 'Passed to PFS' for chargebacks passed within a specific date window.
SELECT
    COALESCE(AVG(days_open_to_pfs), 0)::NUMERIC(10, 2) AS avg_days
FROM
    historical_chargebacks_with_vendor_info
WHERE
    passed_to_pfs_date BETWEEN $1 AND $2;

-- name: GetAverageDaysForPFSCompletionForWindow :one
-- Gets the average days for PFS to complete for chargebacks completed within a specific date window.
SELECT
    COALESCE(AVG(days_pfs_to_complete), 0)::NUMERIC(10, 2) AS avg_days
FROM
    historical_chargebacks_with_vendor_info
WHERE
    pfs_completion_date BETWEEN $1 AND $2;

-- name: GetPFSCountsForWindow :one
-- Gets the count of chargebacks passed to PFS and completed by PFS within a specific date window.
-- This version uses conditional aggregation for better performance and to avoid ambiguity.
SELECT
    COUNT(*) FILTER (WHERE passed_to_pfs_date BETWEEN $1 AND $2) AS passed_to_pfs_count,
    COUNT(*) FILTER (WHERE pfs_completion_date BETWEEN $1 AND $2) AS completed_by_pfs_count
FROM
    historical_chargebacks_with_vendor_info;
