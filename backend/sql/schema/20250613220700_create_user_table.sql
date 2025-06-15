-- +goose Up
-- Create the "user" table

CREATE TABLE "user" (
    "id" UUID PRIMARY KEY,
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "org" user_org NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE,--Account status (active/inactive)
    "is_admin" BOOLEAN NOT NULL DEFAULT FALSE -- indicates userr has admin priviledges
);

-- +goose Down
-- Drop the "user" table

DROP TABLE IF EXISTS "user";
