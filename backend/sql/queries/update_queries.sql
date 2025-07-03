-- name: UserUpdateChargeback :one
-- Updates the user-modifiable fields of a specific chargeback record
UPDATE chargeback
SET
    current_status = $2,
    reason_code = $3,
    action = $4,
    alc_to_rebill = $5,
    tas_to_rebill = $6,
    line_of_accounting_rebill = $7,
    special_instruction = $8,
    updated_at = NOW()
WHERE
    id = $1
RETURNING *;

-- name: PFSUpdateChargeback :one
-- Updates the user-modifiable fields of a specific chargeback record
UPDATE chargeback
SET
    current_status = $2,
    new_ipac_document_ref = $3,
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
    alc_to_rebill = $5,
    tas_to_rebill = $6,
    line_of_accounting_rebill = $7,
    special_instruction = $8,
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


