-- //go:generate mockery --name Querier --output ./mocks --outpkg mocks
-- name: ListActiveChargebacks :many
-- Fetches a paginated list from the active_chargebacks_with_vendor_info view.
-- The view is already filtered by is_active = true.
SELECT *, count(*) OVER() AS total_count
FROM active_chargebacks_with_vendor_info
ORDER BY document_date DESC
LIMIT $1
OFFSET $2;

-- name: ListActiveDelinquencies :many
-- Fetches a paginated list from the active_nonipac_with_vendor_info view.
-- The view is already filtered by is_active = true.
SELECT *, count(*) OVER() AS total_count
FROM active_nonipac_with_vendor_info
ORDER BY document_date DESC
LIMIT $1
OFFSET $2;

-- name: GetActiveChargebackByID :one
-- Fetches a single active chargeback by its primary key from the view.
SELECT * FROM active_chargebacks_with_vendor_info
WHERE id = $1;

-- name: GetActiveDelinquencyByID :one
-- Fetches a single active delinquency by its primary key from the view.
SELECT * FROM active_nonipac_with_vendor_info
WHERE id = $1;

-- name: GetActiveChargebackByBusinessKey :one
-- Fetches a single active chargeback by business key
SELECT * FROM active_chargebacks_with_vendor_info
WHERE bd_doc_num = $1 AND al_num = $2;

-- name: GetActiveDelinquencyByBusinessKey :one
-- Fetches a single active delinquency by business key
SELECT * FROM active_nonipac_with_vendor_info
WHERE document_number = $1;

-- name: GetChargebackForUpdate :one
-- Fetches a single chargeback directly from the base table for updating.
SELECT * FROM chargeback
WHERE id = $1 LIMIT 1;

-- name: GetDelinquencyForUpdate :one
-- Fetches a single chargeback directly from the base table for updating.
SELECT * FROM nonipac
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email = $1;

-- name: GetStatusHistoryForChargeback :many
-- Fetches Status History for Chargebacks
SELECT
    sh.id as status_history_id, 
    sh.status,
    sh.status_date, 
    sh.notes,
    sh.user_id,
    u.first_name AS user_first_name,
    u.last_name AS user_last_name,
    u.email AS user_email
FROM
    "status_history" sh
JOIN
    "chargeback_status_merge" csm ON sh.id = csm.status_history_id
JOIN
    "user" u ON sh.user_id = u.id
WHERE
    csm.chargeback_id = $1 
ORDER BY
    sh.status_date DESC;

-- name: GetStatusHistoryForDelinquencies :many
-- Fetches Status History for Delinquencies
SELECT
    sh.id as status_history_id, 
    sh.status,
    sh.status_date, 
    sh.notes,
    sh.user_id,
    u.first_name AS user_first_name,
    u.last_name AS user_last_name,
    u.email AS user_email
FROM
    "status_history" sh
JOIN
    "nonipac_status_merge" nsm ON sh.id = nsm.status_history_id
JOIN
    "user" u ON sh.user_id = u.id
WHERE
    nsm.nonipac_id = $1 
ORDER BY
    sh.status_date DESC;
