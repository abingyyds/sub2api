-- +migrate Up
CREATE TABLE org_audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    org_id          BIGINT NOT NULL,
    user_id         BIGINT NOT NULL,
    member_id       BIGINT,
    project_id      BIGINT,
    usage_log_id    BIGINT,
    action          VARCHAR(50) NOT NULL,
    model           VARCHAR(100),
    audit_mode      VARCHAR(20) NOT NULL DEFAULT 'metadata',
    request_summary TEXT,
    request_content TEXT,
    response_summary TEXT,
    keywords        JSONB,
    flagged         BOOLEAN NOT NULL DEFAULT false,
    flag_reason     TEXT,
    input_tokens    INT,
    output_tokens   INT,
    cost_usd        DECIMAL(20,10),
    ip_address      VARCHAR(45),
    user_agent      VARCHAR(512),
    detail          JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_org_audit_logs_org_id ON org_audit_logs(org_id);
CREATE INDEX idx_org_audit_logs_user_id ON org_audit_logs(user_id);
CREATE INDEX idx_org_audit_logs_project_id ON org_audit_logs(project_id);
CREATE INDEX idx_org_audit_logs_created_at ON org_audit_logs(created_at);
CREATE INDEX idx_org_audit_logs_flagged ON org_audit_logs(org_id, flagged) WHERE flagged = true;
CREATE INDEX idx_org_audit_logs_action ON org_audit_logs(action);

-- +migrate Down
DROP TABLE IF EXISTS org_audit_logs;
