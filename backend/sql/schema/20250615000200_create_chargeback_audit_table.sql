-- +goose Up
-- Create the audit table for chargeback changes

CREATE TABLE audit.chargeback_changes (
    audit_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_id UUID NOT NULL, -- The ID of the chargeback record being audited
    operation CHAR(1) NOT NULL, -- 'I' (Insert), 'U' (Update), 'D' (Delete)
    changed_by UUID, -- The user who made the change (FK to users.id, can be NULL for automated changes)
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    old_data JSONB, -- Full row snapshot BEFORE the change (for UPDATE/DELETE)
    new_data JSONB   -- Full row snapshot AFTER the change (for INSERT/UPDATE)
);

-- Add index for efficient lookup by chargeback ID
CREATE INDEX idx_audit_chargeback_target_id ON audit.chargeback_changes (target_id);
-- Add index for chronological ordering
CREATE INDEX idx_audit_chargeback_changed_at ON audit.chargeback_changes (changed_at DESC);

-- +goose Down
-- Drop the audit table for chargeback changes

DROP TABLE IF EXISTS audit.chargeback_changes;
