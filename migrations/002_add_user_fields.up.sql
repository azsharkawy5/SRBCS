-- Add missing fields to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE,
ADD COLUMN IF NOT EXISTS otp VARCHAR(6),
ADD COLUMN IF NOT EXISTS otp_expires_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'user';

-- Add check constraint for valid role values
ALTER TABLE users ADD CONSTRAINT check_users_role 
CHECK (role IN ('admin', 'user'));

-- Add check constraint for OTP length
ALTER TABLE users ADD CONSTRAINT check_users_otp_length 
CHECK (otp IS NULL OR LENGTH(otp) = 6);

-- Create index on role for filtering
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Create index on is_active for filtering active users
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
