-- +goose Up
-- Create the "user" table

CREATE TABLE "user" (
    "id" UUID PRIMARY KEY,
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "org" user_org NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL
);

-- +goose Down
-- Drop the "user" table

DROP TABLE IF EXISTS "user";
