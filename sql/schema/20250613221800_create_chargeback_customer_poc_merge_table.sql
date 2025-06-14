-- +goose Up
-- Create the "chargeback_customer_poc_merge" table

CREATE TABLE "chargeback_customer_poc_merge" (
    "chargeback_id" UUID NOT NULL,
    "customer_poc_id" UUID NOT NULL,
    PRIMARY KEY ("chargeback_id", "customer_poc_id")
);

-- +goose Down
-- Drop the "chargeback_customer_poc_merge" table

DROP TABLE IF EXISTS "chargeback_customer_poc_merge";
