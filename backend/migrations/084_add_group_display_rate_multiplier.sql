ALTER TABLE groups
  ADD COLUMN IF NOT EXISTS display_rate_multiplier DECIMAL(10,4);
