-- +goose Up
-- Create the "chargeback_comments_merge" table

CREATE TABLE "chargeback_comments_merge" (
    "chargeback_id" UUID NOT NULL,
    "comment_id" UUID NOT NULL,
    PRIMARY KEY ("chargeback_id", "comment_id")
);

-- +goose Down
-- Drop the "chargeback_comments_merge" table

DROP TABLE IF EXISTS "chargeback_comments_merge";
