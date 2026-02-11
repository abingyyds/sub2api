-- 045: Add referral/invite reward system
-- Each user gets a unique invite code; invitees trigger reward on first balance redemption.

ALTER TABLE users ADD COLUMN IF NOT EXISTS invite_code VARCHAR(16);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_invite_code ON users(invite_code) WHERE invite_code IS NOT NULL AND deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS referrals (
    id BIGSERIAL PRIMARY KEY,
    inviter_id BIGINT NOT NULL REFERENCES users(id),
    invitee_id BIGINT NOT NULL REFERENCES users(id),
    reward_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    reward_amount DECIMAL(20,8) NOT NULL DEFAULT 0,
    rewarded_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_referrals_invitee_id ON referrals(invitee_id);
CREATE INDEX IF NOT EXISTS idx_referrals_inviter_id ON referrals(inviter_id);
CREATE INDEX IF NOT EXISTS idx_referrals_reward_status ON referrals(reward_status);
