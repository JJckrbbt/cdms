-- This file defines the structure of our temporary staging tables.
-- It is NOT a goose migration. It is only here to inform sqlc
-- about the shape of the tables we create at runtime.

-- The structure is identical to the main 'chargeback' table.
CREATE TABLE "temp_chargeback_staging" (LIKE "chargeback" INCLUDING DEFAULTS);

-- The structure is identical to the main 'nonIpac' table.
CREATE TABLE "temp_nonipac_staging" (LIKE "nonipac" INCLUDING DEFAULTS);

-- The structure is identical to the main 'agency_bureau' table.
CREATE TABLE "temp_agency_bureau_staging" (LIKE "agency_bureau" INCLUDING DEFAULTS);
