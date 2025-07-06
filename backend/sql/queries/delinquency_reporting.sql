-- name: GetNonipacStatusSummary :many
-- Gets the count, total value, and percentage of total value for each nonipac status for active items.
SELECT
    current_status,
    COUNT(*) AS status_count,
    SUM(abs_amount)::NUMERIC AS total_value,
    (SUM(abs_amount) * 100.0 / SUM(SUM(abs_amount)) OVER ())::NUMERIC(5, 2) AS percentage_of_total
FROM
    historical_nonipac_with_vendor_info
WHERE
    is_active = TRUE AND current_status != 'Reconciled - Off Report' -- Exclude inactive and reconciled items
GROUP BY
    current_status
ORDER BY
    current_status;

-- name: GetNonipacAgingScheduleByBusinessLine :many
-- Provides an aging schedule for active nonipac items, broken down by business line and age categories.
SELECT
    business_line,
    COUNT(*) FILTER (WHERE days_old <= 180) AS "less_than_180_days_count",
    COALESCE(SUM(abs_amount) FILTER (WHERE days_old <= 180), 0)::NUMERIC(12, 2)::TEXT AS "less_than_180_days_value",
    COUNT(*) FILTER (WHERE days_old BETWEEN 181 AND 365) AS "181_to_365_days_count",
    COALESCE(SUM(abs_amount) FILTER (WHERE days_old BETWEEN 181 AND 365), 0)::NUMERIC(12, 2)::TEXT AS "181_to_365_days_value",
    COUNT(*) FILTER (WHERE days_old BETWEEN 366 AND 730) AS "one_to_two_years_count",
    COALESCE(SUM(abs_amount) FILTER (WHERE days_old BETWEEN 366 AND 730), 0)::NUMERIC(12, 2)::TEXT AS "one_to_two_years_value",
    COUNT(*) FILTER (WHERE days_old > 730) AS "over_two_years_count",
    COALESCE(SUM(abs_amount) FILTER (WHERE days_old > 730), 0)::NUMERIC(12, 2)::TEXT AS "over_two_years_value",
    COUNT(*) AS "total_count",
    SUM(abs_amount)::NUMERIC(12, 2)::TEXT AS "total_value"
FROM
    historical_nonipac_with_vendor_info
WHERE
    is_active = TRUE AND current_status != 'Reconciled - Off Report' -- Filter for active, non-reconciled items
GROUP BY
    business_line
ORDER BY
    business_line;
