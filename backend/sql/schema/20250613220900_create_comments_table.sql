-- +goose Up
-- Create the "comments" table

CREATE TABLE "comments" (
    "id" UUID PRIMARY KEY,
    "comment" TEXT NOT NULL,
    "comment_date" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "user_id" UUID NOT NULL -- Foreign key will be added in a later migration
);

-- +goose Down
-- Drop the "comments" table

DROP TABLE IF EXISTS "comments";
