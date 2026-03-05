-- 050: Add org_id to api_keys for organization-owned API keys
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS org_id BIGINT REFERENCES organizations(id);
CREATE INDEX IF NOT EXISTS idx_api_keys_org_id ON api_keys(org_id);
