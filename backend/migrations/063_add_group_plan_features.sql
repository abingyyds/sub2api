-- Add plan_features column to groups table
ALTER TABLE groups ADD COLUMN IF NOT EXISTS plan_features JSONB DEFAULT '[]'::jsonb;
