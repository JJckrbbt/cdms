-- +goose Up
-- Add audit triggers to the "cdms_user" table

CREATE TRIGGER user_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON "cdms_user"
FOR EACH ROW EXECUTE FUNCTION audit.if_modified_func();

-- +goose Down
-- Drop audit triggers from the "cdms_user" table

DROP TRIGGER IF EXISTS user_audit_trigger ON "cdms_user";
