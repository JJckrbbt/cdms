-- +goose Up
-- Add audit triggers to the "chargeback" table

CREATE TRIGGER chargeback_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON "chargeback"
FOR EACH ROW EXECUTE FUNCTION audit.if_modified_func();

-- +goose Down
-- Drop audit triggers from the "chargeback" table

DROP TRIGGER IF EXISTS chargeback_audit_trigger ON "chargeback";
