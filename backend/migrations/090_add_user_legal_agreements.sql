CREATE TABLE IF NOT EXISTS user_legal_agreements (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    terms_accepted_at TIMESTAMPTZ NOT NULL,
    privacy_accepted_at TIMESTAMPTZ NOT NULL,
    api_terms_accepted_at TIMESTAMPTZ NOT NULL,
    terms_version VARCHAR(32) NOT NULL DEFAULT '',
    privacy_version VARCHAR(32) NOT NULL DEFAULT '',
    api_terms_version VARCHAR(32) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_legal_agreements_updated_at
    ON user_legal_agreements(updated_at DESC);
