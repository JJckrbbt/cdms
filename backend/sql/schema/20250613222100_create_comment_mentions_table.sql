-- +goose Up
-- Create the "comment_mentions" table

CREATE TABLE "comment_mentions" (
    "comment_id" BIGINT NOT NULL,
    "user_id" BIGINT NOT NULL,
    PRIMARY KEY ("comment_id", "user_id")
);

-- +goose Down
-- Drop the "comment_mentions" table

DROP TABLE IF EXISTS "comment_mentions";
