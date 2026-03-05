-- 047: Create organizations table for multi-tenant enterprise support
CREATE TABLE IF NOT EXISTS organizations (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(200) NOT NULL,
    slug            VARCHAR(100) NOT NULL,
    description     TEXT,
    owner_user_id   BIGINT NOT NULL REFERENCES users(id),
    billing_mode    VARCHAR(20) NOT NULL DEFAULT 'balance',
    balance         DECIMAL(20,8) NOT NULL DEFAULT 0,
    monthly_budget_usd  DECIMAL(20,8),
    max_members     INT NOT NULL DEFAULT 50,
    max_api_keys    INT NOT NULL DEFAULT 100,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    audit_mode      VARCHAR(20) NOT NULL DEFAULT 'metadata',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_organizations_slug_active ON organizations(slug) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_organizations_owner_active ON organizations(owner_user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_organizations_status ON organizations(status);
CREATE INDEX IF NOT EXISTS idx_organizations_deleted_at ON organizations(deleted_at);
