-- +goose Up
-- Create the "customer_poc" table with created_at, updated_at, and is_active

CREATE TABLE "customer_poc" (
    "id" BIGSERIAL PRIMARY KEY,
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "email" VARCHAR(255),
    "phone" VARCHAR(50),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE
);

-- +goose Down
-- Drop the "customer_poc" table

DROP TABLE IF EXISTS "customer_poc";
