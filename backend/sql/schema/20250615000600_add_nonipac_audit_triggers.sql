-- +goose Up
-- Add audit triggers to the "nonIpac" table

CREATE TRIGGER nonipac_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON "nonIpac"
FOR EACH ROW EXECUTE FUNCTION audit.if_modified_func();

-- +goose Down
-- Drop audit triggers from the "nonIpac" table

DROP TRIGGER IF EXISTS nonipac_audit_trigger ON "nonIpac";
