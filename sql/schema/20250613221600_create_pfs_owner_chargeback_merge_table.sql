-- +goose Up
-- Create the "issue_owner_pfs_chargeback_merge" table

CREATE TABLE "issue_owner_pfs_chargeback_merge" (
    "user_id" UUID NOT NULL,
    "chargeback_id" UUID NOT NULL,
    PRIMARY KEY ("user_id", "chargeback_id")
);

-- +goose Down
-- Drop the "issue_owner_pfs_chargeback_merge" table

DROP TABLE IF EXISTS "issue_owner_pfs_chargeback_merge";
