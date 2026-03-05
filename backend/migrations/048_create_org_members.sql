-- 048: Create org_members table for organization membership
CREATE TABLE IF NOT EXISTS org_members (
    id              BIGSERIAL PRIMARY KEY,
    org_id          BIGINT NOT NULL REFERENCES organizations(id),
    user_id         BIGINT NOT NULL REFERENCES users(id),
    role            VARCHAR(20) NOT NULL DEFAULT 'member',
    monthly_quota_usd   DECIMAL(20,8),
    daily_quota_usd     DECIMAL(20,8),
    monthly_usage_usd   DECIMAL(20,10) NOT NULL DEFAULT 0,
    daily_usage_usd     DECIMAL(20,10) NOT NULL DEFAULT 0,
    monthly_window_start TIMESTAMPTZ,
    daily_window_start   TIMESTAMPTZ,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_org_members_org_user_active ON org_members(org_id, user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_org_members_org_id ON org_members(org_id);
CREATE INDEX IF NOT EXISTS idx_org_members_user_id ON org_members(user_id);
CREATE INDEX IF NOT EXISTS idx_org_members_deleted_at ON org_members(deleted_at);
