-- Add WeChat Official Account notification bindings and delivery cooldowns.

CREATE TABLE IF NOT EXISTS user_wechat_official_bindings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    openid VARCHAR(128) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    bound_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    unbound_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_wechat_official_bindings_user_active
    ON user_wechat_official_bindings(user_id)
    WHERE unbound_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_wechat_official_bindings_openid_active
    ON user_wechat_official_bindings(openid)
    WHERE unbound_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_wechat_official_bindings_user_id
    ON user_wechat_official_bindings(user_id);

CREATE TABLE IF NOT EXISTS user_notification_events (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel VARCHAR(32) NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    resource_key VARCHAR(128) NOT NULL DEFAULT '',
    last_sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, channel, event_type, resource_key)
);

CREATE INDEX IF NOT EXISTS idx_user_notification_events_user_channel
    ON user_notification_events(user_id, channel, event_type);
