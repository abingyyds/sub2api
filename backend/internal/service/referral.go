package service

import "time"

// Referral represents an invite referral relationship.
type Referral struct {
	ID           int64
	InviterID    int64
	InviteeID    int64
	RewardStatus string
	RewardAmount float64
	RewardedAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Joined fields for display
	InviterEmail string
	InviteeEmail string
}

// ReferralStats holds aggregated referral statistics.
type ReferralStats struct {
	TotalInvites   int64   `json:"total_invites"`
	RewardedCount  int64   `json:"rewarded_count"`
	PendingCount   int64   `json:"pending_count"`
	TotalRewarded  float64 `json:"total_rewarded"`
}
