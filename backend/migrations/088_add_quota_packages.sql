-- Add repeat-purchasable quota package support.
-- Quota packages are independent from user_subscriptions and do not affect
-- subscription refresh/window logic.

ALTER TABLE groups
    ADD COLUMN IF NOT EXISTS quota_package_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS quota_package_quota_usd DECIMAL(20,8),
    ADD COLUMN IF NOT EXISTS quota_package_validity_days INTEGER NOT NULL DEFAULT 30;

UPDATE groups
SET quota_package_validity_days = 30
WHERE quota_package_validity_days IS NULL OR quota_package_validity_days <= 0;

CREATE TABLE IF NOT EXISTS user_quota_packages (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    order_id BIGINT REFERENCES payment_orders(id) ON DELETE SET NULL,
    total_quota_usd DECIMAL(20,8) NOT NULL,
    remaining_quota_usd DECIMAL(20,8) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_quota_packages_order_id
    ON user_quota_packages(order_id)
    WHERE order_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_user_quota_packages_available
    ON user_quota_packages(user_id, group_id, status, expires_at)
    WHERE remaining_quota_usd > 0;

CREATE INDEX IF NOT EXISTS idx_user_quota_packages_user_group
    ON user_quota_packages(user_id, group_id);
