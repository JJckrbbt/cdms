-- +goose Up
-- Create the audit table for nonipac changes

CREATE TABLE audit.nonipac_changes (
    audit_id BIGSERIAL PRIMARY KEY,
    target_id BIGINT NOT NULL, -- The ID of the nonipac record being audited
    operation CHAR(1) NOT NULL, -- 'I' (Insert), 'U' (Update), 'D' (Delete)
    changed_by BIGINT, -- The user who made the change (FK to users.id, can be NULL for automated changes)
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    old_data JSONB, -- Full row snapshot BEFORE the change (for UPDATE/DELETE)
    new_data JSONB   -- Full row snapshot AFTER the change (for INSERT/UPDATE)
);

-- Add index for efficient lookup by nonipac ID
CREATE INDEX idx_audit_nonipac_target_id ON audit.nonipac_changes (target_id);
-- Add index for chronological ordering
CREATE INDEX idx_audit_nonipac_changed_at ON audit.nonipac_changes (changed_at DESC);

-- +goose Down
-- Drop the audit table for nonipac changes

DROP TABLE IF EXISTS audit.nonipac_changes;
