-- name: UserUpdateChargeback :one
-- Updates the user-modifiable fields of a specific chargeback record
UPDATE chargeback
SET
    current_status = $2,
    reason_code = $3,
    action = $4,
    issue_in_research_date = $5,
    alc_to_rebill = $6,
    tas_to_rebill = $7,
    line_of_accounting_rebill = $8,
    special_instruction = $9,
    passed_to_psf = $10,
    updated_at = NOW()
WHERE
    id = $1
RETURNING *;

-- name: PFSUpdateChargeback :one
-- Updates the user-modifiable fields of a specific chargeback record
UPDATE chargeback
SET
    current_status = $2,
    passed_to_psf = $3,
    new_ipac_document_ref = $4,
    pfs_completion_date = $5,
    updated_at = NOW()
WHERE
    id = $1
RETURNING *;

-- name: AdminUpdateChargeback :one
-- Updates the admin-modifiable fields of a specific chargeback record
UPDATE chargeback
SET
    current_status = $2,
    reason_code = $3,
    action = $4,
    issue_in_research_date = $5,
    alc_to_rebill = $6,
    tas_to_rebill = $7,
    line_of_accounting_rebill = $8,
    special_instruction = $9,
    passed_to_psf = $10,
    pfs_completion_date = $11,
    updated_at = NOW()
WHERE
    id = $1
RETURNING *;

-- name: UserUpdateDelinquency :one
-- Updates the user-modifiable fields of a specific delinquency record
UPDATE nonipac
SET
    current_status = $2,
    updated_at = NOW()
WHERE
    id = $1
RETURNING *;

-- name: PFSUpdateDelinquency :one
-- Updates the user-modifiable fields of a specific delinquency record
UPDATE nonipac
SET
    current_status = $2,
    updated_at = NOW()
WHERE
    id = $1
RETURNING *;

-- name: AdminUpdateDelinquency :one
-- Updates the admin-modifiable fields of a specific delinquency record
UPDATE nonipac
SET
    current_status = $2,
    updated_at = NOW()
WHERE
    id = $1
RETURNING *;


