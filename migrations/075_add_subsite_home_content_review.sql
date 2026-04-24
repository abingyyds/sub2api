ALTER TABLE sub_sites
  ADD COLUMN IF NOT EXISTS pending_home_content TEXT,
  ADD COLUMN IF NOT EXISTS home_content_review_status VARCHAR(20) NOT NULL DEFAULT 'none',
  ADD COLUMN IF NOT EXISTS home_content_review_note TEXT,
  ADD COLUMN IF NOT EXISTS home_content_submitted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS home_content_reviewed_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS home_content_reviewed_by BIGINT REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_sub_sites_home_content_review_status
  ON sub_sites(home_content_review_status);
