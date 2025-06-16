-- +goose Up
-- Create the "user" table with created_at, updated_at, is_active, and is_admin

CREATE TABLE "user" (
    "id" BIGSERIAL PRIMARY KEY,
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "org" user_org NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
    "is_admin" BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
-- Drop the "user" table

DROP TABLE IF EXISTS "user";
