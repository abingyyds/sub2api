-- 管理员邀请码表
CREATE TABLE IF NOT EXISTS admin_invite_codes (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    source_name VARCHAR(100) NOT NULL,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    used_count INT NOT NULL DEFAULT 0,
    max_uses INT,
    enabled BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_admin_invite_codes_code ON admin_invite_codes(code);
CREATE INDEX idx_admin_invite_codes_created_by ON admin_invite_codes(created_by);
CREATE INDEX idx_admin_invite_codes_enabled ON admin_invite_codes(enabled);
