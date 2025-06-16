-- +goose Up
-- Create the "non_ipac_comments_merge" table

CREATE TABLE "non_ipac_comments_merge" (
    "nonipac_id" BIGINT NOT NULL,
    "comment_id" BIGINT NOT NULL,
    PRIMARY KEY ("nonipac_id", "comment_id")
);

-- +goose Down
-- Drop the "non_ipac_comments_merge" table

DROP TABLE IF EXISTS "non_ipac_comments_merge";
