-- +goose Up
-- Seed initial administrative user for CDMS application

-- Declare a variable to hold the ID of the inserted user. 
With new_user AS (
  INSERT INTO "cdms_user" (first_name, last_name, org, email, is_admin)
  VALUES ('John', 'Willett', 'GSA', 'jjckrbbt@gmail.com', TRUE)
  RETURNING id
)
-- Grant new user acccess to all business lines
INSERT INTO "user_business_line_access" (user_id, business_line)
SELECT
    id, 
    unnest(enum_range(NULL::chargeback_business_line))
FROM new_user;

-- +goose Down
-- Remove user and their associated rights
DELETE FROM "user_business_line_access"
WHERE user_id = (SELECT id FROM "cdms_user" WHERE email = 'jjckrbbt@gmail.com');

DELETE FROM "cdms_user"
WHERE email = 'jjckrbbt@gmail.com';
