-- +goose Up
-- Create the "comments" table with created_at and updated_at

CREATE TABLE "comments" (
    "id" BIGSERIAL PRIMARY KEY,
    "comment" TEXT NOT NULL,
    "comment_date" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "user_id" BIGINT NOT NULL, -- (FK to user.id)
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
-- Drop the "comments" table

DROP TABLE IF EXISTS "comments";
