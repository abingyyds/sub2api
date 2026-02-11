package service

import "time"

// Referral represents an invite referral relationship.
type Referral struct {
	ID           int64      `json:"id"`
	InviterID    int64      `json:"inviter_id"`
	InviteeID    int64      `json:"invitee_id"`
	RewardStatus string     `json:"reward_status"`
	RewardAmount float64    `json:"reward_amount"`
	RewardedAt   *time.Time `json:"rewarded_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Joined fields for display
	InviterEmail string `json:"inviter_email,omitempty"`
	InviteeEmail string `json:"invitee_email,omitempty"`
}

// ReferralStats holds aggregated referral statistics.
type ReferralStats struct {
	TotalInvites  int64   `json:"total_invitees"`
	RewardedCount int64   `json:"rewarded_count"`
	PendingCount  int64   `json:"pending_count"`
	TotalRewarded float64 `json:"total_reward_amount"`
}
