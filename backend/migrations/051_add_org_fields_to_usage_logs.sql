-- 051: Add org_id and org_member_id to usage_logs for organization usage tracking
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS org_id BIGINT;
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS org_member_id BIGINT;
CREATE INDEX IF NOT EXISTS idx_usage_logs_org_id ON usage_logs(org_id);
CREATE INDEX IF NOT EXISTS idx_usage_logs_org_member_id ON usage_logs(org_member_id);
