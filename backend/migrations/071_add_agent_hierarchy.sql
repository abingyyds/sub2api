-- 071: Add per-user commission rate for agent hierarchy (differential commission)
ALTER TABLE referrals ADD COLUMN IF NOT EXISTS commission_rate DECIMAL(5,4);
