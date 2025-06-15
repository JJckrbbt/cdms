-- +goose Up
-- Create the junction table for many-to-many user-to-business-line access

CREATE TABLE "user_business_line_access" (
    "user_id" UUID NOT NULL,
    "business_line" chargeback_business_line NOT NULL, -- Reuses the ENUM for consistency
    PRIMARY KEY ("user_id", "business_line")
);

-- +goose Down
-- Drop the user_business_line_access table

DROP TABLE IF EXISTS "user_business_line_access";
