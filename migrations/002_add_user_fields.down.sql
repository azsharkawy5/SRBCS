-- Remove added fields from users table
ALTER TABLE users 
DROP COLUMN IF EXISTS role,
DROP COLUMN IF EXISTS otp_expires_at,
DROP COLUMN IF EXISTS otp,
DROP COLUMN IF EXISTS is_active,
DROP COLUMN IF EXISTS is_email_verified;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_role;
