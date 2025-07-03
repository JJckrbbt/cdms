-- +goose Up
-- Create all custom ENUM types used in the schema

CREATE TYPE cdms_status AS ENUM (
  'Open',
  'Hold Pending External Action',
  'Hold Pending Internal Action',
  'In Research',
  'Passed to PFS',
  'Completed by PFS',
  'PFS Return to GSA',
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
  'EIS Issues'
);

CREATE TYPE chargeback_reporting_source AS ENUM (
  'BC1048',
  'BC1300',
  'ApplicationCreated'
);



CREATE TYPE chargeback_reason_code AS ENUM (
  'Incorrect ALC',
  'Incorrect TAS',
  'Incorrect LOA',
  'Need Supporting Documentation',
  'No Funds Available',
  'Billed Wrong Amount',
  'Bill Exceeds Authorized Amount',
  'Funds Expired',
  'Billed Goods or Services Unsatisfactory/Not Received',
  'Missing Customer Order Number',
  'Billed Incorrect Method',
  'PO Canceled or Ended',
  'End of Month Rejection',
  'No or Incorrect FSN',
  'Customer Billing Office Closure/Reorg',
  'Other/Multiple',
  'Funds Not Obligated by Client',
  'Speedpay not updated',
  'COVID-19 Agency Pickup Delay',
  'Mileage Billing Errors',
  'BETC Update Needed',
  'EIS Issues',
  'Wrong PC Code',
  'Wallet not updated'
);

CREATE TYPE chargeback_action AS ENUM (
  'Rebill',
  'Reverse to Income',
  'Reverse to Income & GSA Rebill',
  'Write Off',
  'Return to Treasury',
  'Other - See Special Instructions'
);

CREATE TYPE chargeback_business_line AS ENUM (
  'Procurement',
  'Operations',
  'Research & Dev',
  'IT Services',
  'Logistics',
  'Admin',
  'Cars',
  'Rent',
  'Credit',
  'Hotels',
  'Grocery'
);

CREATE TYPE chargeback_fund AS ENUM (
  'F-100',
  'F-201',
  'F-305',
  'F-410',
  'F-501'
);

CREATE TYPE nonipac_reporting_source AS ENUM (
  'ApplicationCreated',
  'OUTSTANDING_BILLS'
);

CREATE TYPE status_history_status AS ENUM (
  'Open',
  'Hold Pending External Action',
  'Hold Pending Internal Action',
  'In Research',
  'Passed to PFS',
  'Completed by PFS',
  'PFS Return to GSA',
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
  'EIS Issues'
);

CREATE TYPE user_org AS ENUM (
  'GSA',
  'PFS'
);

-- +goose Down
-- Drop all custom ENUM types

DROP TYPE IF EXISTS user_org;
DROP TYPE IF EXISTS status_history_status;
DROP TYPE IF EXISTS nonipac_reporting_source;
DROP TYPE IF EXISTS chargeback_fund;
DROP TYPE IF EXISTS chargeback_business_line;
DROP TYPE IF EXISTS chargeback_action;
DROP TYPE IF EXISTS chargeback_reason_code;
DROP TYPE IF EXISTS chargeback_reporting_source;
DROP TYPE IF EXISTS cdms_status;
