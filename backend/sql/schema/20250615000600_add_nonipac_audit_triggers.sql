-- +goose Up
-- Add audit triggers to the "nonipac" table

CREATE TRIGGER nonipac_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON "nonipac"
FOR EACH ROW EXECUTE FUNCTION audit.if_modified_func();

-- +goose Down
-- Drop audit triggers from the "nonipac" table

DROP TRIGGER IF EXISTS nonipac_audit_trigger ON "nonipac";
