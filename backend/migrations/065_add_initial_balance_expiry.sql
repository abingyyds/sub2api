-- 065: Add initial balance expiry fields to users table
-- Tracks the initial balance given at registration and its expiration time.
-- When initial_balance_expires_at passes, a background job clears the initial balance portion.

ALTER TABLE users ADD COLUMN IF NOT EXISTS initial_balance DECIMAL(20,8) DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS initial_balance_expires_at TIMESTAMP;
