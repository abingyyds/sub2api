package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrAgentDisabled       = infraerrors.Forbidden("AGENT_DISABLED", "agent system is disabled")
	ErrAgentAlreadyApplied = infraerrors.Conflict("AGENT_ALREADY_APPLIED", "you have already applied or are already an agent")
	ErrAgentNotFound       = infraerrors.NotFound("AGENT_NOT_FOUND", "agent not found")
	ErrAgentNotApproved    = infraerrors.Forbidden("AGENT_NOT_APPROVED", "your agent application has not been approved")
	ErrRateExceedsOwn      = infraerrors.BadRequest("RATE_EXCEEDS_OWN", "commission rate cannot exceed your own rate")
	ErrSelfReference       = infraerrors.BadRequest("SELF_REFERENCE", "cannot set user as their own parent")
)

// AgentRepository defines the data access interface for agent operations.
type AgentRepository interface {
	CountSubUsers(ctx context.Context, agentID int64) (int64, error)
	ListSubUsers(ctx context.Context, agentID int64, params pagination.PaginationParams, search string) ([]AgentSubUser, *pagination.PaginationResult, error)
	ListSubUserPaymentOrders(ctx context.Context, agentID int64, params pagination.PaginationParams, search string) ([]AgentFinancialLog, *pagination.PaginationResult, error)
	GetDashboardStats(ctx context.Context, agentID int64) (*AgentDashboardStats, error)
	CreateCommission(ctx context.Context, c *AgentCommission) error
	ListCommissions(ctx context.Context, agentID int64, params pagination.PaginationParams, status string) ([]AgentCommission, *pagination.PaginationResult, error)
	SettlePendingCommissions(ctx context.Context, agentID int64) (float64, error)
	ListAgents(ctx context.Context, params pagination.PaginationParams, status string, search string) ([]AgentInfo, *pagination.PaginationResult, error)
	GetAgentByUserID(ctx context.Context, userID int64) (int64, *float64, error)
	GetParentAgent(ctx context.Context, agentID int64) (int64, *float64, error)
	UpdateReferralCommissionRate(ctx context.Context, inviterID, inviteeID int64, rate float64) error
	UpdateReferralInviter(ctx context.Context, inviteeID, newInviterID int64) error
}

// AgentService handles agent/affiliate business logic.
type AgentService struct {
	agentRepo      AgentRepository
	userRepo       UserRepository
	referralSvc    *ReferralService
	settingService *SettingService
}

// NewAgentService creates a new AgentService.
func NewAgentService(
	agentRepo AgentRepository,
	userRepo UserRepository,
	referralSvc *ReferralService,
	settingService *SettingService,
) *AgentService {
	return &AgentService{
		agentRepo:      agentRepo,
		userRepo:       userRepo,
		referralSvc:    referralSvc,
		settingService: settingService,
	}
}

// ApplyForAgent submits an agent application.
func (s *AgentService) ApplyForAgent(ctx context.Context, userID int64, note string) error {
	if !s.settingService.IsAgentEnabled(ctx) {
		return ErrAgentDisabled
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	// Already applied or approved
	if user.AgentStatus == AgentStatusPending || user.AgentStatus == AgentStatusApproved {
		return ErrAgentAlreadyApplied
	}

	user.IsAgent = true
	user.AgentStatus = AgentStatusPending
	user.AgentNote = note
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	log.Printf("[Agent] user %d applied for agent", userID)
	return nil
}

// GetAgentStatus returns agent status for a user.
func (s *AgentService) GetAgentStatus(ctx context.Context, userID int64) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"is_agent":     user.IsAgent,
		"agent_status": user.AgentStatus,
		"enabled":      s.settingService.IsAgentEnabled(ctx),
	}

	if user.IsAgent && user.AgentStatus == AgentStatusApproved {
		result["commission_rate"] = user.AgentCommissionRate

		// Include invite code
		code, err := s.referralSvc.GetOrCreateInviteCode(ctx, userID)
		if err == nil {
			result["invite_code"] = code
		}
	}

	return result, nil
}

// GetDashboard returns the agent dashboard stats.
func (s *AgentService) GetDashboard(ctx context.Context, userID int64) (*AgentDashboardStats, error) {
	if err := s.requireApprovedAgent(ctx, userID); err != nil {
		return nil, err
	}
	return s.agentRepo.GetDashboardStats(ctx, userID)
}

// GetInviteLink returns the agent's invite link (code).
func (s *AgentService) GetInviteLink(ctx context.Context, userID int64) (string, error) {
	if err := s.requireApprovedAgent(ctx, userID); err != nil {
		return "", err
	}
	return s.referralSvc.GetOrCreateInviteCode(ctx, userID)
}

// ListSubUsers returns the agent's sub-users.
func (s *AgentService) ListSubUsers(ctx context.Context, userID int64, params pagination.PaginationParams, search string) ([]AgentSubUser, *pagination.PaginationResult, error) {
	if err := s.requireApprovedAgent(ctx, userID); err != nil {
		return nil, nil, err
	}
	return s.agentRepo.ListSubUsers(ctx, userID, params, search)
}

// ListSubUserFinancialLogs returns financial logs of the agent's sub-users.
func (s *AgentService) ListSubUserFinancialLogs(ctx context.Context, userID int64, params pagination.PaginationParams, search string) ([]AgentFinancialLog, *pagination.PaginationResult, error) {
	if err := s.requireApprovedAgent(ctx, userID); err != nil {
		return nil, nil, err
	}
	return s.agentRepo.ListSubUserPaymentOrders(ctx, userID, params, search)
}

// ListCommissions returns the agent's commission records.
func (s *AgentService) ListCommissions(ctx context.Context, userID int64, params pagination.PaginationParams, status string) ([]AgentCommission, *pagination.PaginationResult, error) {
	if err := s.requireApprovedAgent(ctx, userID); err != nil {
		return nil, nil, err
	}
	return s.agentRepo.ListCommissions(ctx, userID, params, status)
}

// TriggerCommissionForPayment is called when a payment is completed. It checks if
// the paying user has an agent (via referrals), and if so, creates a commission record.
// It also checks for a parent agent and creates a differential commission if applicable.
func (s *AgentService) TriggerCommissionForPayment(ctx context.Context, userID int64, orderID int64, paymentAmount float64) {
	if !s.settingService.IsAgentEnabled(ctx) {
		return
	}

	// Step 1: Find direct agent + per-user rate
	agentID, perUserRate, err := s.agentRepo.GetAgentByUserID(ctx, userID)
	if err != nil || agentID == 0 {
		return // no agent for this user
	}

	// Step 2: Get agent's global commission rate
	agent, err := s.userRepo.GetByID(ctx, agentID)
	if err != nil || !agent.IsAgent || agent.AgentStatus != AgentStatusApproved {
		return
	}

	// Determine effective rate: per-user rate takes priority, otherwise use agent's global rate
	effectiveRate := agent.AgentCommissionRate
	if perUserRate != nil {
		effectiveRate = *perUserRate
	}

	if effectiveRate <= 0 {
		return
	}

	// Step 3: Create direct agent commission record
	commission := &AgentCommission{
		AgentID:          agentID,
		UserID:           userID,
		OrderID:          &orderID,
		SourceType:       AgentCommissionSourcePayment,
		SourceAmount:     paymentAmount,
		CommissionRate:   effectiveRate,
		CommissionAmount: paymentAmount * effectiveRate,
		Status:           AgentCommissionStatusPending,
	}

	if err := s.agentRepo.CreateCommission(ctx, commission); err != nil {
		log.Printf("[Agent] failed to create commission for agent=%d user=%d order=%d: %v", agentID, userID, orderID, err)
		return
	}

	log.Printf("[Agent] commission created: agent=%d user=%d order=%d amount=%.8f", agentID, userID, orderID, commission.CommissionAmount)

	// Step 4: Find parent agent (the direct agent's inviter who is also an approved agent)
	parentID, parentPerUserRate, err := s.agentRepo.GetParentAgent(ctx, agentID)
	if err != nil || parentID == 0 {
		return // no parent agent
	}

	// Step 5: Get parent agent's global rate
	parentAgent, err := s.userRepo.GetByID(ctx, parentID)
	if err != nil || !parentAgent.IsAgent || parentAgent.AgentStatus != AgentStatusApproved {
		return
	}

	// Determine parent's effective rate for this child agent
	parentEffectiveRate := parentAgent.AgentCommissionRate
	if parentPerUserRate != nil {
		parentEffectiveRate = *parentPerUserRate
	}

	// Step 6: Calculate differential = parent's rate for this child - child's effective rate
	differential := parentEffectiveRate - effectiveRate
	if differential <= 0 {
		return
	}

	// Step 7: Create differential commission record for parent
	diffCommission := &AgentCommission{
		AgentID:          parentID,
		UserID:           userID,
		OrderID:          &orderID,
		SourceType:       AgentCommissionSourceDifferential,
		SourceAmount:     paymentAmount,
		CommissionRate:   differential,
		CommissionAmount: paymentAmount * differential,
		Status:           AgentCommissionStatusPending,
	}

	if err := s.agentRepo.CreateCommission(ctx, diffCommission); err != nil {
		log.Printf("[Agent] failed to create differential commission for parent=%d agent=%d order=%d: %v", parentID, agentID, orderID, err)
		return
	}

	log.Printf("[Agent] differential commission created: parent=%d agent=%d order=%d amount=%.8f", parentID, agentID, orderID, diffCommission.CommissionAmount)
}

// SetSubUserCommissionRate sets a per-user commission rate for a sub-user.
// The rate must not exceed the agent's own commission rate.
func (s *AgentService) SetSubUserCommissionRate(ctx context.Context, agentID int64, subUserID int64, rate float64) error {
	if err := s.requireApprovedAgent(ctx, agentID); err != nil {
		return err
	}

	// Get the agent's own commission rate
	agent, err := s.userRepo.GetByID(ctx, agentID)
	if err != nil {
		return err
	}

	if rate > agent.AgentCommissionRate {
		return ErrRateExceedsOwn
	}

	if rate < 0 {
		return infraerrors.BadRequest("INVALID_RATE", "commission rate must be non-negative")
	}

	// Verify the sub-user is indeed a sub-user of this agent
	if err := s.agentRepo.UpdateReferralCommissionRate(ctx, agentID, subUserID, rate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrReferralNotFound
		}
		return fmt.Errorf("update referral commission rate: %w", err)
	}

	log.Printf("[Agent] set commission rate for agent=%d sub-user=%d rate=%.4f", agentID, subUserID, rate)
	return nil
}

// AdminUpdateParentAgent updates the parent agent (inviter) for a user.
func (s *AgentService) AdminUpdateParentAgent(ctx context.Context, userID int64, newParentID int64) error {
	if userID == newParentID {
		return ErrSelfReference
	}

	// Verify the new parent exists and is an approved agent
	parent, err := s.userRepo.GetByID(ctx, newParentID)
	if err != nil {
		return ErrAgentNotFound
	}
	if !parent.IsAgent || parent.AgentStatus != AgentStatusApproved {
		return ErrAgentNotApproved
	}

	// Verify the target user exists
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return infraerrors.NotFound("USER_NOT_FOUND", "user not found")
	}

	if err := s.agentRepo.UpdateReferralInviter(ctx, userID, newParentID); err != nil {
		return fmt.Errorf("update referral inviter: %w", err)
	}

	log.Printf("[Agent] admin updated parent for user=%d to parent=%d", userID, newParentID)
	return nil
}

// --- Admin operations ---

// AdminListAgents returns all agents (for admin).
func (s *AgentService) AdminListAgents(ctx context.Context, params pagination.PaginationParams, status string, search string) ([]AgentInfo, *pagination.PaginationResult, error) {
	return s.agentRepo.ListAgents(ctx, params, status, search)
}

// AdminApproveAgent approves an agent application.
func (s *AgentService) AdminApproveAgent(ctx context.Context, userID int64, commissionRate float64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrAgentNotFound
	}

	if commissionRate <= 0 {
		commissionRate = s.settingService.GetAgentDefaultCommissionRate(ctx)
	}

	now := time.Now()
	user.IsAgent = true
	user.AgentStatus = AgentStatusApproved
	user.AgentCommissionRate = commissionRate
	user.AgentApprovedAt = &now

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	// Ensure the agent has an invite code
	if _, err := s.referralSvc.GetOrCreateInviteCode(ctx, userID); err != nil {
		log.Printf("[Agent] failed to create invite code for agent %d: %v", userID, err)
	}

	log.Printf("[Agent] approved agent user=%d rate=%.4f", userID, commissionRate)
	return nil
}

// AdminRejectAgent rejects an agent application.
func (s *AgentService) AdminRejectAgent(ctx context.Context, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrAgentNotFound
	}

	user.AgentStatus = AgentStatusRejected
	user.IsAgent = false
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	log.Printf("[Agent] rejected agent user=%d", userID)
	return nil
}

// AdminUpdateCommissionRate updates the agent's commission rate.
func (s *AgentService) AdminUpdateCommissionRate(ctx context.Context, userID int64, rate float64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrAgentNotFound
	}

	user.AgentCommissionRate = rate
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	log.Printf("[Agent] updated commission rate for agent %d to %.4f", userID, rate)
	return nil
}

// AdminSettleCommissions settles all pending commissions for an agent.
func (s *AgentService) AdminSettleCommissions(ctx context.Context, agentID int64) (float64, error) {
	amount, err := s.agentRepo.SettlePendingCommissions(ctx, agentID)
	if err != nil {
		return 0, fmt.Errorf("settle commissions: %w", err)
	}

	if amount > 0 {
		if err := s.userRepo.UpdateBalance(ctx, agentID, amount); err != nil {
			return 0, fmt.Errorf("add balance: %w", err)
		}
		log.Printf("[Agent] settled commissions for agent %d, amount=%.8f", agentID, amount)
	}

	return amount, nil
}

// --- helpers ---

func (s *AgentService) requireApprovedAgent(ctx context.Context, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !user.IsAgent || user.AgentStatus != AgentStatusApproved {
		return ErrAgentNotApproved
	}
	return nil
}
