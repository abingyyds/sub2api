-- +migrate Up
CREATE TABLE org_projects (
    id              BIGSERIAL PRIMARY KEY,
    org_id          BIGINT NOT NULL REFERENCES organizations(id),
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    group_id        BIGINT REFERENCES groups(id),
    allowed_models  JSONB,
    monthly_budget_usd DECIMAL(20,8),
    monthly_usage_usd  DECIMAL(20,10) NOT NULL DEFAULT 0,
    monthly_window_start TIMESTAMPTZ,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_org_projects_org_name_active ON org_projects(org_id, name) WHERE deleted_at IS NULL;
CREATE INDEX idx_org_projects_org_id ON org_projects(org_id);
CREATE INDEX idx_org_projects_deleted_at ON org_projects(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS org_projects;
