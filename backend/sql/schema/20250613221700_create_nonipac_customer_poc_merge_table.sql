-- +goose Up
-- Create the "non_ipac_customer_poc_merge" table

CREATE TABLE "non_ipac_customer_poc_merge" (
    "nonipac_id" BIGINT NOT NULL,
    "customer_poc_id" BIGINT NOT NULL,
    PRIMARY KEY ("nonipac_id", "customer_poc_id")
);

-- +goose Down
-- Drop the "non_ipac_customer_poc_merge" table

DROP TABLE IF EXISTS "non_ipac_customer_poc_merge";
