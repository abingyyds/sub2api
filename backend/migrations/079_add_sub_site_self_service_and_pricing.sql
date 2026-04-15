ALTER TABLE sub_sites
    ADD COLUMN IF NOT EXISTS parent_sub_site_id BIGINT REFERENCES sub_sites(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS level SMALLINT NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS theme_template VARCHAR(64) NOT NULL DEFAULT 'starter',
    ADD COLUMN IF NOT EXISTS custom_config TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS registration_mode VARCHAR(16) NOT NULL DEFAULT 'open',
    ADD COLUMN IF NOT EXISTS enable_topup BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS allow_sub_site BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS sub_site_price_fen INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS subscription_expired_at TIMESTAMPTZ NULL;

CREATE INDEX IF NOT EXISTS idx_sub_sites_parent_sub_site_id ON sub_sites(parent_sub_site_id);
CREATE INDEX IF NOT EXISTS idx_sub_sites_level ON sub_sites(level);

CREATE TABLE IF NOT EXISTS sub_site_group_prices (
    id BIGSERIAL PRIMARY KEY,
    sub_site_id BIGINT NOT NULL REFERENCES sub_sites(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    price_fen INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(sub_site_id, group_id)
);

CREATE INDEX IF NOT EXISTS idx_sub_site_group_prices_sub_site_id ON sub_site_group_prices(sub_site_id);

CREATE TABLE IF NOT EXISTS sub_site_recharge_prices (
    id BIGSERIAL PRIMARY KEY,
    sub_site_id BIGINT NOT NULL REFERENCES sub_sites(id) ON DELETE CASCADE,
    plan_key VARCHAR(100) NOT NULL,
    pay_amount_fen INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(sub_site_id, plan_key)
);

CREATE INDEX IF NOT EXISTS idx_sub_site_recharge_prices_sub_site_id ON sub_site_recharge_prices(sub_site_id);

CREATE TABLE IF NOT EXISTS sub_site_activation_orders (
    payment_order_id BIGINT PRIMARY KEY REFERENCES payment_orders(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_sub_site_id BIGINT REFERENCES sub_sites(id) ON DELETE SET NULL,
    level SMALLINT NOT NULL DEFAULT 1,
    validity_days INTEGER NOT NULL DEFAULT 365,
    payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    sub_site_id BIGINT REFERENCES sub_sites(id) ON DELETE SET NULL,
    activated_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sub_site_activation_orders_user_id ON sub_site_activation_orders(user_id);
