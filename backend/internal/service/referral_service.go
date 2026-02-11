package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrReferralNotFound    = infraerrors.NotFound("REFERRAL_NOT_FOUND", "referral not found")
	ErrReferralSelfInvite  = infraerrors.BadRequest("REFERRAL_SELF_INVITE", "cannot invite yourself")
	ErrReferralDisabled    = infraerrors.Forbidden("REFERRAL_DISABLED", "referral feature is disabled")
	ErrInviteCodeNotFound  = infraerrors.NotFound("INVITE_CODE_NOT_FOUND", "invite code not found")
	ErrAlreadyReferred     = infraerrors.Conflict("ALREADY_REFERRED", "user already has a referrer")
)

// ReferralRepository defines the data access interface for referrals.
type ReferralRepository interface {
	Create(ctx context.Context, referral *Referral) error
	GetByInviteeID(ctx context.Context, inviteeID int64) (*Referral, error)
	UpdateRewardStatus(ctx context.Context, id int64, status string, amount float64) error
	ListByInviterID(ctx context.Context, inviterID int64, params pagination.PaginationParams) ([]Referral, *pagination.PaginationResult, error)
	ListAll(ctx context.Context, params pagination.PaginationParams, search string) ([]Referral, *pagination.PaginationResult, error)
	GetStatsByInviterID(ctx context.Context, inviterID int64) (*ReferralStats, error)
	GetGlobalStats(ctx context.Context) (*ReferralStats, error)
}

// ReferralService handles referral/invite reward logic.
type ReferralService struct {
	referralRepo   ReferralRepository
	userRepo       UserRepository
	settingService *SettingService
}

// NewReferralService creates a new ReferralService.
func NewReferralService(
	referralRepo ReferralRepository,
	userRepo UserRepository,
	settingService *SettingService,
) *ReferralService {
	return &ReferralService{
		referralRepo:   referralRepo,
		userRepo:       userRepo,
		settingService: settingService,
	}
}

// IsReferralEnabled checks if the referral feature is enabled.
func (s *ReferralService) IsReferralEnabled(ctx context.Context) bool {
	return s.settingService.IsReferralEnabled(ctx)
}

// generateInviteCode generates a random 8-char uppercase hex invite code.
func generateInviteCode() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(b)), nil
}

// GetOrCreateInviteCode returns the user's invite code, generating one if needed.
func (s *ReferralService) GetOrCreateInviteCode(ctx context.Context, userID int64) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	if user.InviteCode != nil && *user.InviteCode != "" {
		return *user.InviteCode, nil
	}

	// Generate and save a new invite code (retry on collision)
	for i := 0; i < 5; i++ {
		code, err := generateInviteCode()
		if err != nil {
			return "", fmt.Errorf("generate invite code: %w", err)
		}
		user.InviteCode = &code
		if err := s.userRepo.Update(ctx, user); err != nil {
			// Retry on unique constraint violation
			log.Printf("[Referral] invite code collision, retrying: %v", err)
			continue
		}
		return code, nil
	}
	return "", fmt.Errorf("failed to generate unique invite code after retries")
}

// RecordReferral records a referral relationship when invitee registers.
func (s *ReferralService) RecordReferral(ctx context.Context, inviteCode string, inviteeID int64) error {
	if !s.IsReferralEnabled(ctx) {
		return nil // silently skip if disabled
	}

	inviteCode = strings.TrimSpace(inviteCode)
	if inviteCode == "" {
		return nil
	}

	// Find inviter by invite code
	inviter, err := s.userRepo.GetByInviteCode(ctx, inviteCode)
	if err != nil {
		log.Printf("[Referral] invite code %q not found: %v", inviteCode, err)
		return nil // don't block registration
	}

	// Prevent self-invite
	if inviter.ID == inviteeID {
		return nil
	}

	// Check if invitee already has a referrer
	existing, _ := s.referralRepo.GetByInviteeID(ctx, inviteeID)
	if existing != nil {
		return nil // already referred
	}

	referral := &Referral{
		InviterID:    inviter.ID,
		InviteeID:    inviteeID,
		RewardStatus: ReferralRewardPending,
	}
	if err := s.referralRepo.Create(ctx, referral); err != nil {
		log.Printf("[Referral] failed to record referral for invitee %d: %v", inviteeID, err)
		return nil // don't block registration
	}

	log.Printf("[Referral] recorded referral: inviter=%d invitee=%d", inviter.ID, inviteeID)
	return nil
}

// TriggerRewardIfEligible triggers the referral reward when invitee first redeems balance.
func (s *ReferralService) TriggerRewardIfEligible(ctx context.Context, inviteeID int64) {
	if !s.IsReferralEnabled(ctx) {
		return
	}

	referral, err := s.referralRepo.GetByInviteeID(ctx, inviteeID)
	if err != nil || referral == nil {
		return // no referral relationship
	}

	if referral.RewardStatus != ReferralRewardPending {
		return // already rewarded
	}

	rewardAmount := s.settingService.GetReferralRewardAmount(ctx)
	if rewardAmount <= 0 {
		return
	}

	// Update referral status and add balance to inviter
	if err := s.referralRepo.UpdateRewardStatus(ctx, referral.ID, ReferralRewardRewarded, rewardAmount); err != nil {
		log.Printf("[Referral] failed to update reward status for referral %d: %v", referral.ID, err)
		return
	}

	if err := s.userRepo.UpdateBalance(ctx, referral.InviterID, rewardAmount); err != nil {
		log.Printf("[Referral] failed to add reward balance to inviter %d: %v", referral.InviterID, err)
		return
	}

	log.Printf("[Referral] rewarded inviter=%d amount=%.8f for invitee=%d", referral.InviterID, rewardAmount, inviteeID)
}

// ListInvitees returns the list of invitees for a given user.
func (s *ReferralService) ListInvitees(ctx context.Context, userID int64, params pagination.PaginationParams) ([]Referral, *pagination.PaginationResult, error) {
	return s.referralRepo.ListByInviterID(ctx, userID, params)
}

// ListAll returns all referral records (admin).
func (s *ReferralService) ListAll(ctx context.Context, params pagination.PaginationParams, search string) ([]Referral, *pagination.PaginationResult, error) {
	return s.referralRepo.ListAll(ctx, params, search)
}

// GetMyStats returns referral stats for a specific user.
func (s *ReferralService) GetMyStats(ctx context.Context, userID int64) (*ReferralStats, error) {
	return s.referralRepo.GetStatsByInviterID(ctx, userID)
}

// GetStats returns global referral stats (admin).
func (s *ReferralService) GetStats(ctx context.Context) (*ReferralStats, error) {
	return s.referralRepo.GetGlobalStats(ctx)
}
