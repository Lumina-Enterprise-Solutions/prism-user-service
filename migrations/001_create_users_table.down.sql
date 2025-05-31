-- Drop trigger and function
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_email;

-- Drop table
DROP TABLE IF EXISTS users;
