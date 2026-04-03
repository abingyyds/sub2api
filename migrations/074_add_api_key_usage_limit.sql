-- Add usage_limit field to api_keys table
-- Usage limit in USD, null means unlimited
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS usage_limit DOUBLE PRECISION;
