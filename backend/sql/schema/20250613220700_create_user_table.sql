-- +goose Up
-- Create the "cdms_user" table with created_at, updated_at, is_active, and is_admin

CREATE TABLE "cdms_user" (
    "id" BIGSERIAL PRIMARY KEY,
    --External Auth Provider ID & Email Provided---
    "auth_provider_subject" VARCHAR(255) UNIQUE NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    --Internal Application Fields---
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "org" user_org NOT NULL,
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
    "is_admin" BOOLEAN NOT NULL DEFAULT FALSE,

    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
-- Drop the "user" table

DROP TABLE IF EXISTS "cdms_user";
