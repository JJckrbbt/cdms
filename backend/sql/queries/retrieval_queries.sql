-- //go:generate mockery --name Querier --output ./mocks --outpkg mocks
-- name: ListActiveChargebacks :many
-- Fetches a paginated list from the active_chargebacks_with_vendor_info view.
-- The view is already filtered by is_active = true.
SELECT * FROM active_chargebacks_with_vendor_info
ORDER BY document_date DESC
LIMIT $1
OFFSET $2;

-- name: ListActiveDelinquencies :many
-- Fetches a paginated list from the active_nonipac_with_vendor_info view.
-- The view is already filtered by is_active = true.
SELECT * FROM active_nonipac_with_vendor_info
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

