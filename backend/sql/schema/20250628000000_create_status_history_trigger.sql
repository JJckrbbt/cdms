-- +goose Up

-- Insert the system user if it doesn't already exist.
INSERT INTO "cdms_user" (first_name, last_name, org, email, is_admin, auth_provider_subject)
SELECT 'System', 'User', 'GSA', 'system@cdms.local', TRUE, 'system|internal'
WHERE NOT EXISTS (SELECT 1 FROM "cdms_user" WHERE email = 'system@cdms.local');

-- +goose StatementBegin
-- Re-create the function to include all logic, removing the need for a WHEN clause on the trigger.
CREATE OR REPLACE FUNCTION log_status_history_func()
RETURNS TRIGGER AS $$
DECLARE
    last_status_history_id BIGINT;
    current_user_id BIGINT;
    system_user_id BIGINT;
BEGIN
    -- The logic from the WHEN clause is now inside the function.
    -- The function will do nothing if the conditions aren't met.
    IF (TG_OP = 'UPDATE' AND OLD.current_status IS NOT DISTINCT FROM NEW.current_status) THEN
        RETURN NEW; -- Exit if status has not changed on an UPDATE
    END IF;

    -- At this point, we know it's an INSERT or a status-changing UPDATE.

    -- Find the ID of our stable system user.
    SELECT id INTO system_user_id FROM "cdms_user" WHERE email = 'system@cdms.local';

    -- Try to get the user_id from the current session's settings.
    BEGIN
        current_user_id := current_setting('app.user_id')::BIGINT;
    EXCEPTION WHEN OTHERS THEN
        current_user_id := system_user_id;
    END;

    -- Insert the new status into the history table.
    INSERT INTO "status_history" (status, notes, user_id, status_date)
    VALUES (NEW.current_status, 'Status ' || NEW.current_status || ' logged via trigger for ' || TG_TABLE_NAME, current_user_id, NOW())
    RETURNING id INTO last_status_history_id;

    -- Link the new history record.
    IF TG_TABLE_NAME = 'chargeback' THEN
        INSERT INTO "chargeback_status_merge" (chargeback_id, status_history_id)
        VALUES (NEW.id, last_status_history_id);
    ELSIF TG_TABLE_NAME = 'nonipac' THEN
        INSERT INTO "nonipac_status_merge" (nonipac_id, status_history_id)
        VALUES (NEW.id, last_status_history_id);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Create the trigger for chargeback, now without the complex WHEN clause.
CREATE TRIGGER log_chargeback_status_change
AFTER INSERT OR UPDATE ON "chargeback"
FOR EACH ROW
EXECUTE FUNCTION log_status_history_func();

-- Create the trigger for nonipac, also without the WHEN clause.
CREATE TRIGGER log_nonipac_status_change
AFTER INSERT OR UPDATE ON "nonipac"
FOR EACH ROW
EXECUTE FUNCTION log_status_history_func();


-- +goose Down
DROP TRIGGER IF EXISTS log_chargeback_status_change ON "chargeback";
DROP TRIGGER IF EXISTS log_nonipac_status_change ON "nonipac";
DROP FUNCTION IF EXISTS log_status_history_func();
