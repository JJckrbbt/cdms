-- +goose Up
-- +goose StatementBegin
-- Create a function to set the updated_at timestamp
CREATE OR REPLACE FUNCTION set_updated_at_timestamp_func()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
-- Add triggers to relevant tables
CREATE TRIGGER set_chargeback_updated_at
BEFORE UPDATE ON "chargeback"
FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp_func();

CREATE TRIGGER set_nonipac_updated_at
BEFORE UPDATE ON "nonIpac"
FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp_func();

CREATE TRIGGER set_user_updated_at
BEFORE UPDATE ON "user"
FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp_func();

CREATE TRIGGER set_customer_poc_updated_at
BEFORE UPDATE ON "customer_poc"
FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp_func();

CREATE TRIGGER set_comments_updated_at
BEFORE UPDATE ON "comments"
FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp_func();

CREATE TRIGGER set_agency_bureau_updated_at
BEFORE UPDATE ON "agency_bureau"
FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp_func();
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Drop triggers and function (order matters for triggers, function last)

DROP TRIGGER IF EXISTS set_agency_bureau_updated_at ON "agency_bureau";
DROP TRIGGER IF EXISTS set_comments_updated_at ON "comments";
DROP TRIGGER IF EXISTS set_customer_poc_updated_at ON "customer_poc";
DROP TRIGGER IF EXISTS set_user_updated_at ON "user";
DROP TRIGGER IF EXISTS set_nonipac_updated_at ON "nonIpac";
DROP TRIGGER IF EXISTS set_chargeback_updated_at ON "chargeback";

DROP FUNCTION IF EXISTS set_updated_at_timestamp_func();
-- +goose StatementEnd
