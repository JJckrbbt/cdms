-- name: GetUserWithAuthorizationContext :one
-- Fetches a single user by their ID, including their roles, permissions, and business lines.
SELECT
    u.id,
    u.email,
    u.first_name,
    u.last_name,
    u.org,
    u.is_active,
    u.created_at,
    COALESCE((SELECT array_agg(r.name) FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = u.id), '{}') AS roles,
    COALESCE((SELECT array_agg(p.action) FROM permissions p JOIN role_permissions rp ON p.id = rp.permission_id JOIN user_roles ur ON rp.role_id = ur.role_id WHERE ur.user_id = u.id), '{}') AS permissions,
    COALESCE((SELECT array_agg(ubla.business_line) FROM user_business_line_access ubla WHERE ubla.user_id = u.id), '{}') AS business_lines
FROM
    "cdms_user" u
WHERE
    u.id = $1;

-- name: ListUsersByBusinessLines :many
-- Fetches a paginated list of users who are associated with a given set of business lines.
-- This is for scoped admins.
SELECT
    u.id,
    u.email,
    u.first_name,
    u.last_name,
    u.org,
    u.is_active,
    u.created_at,
    COALESCE((SELECT array_agg(r.name) FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = u.id), '{}') AS roles
FROM
    "cdms_user" u
WHERE u.id IN (
    SELECT DISTINCT ubla.user_id
    FROM user_business_line_access ubla
    WHERE ubla.business_line = ANY(@business_lines::chargeback_business_line[])
)
ORDER BY
    u.last_name, u.first_name
LIMIT $1
OFFSET $2;

-- name: ListAllUsers :many
-- Fetches a paginated list of all users. For super_admins and global admins.
SELECT
    u.id,
    u.email,
    u.first_name,
    u.last_name,
    u.org,
    u.is_active,
    u.created_at,
    COALESCE((SELECT array_agg(r.name) FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = u.id), '{}') AS roles
FROM
    "cdms_user" u
ORDER BY
    u.last_name, u.first_name
LIMIT $1
OFFSET $2;


-- name: UpdateUser :one
-- Updates a user's mutable details.
UPDATE "cdms_user"
SET
    first_name = $2,
    last_name = $3,
    org = $4,
    is_active = $5
WHERE
    id = $1
RETURNING *;

-- name: ListRoles :many
-- Fetches all available roles in the system.
SELECT id, name, description FROM "roles" ORDER BY id;

-- name: AssignRoleToUser :exec
-- Assigns a specific role to a user.
INSERT INTO "user_roles" (user_id, role_id) VALUES ($1, $2)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- name: RemoveRoleFromUser :exec
-- Removes a specific role from a user.
DELETE FROM "user_roles" WHERE user_id = $1 AND role_id = $2;

-- name: AssignBusinessLinesToUser :exec
-- Assigns a set of business lines to a user, replacing existing ones.
-- This uses a CTE to first delete old assignments, then insert new ones.
WITH deleted AS (
    DELETE FROM "user_business_line_access"
    WHERE user_id = $1
)
INSERT INTO "user_business_line_access" (user_id, business_line)
SELECT $1, unnest(@business_lines::chargeback_business_line[]);

-- name: RemoveAllRolesFromUser :exec
-- Removes all roles from a user.
DELETE FROM "user_roles" WHERE user_id = $1;
