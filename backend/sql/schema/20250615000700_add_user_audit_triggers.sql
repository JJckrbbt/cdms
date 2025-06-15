-- +goose Up
-- Add audit triggers to the "user" table

CREATE TRIGGER user_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON "user"
FOR EACH ROW EXECUTE FUNCTION audit.if_modified_func();

-- +goose Down
-- Drop audit triggers from the "user" table

DROP TRIGGER IF EXISTS user_audit_trigger ON "user";
