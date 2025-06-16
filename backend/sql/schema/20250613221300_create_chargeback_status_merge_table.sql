-- +goose Up
-- Create the "chargeback_status_merge" table

CREATE TABLE "chargeback_status_merge" (
    "chargeback_id" BIGINT NOT NULL,
    "status_history_id" BIGINT NOT NULL,
    PRIMARY KEY ("chargeback_id", "status_history_id")
);

-- +goose Down
-- Drop the "chargeback_status_merge" table

DROP TABLE IF EXISTS "chargeback_status_merge";
