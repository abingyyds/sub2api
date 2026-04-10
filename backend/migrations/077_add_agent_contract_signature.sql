-- 077: add contract signature evidence for agent onboarding

ALTER TABLE agent_profiles
    ADD COLUMN IF NOT EXISTS contract_signature_data TEXT NOT NULL DEFAULT '';
