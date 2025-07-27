-- +goose Up
-- Create the "cdms_user" table with created_at, updated_at, is_active, is_admin, and is_super_admin

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
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Core tables for RBAC
CREATE TABLE "roles" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(50) UNIQUE NOT NULL,
    "description" TEXT
);

CREATE TABLE "permissions" (
    "id" SERIAL PRIMARY KEY,
    "action" VARCHAR(50) UNIQUE NOT NULL,
    "description" TEXT
);

CREATE TABLE "role_permissions" (
    "role_id" INTEGER NOT NULL REFERENCES "roles"("id") ON DELETE CASCADE,
    "permission_id" INTEGER NOT NULL REFERENCES "permissions"("id") ON DELETE CASCADE,
    PRIMARY KEY ("role_id", "permission_id")
);

CREATE TABLE "user_roles" (
    "user_id" BIGINT NOT NULL REFERENCES "cdms_user"("id") ON DELETE CASCADE,
    "role_id" INTEGER NOT NULL REFERENCES "roles"("id") ON DELETE CASCADE,
    PRIMARY KEY ("user_id", "role_id")
);

-- +goose Down
-- Drop the "user" table

DROP TABLE IF EXISTS "user_roles";
DROP TABLE IF EXISTS "role_permissions";
DROP TABLE IF EXISTS "permissions";
DROP TABLE IF EXISTS "roles";
DROP TABLE IF EXISTS "cdms_user";
