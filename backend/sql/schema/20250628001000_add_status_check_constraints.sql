-- +goose Up
-- Add a constraint to the chargeback table to only allow valid statuses for chargebacks.
ALTER TABLE "chargeback" ADD CONSTRAINT "check_chargeback_status"
CHECK (current_status IN (
    'Open',
    'Hold Pending External Action',
    'Hold Pending Internal Action',
    'In Research',
    'Passed to PFS',
    'Completed by PFS',
    'PFS Return to GSA',
    'Reconciled - Off Report'
));

-- Add a constraint to the nonipac table to only allow valid statuses for non-ipac items.
ALTER TABLE "nonipac" ADD CONSTRAINT "check_nonipac_status"
CHECK (current_status IN (
    'Open',
    'Refund',
    'Offset',
    'In Process',
    'Write Off',
    'Referred to Treasury for Collections',
    'Return Credit to Treasury',
    'Waiting on Customer Response',
    'Waiting on GSA Response Pending Payment',
    'Closed - Payment Received',
    'Reverse to Income',
    'Bill as IPAC',
    'Bill as DoD',
    'EIS Issues',
    'Reconciled - Off Report'
));


-- +goose Down
ALTER TABLE "chargeback" DROP CONSTRAINT IF EXISTS "check_chargeback_status";
ALTER TABLE "nonipac" DROP CONSTRAINT IF EXISTS "check_nonipac_status";
