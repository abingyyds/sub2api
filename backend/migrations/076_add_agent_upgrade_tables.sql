-- 076: agent upgrade tables for simplified onboarding, wallet and withdraw workflow

CREATE TABLE IF NOT EXISTS agent_profiles (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    real_name VARCHAR(100) NOT NULL DEFAULT '',
    id_card_no VARCHAR(64) NOT NULL DEFAULT '',
    phone VARCHAR(32) NOT NULL DEFAULT '',
    identity_status VARCHAR(20) NOT NULL DEFAULT 'unsubmitted',
    identity_submitted_at TIMESTAMPTZ,
    contract_status VARCHAR(20) NOT NULL DEFAULT 'unsigned',
    contract_version VARCHAR(32) NOT NULL DEFAULT 'v1',
    contract_signed_at TIMESTAMPTZ,
    contract_ip VARCHAR(64) NOT NULL DEFAULT '',
    activation_order_id BIGINT REFERENCES payment_orders(id),
    activation_fee_paid_at TIMESTAMPTZ,
    frozen_balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    withdrawable_balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    total_withdrawn DECIMAL(20,8) NOT NULL DEFAULT 0,
    is_frozen BOOLEAN NOT NULL DEFAULT false,
    frozen_reason TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agent_profiles_identity_status ON agent_profiles(identity_status);
CREATE INDEX IF NOT EXISTS idx_agent_profiles_contract_status ON agent_profiles(contract_status);
CREATE INDEX IF NOT EXISTS idx_agent_profiles_paid_at ON agent_profiles(activation_fee_paid_at);
CREATE INDEX IF NOT EXISTS idx_agent_profiles_is_frozen ON agent_profiles(is_frozen);

CREATE TABLE IF NOT EXISTS agent_wallet_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance_type VARCHAR(20) NOT NULL,
    change_type VARCHAR(40) NOT NULL,
    amount DECIMAL(20,8) NOT NULL,
    related_user_id BIGINT REFERENCES users(id),
    related_order_id BIGINT REFERENCES payment_orders(id),
    withdraw_request_id BIGINT,
    unlock_at TIMESTAMPTZ,
    remark TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agent_wallet_logs_user_id ON agent_wallet_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_agent_wallet_logs_created_at ON agent_wallet_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_agent_wallet_logs_unlock_at ON agent_wallet_logs(unlock_at);

CREATE TABLE IF NOT EXISTS agent_withdraw_requests (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(20,8) NOT NULL,
    alipay_name VARCHAR(100) NOT NULL DEFAULT '',
    alipay_account VARCHAR(255) NOT NULL DEFAULT '',
    alipay_qr_image TEXT NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    review_note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_agent_withdraw_requests_user_id ON agent_withdraw_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_agent_withdraw_requests_status ON agent_withdraw_requests(status);
