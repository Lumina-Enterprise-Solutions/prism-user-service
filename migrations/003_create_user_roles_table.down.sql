-- Drop indexes
DROP INDEX IF EXISTS idx_user_roles_role_id;
DROP INDEX IF EXISTS idx_user_roles_user_id;

-- Drop table
DROP TABLE IF EXISTS user_roles;
