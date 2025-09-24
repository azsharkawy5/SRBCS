-- Drop indexes
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_email;

-- Drop users table
DROP TABLE IF EXISTS users;

-- Drop uuid-ossp extension (optional, as it might be used by other tables)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
