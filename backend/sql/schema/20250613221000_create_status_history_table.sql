-- +goose Up
-- Create the "status_history" table

CREATE TABLE "status_history" (
    "id" UUID PRIMARY KEY,
    "status" status_history_status NOT NULL,
    "status_date" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "notes" TEXT,
    "user_id" UUID NOT NULL -- Foreign key will be added in a later migration
);

-- +goose Down
-- Drop the "status_history" table

DROP TABLE IF EXISTS "status_history";
