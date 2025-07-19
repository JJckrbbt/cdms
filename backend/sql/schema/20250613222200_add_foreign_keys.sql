-- +goose Up
-- Add all foreign key constraints

-- chargeback foreign keys
ALTER TABLE "chargeback" ADD FOREIGN KEY ("vendor") REFERENCES "agency_bureau" ("vendor_code");

-- nonipac foreign keys
ALTER TABLE "nonipac" ADD FOREIGN KEY ("address_code") REFERENCES "agency_bureau" ("vendor_code");
ALTER TABLE "nonipac" ADD FOREIGN KEY ("pfs_poc") REFERENCES "cdms_user" ("id");
ALTER TABLE "nonipac" ADD FOREIGN KEY ("gsa_poc") REFERENCES "cdms_user" ("id");
ALTER TABLE "nonipac" ADD FOREIGN KEY ("customer_poc") REFERENCES "customer_poc" ("id");

-- status_history foreign keys
ALTER TABLE "status_history" ADD FOREIGN KEY ("user_id") REFERENCES "cdms_user" ("id");

-- chargeback_status_merge foreign keys
ALTER TABLE "chargeback_status_merge" ADD FOREIGN KEY ("chargeback_id") REFERENCES "chargeback" ("id");
ALTER TABLE "chargeback_status_merge" ADD FOREIGN KEY ("status_history_id") REFERENCES "status_history" ("id");

-- nonipac_status_merge foreign keys
ALTER TABLE "nonipac_status_merge" ADD FOREIGN KEY ("nonipac_id") REFERENCES "nonipac" ("id");
ALTER TABLE "nonipac_status_merge" ADD FOREIGN KEY ("status_history_id") REFERENCES "status_history" ("id");

-- issue_owner_gsa_chargeback_merge foreign keys
ALTER TABLE "issue_owner_gsa_chargeback_merge" ADD FOREIGN KEY ("user_id") REFERENCES "cdms_user" ("id");
ALTER TABLE "issue_owner_gsa_chargeback_merge" ADD FOREIGN KEY ("chargeback_id") REFERENCES "chargeback" ("id");

-- issue_owner_pfs_chargeback_merge foreign keys
ALTER TABLE "issue_owner_pfs_chargeback_merge" ADD FOREIGN KEY ("user_id") REFERENCES "cdms_user" ("id");
ALTER TABLE "issue_owner_pfs_chargeback_merge" ADD FOREIGN KEY ("chargeback_id") REFERENCES "chargeback" ("id");

-- non_ipac_customer_poc_merge foreign keys
ALTER TABLE "non_ipac_customer_poc_merge" ADD FOREIGN KEY ("nonipac_id") REFERENCES "nonipac" ("id");
ALTER TABLE "non_ipac_customer_poc_merge" ADD FOREIGN KEY ("customer_poc_id") REFERENCES "customer_poc" ("id");

-- chargeback_customer_poc_merge foreign keys
ALTER TABLE "chargeback_customer_poc_merge" ADD FOREIGN KEY ("chargeback_id") REFERENCES "chargeback" ("id");
ALTER TABLE "chargeback_customer_poc_merge" ADD FOREIGN KEY ("customer_poc_id") REFERENCES "customer_poc" ("id");

-- comments foreign keys
ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "cdms_user" ("id");

-- chargeback_comments_merge foreign keys
ALTER TABLE "chargeback_comments_merge" ADD FOREIGN KEY ("chargeback_id") REFERENCES "chargeback" ("id");
ALTER TABLE "chargeback_comments_merge" ADD FOREIGN KEY ("comment_id") REFERENCES "comments" ("id");

-- non_ipac_comments_merge foreign keys
ALTER TABLE "non_ipac_comments_merge" ADD FOREIGN KEY ("nonipac_id") REFERENCES "nonipac" ("id");
ALTER TABLE "non_ipac_comments_merge" ADD FOREIGN KEY ("comment_id") REFERENCES "comments" ("id");

-- comment_mentions foreign keys
ALTER TABLE "comment_mentions" ADD FOREIGN KEY ("comment_id") REFERENCES "comments" ("id");
ALTER TABLE "comment_mentions" ADD FOREIGN KEY ("user_id") REFERENCES "cdms_user" ("id");


-- +goose Down
-- Drop all foreign key constraints (in reverse order of creation if dependencies exist)

ALTER TABLE "comment_mentions" DROP CONSTRAINT IF EXISTS "comment_mentions_user_id_fkey";
ALTER TABLE "comment_mentions" DROP CONSTRAINT IF EXISTS "comment_mentions_comment_id_fkey";

ALTER TABLE "non_ipac_comments_merge" DROP CONSTRAINT IF EXISTS "non_ipac_comments_merge_comment_id_fkey";
ALTER TABLE "non_ipac_comments_merge" DROP CONSTRAINT IF EXISTS "non_ipac_comments_merge_nonipac_id_fkey";

ALTER TABLE "chargeback_comments_merge" DROP CONSTRAINT IF EXISTS "chargeback_comments_merge_comment_id_fkey";
ALTER TABLE "chargeback_comments_merge" DROP CONSTRAINT IF EXISTS "chargeback_comments_merge_chargeback_id_fkey";

ALTER TABLE "comments" DROP CONSTRAINT IF EXISTS "comments_user_id_fkey";

ALTER TABLE "chargeback_customer_poc_merge" DROP CONSTRAINT IF EXISTS "chargeback_customer_poc_merge_customer_poc_id_fkey";
ALTER TABLE "chargeback_customer_poc_merge" DROP CONSTRAINT IF EXISTS "chargeback_customer_poc_merge_chargeback_id_fkey";

ALTER TABLE "non_ipac_customer_poc_merge" DROP CONSTRAINT IF EXISTS "non_ipac_customer_poc_merge_customer_poc_id_fkey";
ALTER TABLE "non_ipac_customer_poc_merge" DROP CONSTRAINT IF EXISTS "non_ipac_customer_poc_merge_nonipac_id_fkey";

ALTER TABLE "issue_owner_pfs_chargeback_merge" DROP CONSTRAINT IF EXISTS "issue_owner_pfs_chargeback_merge_chargeback_id_fkey";
ALTER TABLE "issue_owner_pfs_chargeback_merge" DROP CONSTRAINT IF EXISTS "issue_owner_pfs_chargeback_merge_user_id_fkey";

ALTER TABLE "issue_owner_gsa_chargeback_merge" DROP CONSTRAINT IF EXISTS "issue_owner_gsa_chargeback_merge_chargeback_id_fkey";
ALTER TABLE "issue_owner_gsa_chargeback_merge" DROP CONSTRAINT IF EXISTS "issue_owner_gsa_chargeback_merge_user_id_fkey";

ALTER TABLE "nonipac_status_merge" DROP CONSTRAINT IF EXISTS "nonipac_status_merge_status_history_id_fkey";
ALTER TABLE "nonipac_status_merge" DROP CONSTRAINT IF EXISTS "nonipac_status_merge_nonipac_id_fkey";

ALTER TABLE "chargeback_status_merge" DROP CONSTRAINT IF EXISTS "chargeback_status_merge_status_history_id_fkey";
ALTER TABLE "chargeback_status_merge" DROP CONSTRAINT IF EXISTS "chargeback_status_merge_chargeback_id_fkey";

ALTER TABLE "status_history" DROP CONSTRAINT IF EXISTS "status_history_user_id_fkey";

ALTER TABLE "nonipac" DROP CONSTRAINT IF EXISTS "nonipac_customer_poc_fkey";
ALTER TABLE "nonipac" DROP CONSTRAINT IF EXISTS "nonipac_gsa_poc_fkey";
ALTER TABLE "nonipac" DROP CONSTRAINT IF EXISTS "nonipac_pfs_poc_fkey";
ALTER TABLE "nonipac" DROP CONSTRAINT IF EXISTS "nonipac_address_code_fkey";

ALTER TABLE "chargeback" DROP CONSTRAINT IF EXISTS "chargeback_vendor_fkey";
