-- +goose Up
-- Create the "nonipac_status_merge" table

CREATE TABLE "nonipac_status_merge" (
    "nonipac_id" BIGINT NOT NULL,
    "status_history_id" BIGINT NOT NULL,
    PRIMARY KEY ("nonipac_id", "status_history_id")
);

-- +goose Down
-- Drop the "nonipac_status_merge" table

DROP TABLE IF EXISTS "nonipac_status_merge";
