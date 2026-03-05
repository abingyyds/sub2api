-- +migrate Up
ALTER TABLE api_keys ADD COLUMN org_project_id BIGINT REFERENCES org_projects(id);
CREATE INDEX idx_api_keys_org_project_id ON api_keys(org_project_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_api_keys_org_project_id;
ALTER TABLE api_keys DROP COLUMN IF EXISTS org_project_id;
