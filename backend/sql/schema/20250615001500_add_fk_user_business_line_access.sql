-- +goose Up
-- Add foreign key constraint for user_business_line_access table

ALTER TABLE "user_business_line_access" ADD CONSTRAINT fk_user_business_line_access_user
FOREIGN KEY ("user_id") REFERENCES "user" ("id");

-- Note: No FK for business_line as it's an ENUM type, not a separate table.

-- +goose Down
-- Drop foreign key constraint for user_business_line_access table

ALTER TABLE "user_business_line_access" DROP CONSTRAINT IF EXISTS fk_user_business_line_access_user;
