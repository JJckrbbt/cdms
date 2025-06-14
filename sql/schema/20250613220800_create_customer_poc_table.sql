-- +goose Up
-- Create the "customer_poc" table

CREATE TABLE "customer_poc" (
    "id" UUID PRIMARY KEY,
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "email" VARCHAR(255),
    "phone" VARCHAR(50)
);

-- +goose Down
-- Drop the "customer_poc" table

DROP TABLE IF EXISTS "customer_poc";
