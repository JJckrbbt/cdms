-- +goose Up
-- Create simple, operational views for ACTIVE data.

CREATE OR REPLACE VIEW active_chargebacks_with_vendor_info AS
WITH chargeback_with_age AS (
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
    cwa.*,
    ab.agency AS agency_id,
    ab.bureau_code AS bureau_code
FROM
    chargeback_with_age cwa
JOIN
    "agency_bureau" ab ON cwa.vendor = ab."vendor_code"
WHERE
    cwa.is_active = TRUE;


CREATE OR REPLACE VIEW active_nonipac_with_vendor_info AS
WITH nonipac_with_age AS (
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
    nwa.*,
    ab.agency AS agency_id,
    ab.bureau_code AS bureau_code
FROM
    nonipac_with_age nwa
JOIN
    "agency_bureau" ab ON nwa.address_code = ab."vendor_code"
WHERE
    nwa.is_active = TRUE;

-- +goose Down
DROP VIEW IF EXISTS active_chargebacks_with_vendor_info;
DROP VIEW IF EXISTS active_nonipac_with_vendor_info;
