-- +migrate Up
-- Only add the column if org_projects table exists
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'org_projects') THEN
        ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS org_project_id BIGINT REFERENCES org_projects(id);
        CREATE INDEX IF NOT EXISTS idx_api_keys_org_project_id ON api_keys(org_project_id);
    END IF;
END $$;

-- +migrate Down
DROP INDEX IF EXISTS idx_api_keys_org_project_id;
ALTER TABLE api_keys DROP COLUMN IF EXISTS org_project_id;
