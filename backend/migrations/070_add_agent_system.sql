-- 070: 代理系统 (Agent/Affiliate System)
-- 在 users 表增加代理相关字段，创建佣金记录表

-- 用户表增加代理字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_agent BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS agent_status VARCHAR(20) NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS agent_commission_rate DECIMAL(5,4) NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS agent_note TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS agent_approved_at TIMESTAMPTZ;

-- 代理佣金记录表
CREATE TABLE IF NOT EXISTS agent_commissions (
    id BIGSERIAL PRIMARY KEY,
    agent_id BIGINT NOT NULL REFERENCES users(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    order_id BIGINT REFERENCES payment_orders(id),
    source_type VARCHAR(20) NOT NULL,
    source_amount DECIMAL(20,8) NOT NULL DEFAULT 0,
    commission_rate DECIMAL(5,4) NOT NULL DEFAULT 0,
    commission_amount DECIMAL(20,8) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    settled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agent_commissions_agent_id ON agent_commissions(agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_commissions_user_id ON agent_commissions(user_id);
CREATE INDEX IF NOT EXISTS idx_agent_commissions_order_id ON agent_commissions(order_id);
CREATE INDEX IF NOT EXISTS idx_agent_commissions_status ON agent_commissions(status);
CREATE INDEX IF NOT EXISTS idx_agent_commissions_created_at ON agent_commissions(created_at);

-- users 表索引
CREATE INDEX IF NOT EXISTS idx_users_is_agent ON users(is_agent) WHERE is_agent = true;
CREATE INDEX IF NOT EXISTS idx_users_agent_status ON users(agent_status) WHERE agent_status != '';
