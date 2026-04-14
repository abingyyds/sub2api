CREATE TABLE IF NOT EXISTS sub_sites (
    id BIGSERIAL PRIMARY KEY,
    owner_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(120) NOT NULL,
    slug VARCHAR(64) NOT NULL UNIQUE,
    custom_domain VARCHAR(255) UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    site_logo TEXT NOT NULL DEFAULT '',
    site_favicon TEXT NOT NULL DEFAULT '',
    site_subtitle TEXT NOT NULL DEFAULT '',
    announcement TEXT NOT NULL DEFAULT '',
    contact_info TEXT NOT NULL DEFAULT '',
    doc_url TEXT NOT NULL DEFAULT '',
    home_content TEXT NOT NULL DEFAULT '',
    theme_config TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sub_sites_owner_user_id ON sub_sites(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_sub_sites_status ON sub_sites(status);
CREATE INDEX IF NOT EXISTS idx_sub_sites_custom_domain ON sub_sites(custom_domain);

CREATE TABLE IF NOT EXISTS sub_site_users (
    id BIGSERIAL PRIMARY KEY,
    sub_site_id BIGINT NOT NULL REFERENCES sub_sites(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source VARCHAR(32) NOT NULL DEFAULT 'register',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);

CREATE INDEX IF NOT EXISTS idx_sub_site_users_sub_site_id ON sub_site_users(sub_site_id);
CREATE INDEX IF NOT EXISTS idx_sub_site_users_created_at ON sub_site_users(created_at DESC);
