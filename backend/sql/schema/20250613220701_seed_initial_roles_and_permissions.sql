-- +goose Up
-- Seed the RBAC tables with the final, granular, scoped role structure

-- 1. Create the roles
INSERT INTO "roles" (name, description) VALUES
('super_admin', 'Has all permissions, including managing other users and assigning admin roles.'),
('admin', 'Can manage data, uploads, and assign non-admin roles to any user.'),
('business_line_admin', 'Can assign analyst and viewer roles to users within their own business line(s).'),
('maintainer', 'Can upload reports and edit chargeback/delinquency data.'),
('analyst', 'Can view and edit chargeback and delinquency records.'),
('viewer', 'Has read-only access to dashboards and data.');

-- 2. Create the granular permissions
INSERT INTO "permissions" (action, description) VALUES
('roles:manage_admins', 'Ability to assign the admin or super_admin roles.'),
('roles:assign_global', 'Ability to assign non-admin roles to any user.'),
('roles:assign_scoped', 'Ability to assign non-admin roles to users within the same business line.'),
('users:edit', 'Ability to edit a user''s details (e.g., active status).'),
('users:view_scoped', 'Ability to view users''s within the same business line(s).'),
('reports:upload', 'Ability to upload new reports.'),
('chargebacks:edit', 'Ability to edit chargeback records.'),
('delinquencies:edit', 'Ability to edit delinquency records.'),
('data:view', 'Ability to view all chargeback, delinquency, and dashboard data.');

-- 3. Assign permissions to roles
-- Super Admin gets everything
INSERT INTO "role_permissions" (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'super_admin'), p.id FROM permissions p;

-- Admin gets everything EXCEPT managing other admins
INSERT INTO "role_permissions" (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'admin'), p.id FROM permissions p WHERE p.action != 'roles:manage_admins';

-- Business Line Admin gets scoped user management and view rights
INSERT INTO "role_permissions" (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'business_line_admin'), p.id FROM permissions p WHERE p.action IN ('roles:assign_scoped', 'users:edit', 'users:view_scoped', 'data:view');

-- Maintainer can upload, edit, and view
INSERT INTO "role_permissions" (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'maintainer'), p.id FROM permissions p WHERE p.action IN ('reports:upload', 'chargebacks:edit', 'delinquencies:edit', 'data:view');

-- Analyst can edit and view
INSERT INTO "role_permissions" (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'analyst'), p.id FROM permissions p WHERE p.action IN ('chargebacks:edit', 'delinquencies:edit', 'data:view');

-- Viewer can only view
INSERT INTO "role_permissions" (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'viewer'), p.id FROM permissions p WHERE p.action = 'data:view';


-- +goose Down
-- Clear out all the seeded data
DELETE FROM "role_permissions";
DELETE FROM "permissions";
DELETE FROM "roles";
