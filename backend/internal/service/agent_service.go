package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrAgentDisabled                  = infraerrors.Forbidden("AGENT_DISABLED", "agent system is disabled")
	ErrAgentAlreadyApplied            = infraerrors.Conflict("AGENT_ALREADY_APPLIED", "you have already applied or are already an agent")
	ErrAgentNotFound                  = infraerrors.NotFound("AGENT_NOT_FOUND", "agent not found")
	ErrAgentNotApproved               = infraerrors.Forbidden("AGENT_NOT_APPROVED", "your agent application has not been approved")
	ErrAgentFrozen                    = infraerrors.Forbidden("AGENT_FROZEN", "your agent account is frozen")
	ErrAgentPrerequisitesMissing      = infraerrors.BadRequest("AGENT_PREREQUISITES_MISSING", "complete identity, contract and activation fee payment before applying")
	ErrAgentContractSignatureRequired = infraerrors.BadRequest("AGENT_CONTRACT_SIGNATURE_REQUIRED", "please sign the contract before saving your profile")
	ErrAgentContractTemplateMissing   = infraerrors.BadRequest("AGENT_CONTRACT_TEMPLATE_MISSING", "contract template is not configured yet")
	ErrRateExceedsOwn                 = infraerrors.BadRequest("RATE_EXCEEDS_OWN", "commission rate cannot exceed your own rate")
	ErrSelfReference                  = infraerrors.BadRequest("SELF_REFERENCE", "cannot set user as their own parent")
)

// AgentRepository defines the data access interface for agent operations.
type AgentRepository interface {
	CountSubUsers(ctx context.Context, agentID int64) (int64, error)
	ListSubUsers(ctx context.Context, agentID int64, params pagination.PaginationParams, search string) ([]AgentSubUser, *pagination.PaginationResult, error)
	ListSubUserPaymentOrders(ctx context.Context, agentID int64, params pagination.PaginationParams, search string) ([]AgentFinancialLog, *pagination.PaginationResult, error)
	GetDashboardStats(ctx context.Context, agentID int64, siteBalance float64) (*AgentDashboardStats, error)
	CreateCommission(ctx context.Context, c *AgentCommission) error
	ListCommissions(ctx context.Context, agentID int64, params pagination.PaginationParams, status string) ([]AgentCommission, *pagination.PaginationResult, error)
	SettlePendingCommissions(ctx context.Context, agentID int64) (float64, error)
	ListAgents(ctx context.Context, params pagination.PaginationParams, status string, search string) ([]AgentInfo, *pagination.PaginationResult, error)
	GetAgentByUserID(ctx context.Context, userID int64) (int64, *float64, error)
	GetProfile(ctx context.Context, userID int64) (*AgentProfile, error)
	UpsertProfile(ctx context.Context, profile *AgentProfile) error
	MarkActivationFeePaid(ctx context.Context, userID int64, orderID int64) error
	SetAgentFrozen(ctx context.Context, userID int64, frozen bool, reason string) error
	AddWalletLog(ctx context.Context, userID int64, balanceType, changeType string, amount float64, relatedUserID *int64, relatedOrderID *int64, remark string, unlockAt *time.Time) error
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

// SaveProfile stores identity information and contract confirmation for a user who wants to become an agent.
func (s *AgentService) SaveProfile(ctx context.Context, userID int64, realName string, idCardNo string, phone string, contractAccepted bool, contractSignatureData string, clientIP string) error {
	if !s.settingService.IsAgentEnabled(ctx) {
		return ErrAgentDisabled
	}

	profile, err := s.agentRepo.GetProfile(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now()
	profile.RealName = realName
	profile.IDCardNo = idCardNo
	profile.Phone = phone
	contractSignatureData = strings.TrimSpace(contractSignatureData)
	if realName != "" && idCardNo != "" && phone != "" {
		profile.IdentityStatus = AgentIdentityStatusSubmitted
		profile.IdentitySubmittedAt = &now
	}
	if contractAccepted {
		if strings.TrimSpace(s.settingService.GetAgentContractTemplate(ctx)) == "" {
			return ErrAgentContractTemplateMissing
		}
		if contractSignatureData == "" && strings.TrimSpace(profile.ContractSignatureData) == "" {
			return ErrAgentContractSignatureRequired
		}
		profile.ContractStatus = AgentContractStatusSigned
		profile.ContractVersion = s.settingService.GetAgentContractVersion(ctx)
		profile.ContractSignedAt = &now
		profile.ContractIP = clientIP
		if contractSignatureData != "" {
			profile.ContractSignatureData = contractSignatureData
		}
	}

	return s.agentRepo.UpsertProfile(ctx, profile)
}

// AdminGetAgentDetail returns the full onboarding detail for a single agent/applicant.
func (s *AgentService) AdminGetAgentDetail(ctx context.Context, userID int64) (*AgentInfo, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrAgentNotFound
	}

	profile, err := s.agentRepo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	detail := &AgentInfo{
		ID:                    user.ID,
		Email:                 user.Email,
		Username:              user.Username,
		IsAgent:               user.IsAgent,
		AgentStatus:           user.AgentStatus,
		CommissionRate:        user.AgentCommissionRate,
		AgentNote:             user.AgentNote,
		ApprovedAt:            user.AgentApprovedAt,
		RealName:              profile.RealName,
		IDCardNo:              profile.IDCardNo,
		Phone:                 profile.Phone,
		IdentityStatus:        profile.IdentityStatus,
		IdentitySubmittedAt:   profile.IdentitySubmittedAt,
		ContractStatus:        profile.ContractStatus,
		ContractVersion:       profile.ContractVersion,
		ContractSignedAt:      profile.ContractSignedAt,
		ContractIP:            profile.ContractIP,
		ContractSignatureData: profile.ContractSignatureData,
		ActivationFeePaidAt:   profile.ActivationFeePaidAt,
		IsFrozen:              profile.IsFrozen,
		FrozenReason:          profile.FrozenReason,
		FrozenBalance:         profile.FrozenBalance,
		WithdrawableBalance:   profile.WithdrawableBalance,
		TotalWithdrawn:        profile.TotalWithdrawn,
		CreatedAt:             user.CreatedAt,
	}

	if user.AgentStatus == AgentStatusApproved {
		if code, codeErr := s.referralSvc.GetOrCreateInviteCode(ctx, userID); codeErr == nil {
			detail.InviteCode = code
		}
	}

	return detail, nil
}

// ApplyForAgent submits an agent application after all prerequisites are completed.
func (s *AgentService) ApplyForAgent(ctx context.Context, userID int64) error {
	if !s.settingService.IsAgentEnabled(ctx) {
		return ErrAgentDisabled
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user.AgentStatus == AgentStatusPending || user.AgentStatus == AgentStatusApproved {
		return ErrAgentAlreadyApplied
	}

	profile, err := s.agentRepo.GetProfile(ctx, userID)
	if err != nil {
		return err
	}
	if !s.profileReadyForApplication(profile) {
		return ErrAgentPrerequisitesMissing
	}

	user.IsAgent = false
	user.AgentStatus = AgentStatusPending
	user.AgentNote = ""
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	log.Printf("[Agent] user %d applied for agent", userID)
	return nil
}

// MarkActivationFeePaid records that the user has paid the agent activation fee.
func (s *AgentService) MarkActivationFeePaid(ctx context.Context, userID int64, orderID int64) error {
	return s.agentRepo.MarkActivationFeePaid(ctx, userID, orderID)
}

// GetAgentStatus returns the aggregated agent status for a user.
func (s *AgentService) GetAgentStatus(ctx context.Context, userID int64) (*AgentStatusView, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile, err := s.agentRepo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	weekday, startHour, endHour := s.settingService.GetAgentWithdrawWindow(ctx)
	status := &AgentStatusView{
		Enabled:          s.settingService.IsAgentEnabled(ctx),
		IsAgent:          user.IsAgent,
		AgentStatus:      user.AgentStatus,
		CommissionRate:   user.AgentCommissionRate,
		CanApply:         s.profileReadyForApplication(profile) && user.AgentStatus != AgentStatusPending && user.AgentStatus != AgentStatusApproved,
		ActivationFee:    s.settingService.GetAgentActivationFee(ctx),
		ContractTemplate: s.settingService.GetAgentContractTemplate(ctx),
		Profile:          profile,
		Wallet: AgentWalletSummary{
			SiteBalance:         user.Balance,
			FrozenBalance:       profile.FrozenBalance,
			WithdrawableBalance: profile.WithdrawableBalance,
			TotalWithdrawn:      profile.TotalWithdrawn,
		},
		WithdrawFreezeDays: s.settingService.GetAgentWithdrawFreezeDays(ctx),
		WithdrawWindow: AgentWithdrawWindow{
			Weekday:   weekday,
			StartHour: startHour,
			EndHour:   endHour,
			Label:     fmt.Sprintf("weekly-%d %02d:00-%02d:00", weekday, startHour, endHour),
		},
	}

	if user.AgentStatus == AgentStatusApproved {
		if code, codeErr := s.referralSvc.GetOrCreateInviteCode(ctx, userID); codeErr == nil {
			status.InviteCode = code
		}
	}

	return status, nil
}

// GetDashboard returns the agent dashboard stats.
func (s *AgentService) GetDashboard(ctx context.Context, userID int64) (*AgentDashboardStats, error) {
	user, err := s.requireApprovedAgent(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.agentRepo.GetDashboardStats(ctx, userID, user.Balance)
}

// GetInviteLink returns the agent's invite link (code).
func (s *AgentService) GetInviteLink(ctx context.Context, userID int64) (string, error) {
	if _, err := s.requireApprovedAgent(ctx, userID); err != nil {
		return "", err
	}
	return s.referralSvc.GetOrCreateInviteCode(ctx, userID)
}

// ListSubUsers returns the agent's sub-users.
func (s *AgentService) ListSubUsers(ctx context.Context, userID int64, params pagination.PaginationParams, search string) ([]AgentSubUser, *pagination.PaginationResult, error) {
	if _, err := s.requireApprovedAgent(ctx, userID); err != nil {
		return nil, nil, err
	}
	return s.agentRepo.ListSubUsers(ctx, userID, params, search)
}

// ListSubUserFinancialLogs returns financial logs of the agent's sub-users.
func (s *AgentService) ListSubUserFinancialLogs(ctx context.Context, userID int64, params pagination.PaginationParams, search string) ([]AgentFinancialLog, *pagination.PaginationResult, error) {
	if _, err := s.requireApprovedAgent(ctx, userID); err != nil {
		return nil, nil, err
	}
	return s.agentRepo.ListSubUserPaymentOrders(ctx, userID, params, search)
}

// ListCommissions returns the agent's commission records.
func (s *AgentService) ListCommissions(ctx context.Context, userID int64, params pagination.PaginationParams, status string) ([]AgentCommission, *pagination.PaginationResult, error) {
	if _, err := s.requireApprovedAgent(ctx, userID); err != nil {
		return nil, nil, err
	}
	return s.agentRepo.ListCommissions(ctx, userID, params, status)
}

// TriggerCommissionForPayment creates a single-level invite commission for agent activation fee payments.
func (s *AgentService) TriggerCommissionForPayment(ctx context.Context, userID int64, orderID int64, orderType string, paymentAmount float64) {
	if !s.settingService.IsAgentEnabled(ctx) {
		return
	}
	if orderType != PaymentOrderTypeAgentActivation {
		return
	}

	agentID, perUserRate, err := s.agentRepo.GetAgentByUserID(ctx, userID)
	if err != nil || agentID == 0 {
		return
	}

	agent, err := s.userRepo.GetByID(ctx, agentID)
	if err != nil || !agent.IsAgent || agent.AgentStatus != AgentStatusApproved {
		return
	}

	effectiveRate := agent.AgentCommissionRate
	if perUserRate != nil {
		effectiveRate = *perUserRate
	}
	if effectiveRate <= 0 {
		return
	}

	commissionAmount := paymentAmount * effectiveRate
	now := time.Now()
	commission := &AgentCommission{
		AgentID:          agentID,
		UserID:           userID,
		OrderID:          &orderID,
		SourceType:       AgentCommissionSourcePayment,
		SourceAmount:     paymentAmount,
		CommissionRate:   effectiveRate,
		CommissionAmount: commissionAmount,
		Status:           AgentCommissionStatusSettled,
		SettledAt:        &now,
	}

	if err := s.agentRepo.CreateCommission(ctx, commission); err != nil {
		log.Printf("[Agent] failed to create commission for agent=%d user=%d order=%d: %v", agentID, userID, orderID, err)
		return
	}
	if err := s.userRepo.UpdateBalance(ctx, agentID, commissionAmount); err != nil {
		log.Printf("[Agent] failed to credit site balance for agent=%d order=%d: %v", agentID, orderID, err)
		return
	}
	relatedUserID := userID
	relatedOrderID := orderID
	if err := s.agentRepo.AddWalletLog(ctx, agentID, AgentBalanceTypeSite, AgentWalletChangeInviteCommission, commissionAmount, &relatedUserID, &relatedOrderID, "agent activation invite commission", nil); err != nil {
		log.Printf("[Agent] failed to write wallet log for agent=%d order=%d: %v", agentID, orderID, err)
	}

	log.Printf("[Agent] agent activation commission credited: agent=%d user=%d order=%d amount=%.8f", agentID, userID, orderID, commissionAmount)
}

// SetSubUserCommissionRate sets a per-user commission rate for a sub-user.
// The rate must not exceed the agent's own commission rate.
func (s *AgentService) SetSubUserCommissionRate(ctx context.Context, agentID int64, subUserID int64, rate float64) error {
	if _, err := s.requireApprovedAgent(ctx, agentID); err != nil {
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
	profile, err := s.agentRepo.GetProfile(ctx, userID)
	if err != nil {
		return err
	}
	if !s.profileReadyForApplication(profile) {
		return ErrAgentPrerequisitesMissing
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

// AdminSetFrozen updates an agent's frozen status.
func (s *AgentService) AdminSetFrozen(ctx context.Context, userID int64, frozen bool, reason string) error {
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return ErrAgentNotFound
	}
	if err := s.agentRepo.SetAgentFrozen(ctx, userID, frozen, reason); err != nil {
		return err
	}
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

func (s *AgentService) requireApprovedAgent(ctx context.Context, userID int64) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !user.IsAgent || user.AgentStatus != AgentStatusApproved {
		return nil, ErrAgentNotApproved
	}
	profile, err := s.agentRepo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile.IsFrozen {
		return nil, ErrAgentFrozen
	}
	return user, nil
}

func (s *AgentService) profileReadyForApplication(profile *AgentProfile) bool {
	if profile == nil {
		return false
	}
	return profile.IdentityStatus == AgentIdentityStatusSubmitted &&
		profile.ContractStatus == AgentContractStatusSigned &&
		strings.TrimSpace(profile.ContractSignatureData) != "" &&
		profile.ActivationFeePaidAt != nil
}
