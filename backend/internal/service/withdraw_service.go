package service

import (
	"context"
	"fmt"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrWithdrawNotFound      = infraerrors.NotFound("WITHDRAW_NOT_FOUND", "withdraw request not found")
	ErrWithdrawInvalidStatus = infraerrors.BadRequest("WITHDRAW_INVALID_STATUS", "withdraw request is not in the expected status")
	ErrWithdrawAmountInvalid = infraerrors.BadRequest("WITHDRAW_AMOUNT_INVALID", "withdraw amount must be positive")
	ErrWithdrawInsufficient  = infraerrors.BadRequest("WITHDRAW_INSUFFICIENT", "insufficient balance for withdrawal")
)

type WithdrawService struct {
	agentRepo      AgentRepository
	subSiteService *SubSiteService
}

func NewWithdrawService(agentRepo AgentRepository, subSiteService *SubSiteService) *WithdrawService {
	return &WithdrawService{
		agentRepo:      agentRepo,
		subSiteService: subSiteService,
	}
}

// ApplyWithdraw 用户申请提现。
// sourceType = "agent_commission" → 扣 agent_profiles.withdrawable_balance
// sourceType = "sub_site_profit"  → 扣 sub_sites.balance_fen（rate 模式分站利润）
func (s *WithdrawService) ApplyWithdraw(ctx context.Context, userID int64, amount float64,
	alipayName, alipayAccount, alipayQR, sourceType string, sourceSubSiteID *int64) (*WithdrawRequest, error) {
	if amount <= 0 {
		return nil, ErrWithdrawAmountInvalid
	}
	switch sourceType {
	case WithdrawSourceAgentCommission:
		profile, err := s.agentRepo.GetProfile(ctx, userID)
		if err != nil {
			return nil, err
		}
		if profile.WithdrawableBalance < amount {
			return nil, ErrWithdrawInsufficient
		}
		if err := s.agentRepo.AdjustWithdrawableBalance(ctx, userID, -amount); err != nil {
			return nil, err
		}
		if err := s.agentRepo.AddWalletLog(ctx, userID,
			AgentBalanceTypeWithdrawable, AgentWalletChangeWithdrawApply,
			-amount, nil, nil, fmt.Sprintf("提现申请 %.2f 元", amount), nil); err != nil {
			return nil, err
		}
	case WithdrawSourceSubSiteProfit:
		if sourceSubSiteID == nil || *sourceSubSiteID <= 0 {
			return nil, infraerrors.BadRequest("WITHDRAW_SITE_REQUIRED", "source sub-site is required for sub_site_profit withdrawal")
		}
		site, err := s.subSiteService.GetOwnedSite(ctx, userID, *sourceSubSiteID)
		if err != nil {
			return nil, err
		}
		if site.Mode != SubSiteModeRate {
			return nil, infraerrors.BadRequest("WITHDRAW_MODE_INVALID", "only rate-mode sub-sites support profit withdrawal")
		}
		amountFen := int64(amount * 100)
		if site.BalanceFen < amountFen {
			return nil, ErrWithdrawInsufficient
		}
		entry := SubSiteLedgerEntry{TxType: SubSiteLedgerWithdrawApply}
		if _, err := s.subSiteService.AdjustPoolBalance(ctx, site.ID, -amountFen, entry); err != nil {
			return nil, err
		}
	default:
		return nil, infraerrors.BadRequest("WITHDRAW_SOURCE_INVALID", "source_type must be agent_commission or sub_site_profit")
	}

	req := &WithdrawRequest{
		UserID:          userID,
		Amount:          amount,
		AlipayName:      alipayName,
		AlipayAccount:   alipayAccount,
		AlipayQRImage:   alipayQR,
		SourceType:      sourceType,
		SourceSubSiteID: sourceSubSiteID,
		Status:          WithdrawStatusPending,
	}
	if err := s.agentRepo.CreateWithdrawRequest(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

// ReviewWithdraw 管理员审核提现。approve=true → approved；approve=false → rejected（退款回原账户）。
func (s *WithdrawService) ReviewWithdraw(ctx context.Context, requestID int64, approve bool, reviewNote string) (*WithdrawRequest, error) {
	req, err := s.agentRepo.GetWithdrawRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrWithdrawNotFound
	}
	if req.Status != WithdrawStatusPending {
		return nil, ErrWithdrawInvalidStatus
	}
	if approve {
		if err := s.agentRepo.UpdateWithdrawRequestStatus(ctx, requestID, WithdrawStatusApproved, reviewNote); err != nil {
			return nil, err
		}
		req.Status = WithdrawStatusApproved
	} else {
		if err := s.refundWithdraw(ctx, req); err != nil {
			return nil, err
		}
		if err := s.agentRepo.UpdateWithdrawRequestStatus(ctx, requestID, WithdrawStatusRejected, reviewNote); err != nil {
			return nil, err
		}
		req.Status = WithdrawStatusRejected
	}
	return req, nil
}

// PayWithdraw 管理员确认打款完成。仅 approved → paid。
func (s *WithdrawService) PayWithdraw(ctx context.Context, requestID int64) (*WithdrawRequest, error) {
	req, err := s.agentRepo.GetWithdrawRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrWithdrawNotFound
	}
	if req.Status != WithdrawStatusApproved {
		return nil, ErrWithdrawInvalidStatus
	}
	if err := s.agentRepo.UpdateWithdrawRequestStatus(ctx, requestID, WithdrawStatusPaid, ""); err != nil {
		return nil, err
	}
	switch req.SourceType {
	case WithdrawSourceAgentCommission:
		_ = s.agentRepo.IncrementTotalWithdrawn(ctx, req.UserID, req.Amount)
	case WithdrawSourceSubSiteProfit:
		if req.SourceSubSiteID != nil && *req.SourceSubSiteID > 0 {
			_ = s.subSiteService.repo.IncrementTotalWithdrawnFen(ctx, *req.SourceSubSiteID, int64(req.Amount*100))
			entry := SubSiteLedgerEntry{TxType: SubSiteLedgerWithdrawPaid}
			_, _ = s.subSiteService.AdjustPoolBalance(ctx, *req.SourceSubSiteID, 0, entry)
		}
	}
	req.Status = WithdrawStatusPaid
	return req, nil
}

// CancelWithdraw 用户自行撤销 pending 提现。
func (s *WithdrawService) CancelWithdraw(ctx context.Context, userID, requestID int64) (*WithdrawRequest, error) {
	req, err := s.agentRepo.GetWithdrawRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrWithdrawNotFound
	}
	if req.UserID != userID {
		return nil, infraerrors.Forbidden("WITHDRAW_FORBIDDEN", "you can only cancel your own withdraw request")
	}
	if req.Status != WithdrawStatusPending {
		return nil, ErrWithdrawInvalidStatus
	}
	if err := s.refundWithdraw(ctx, req); err != nil {
		return nil, err
	}
	if err := s.agentRepo.UpdateWithdrawRequestStatus(ctx, requestID, WithdrawStatusCancelled, "用户撤销"); err != nil {
		return nil, err
	}
	req.Status = WithdrawStatusCancelled
	return req, nil
}

func (s *WithdrawService) refundWithdraw(ctx context.Context, req *WithdrawRequest) error {
	switch req.SourceType {
	case WithdrawSourceAgentCommission:
		if err := s.agentRepo.AdjustWithdrawableBalance(ctx, req.UserID, req.Amount); err != nil {
			return err
		}
		return s.agentRepo.AddWalletLog(ctx, req.UserID,
			AgentBalanceTypeWithdrawable, AgentWalletChangeWithdrawReject,
			req.Amount, nil, nil, "提现退回", nil)
	case WithdrawSourceSubSiteProfit:
		if req.SourceSubSiteID != nil && *req.SourceSubSiteID > 0 {
			amountFen := int64(req.Amount * 100)
			entry := SubSiteLedgerEntry{TxType: SubSiteLedgerWithdrawReject}
			if _, err := s.subSiteService.AdjustPoolBalance(ctx, *req.SourceSubSiteID, amountFen, entry); err != nil {
				return err
			}
		}
	}
	return nil
}

// ListMyWithdraws 用户查看自己的提现记录。
func (s *WithdrawService) ListMyWithdraws(ctx context.Context, userID int64, params pagination.PaginationParams, sourceType, status string) ([]WithdrawRequest, *pagination.PaginationResult, error) {
	return s.agentRepo.ListWithdrawRequests(ctx, params, userID, sourceType, status)
}

// AdminListWithdraws 管理员查看所有提现记录。
func (s *WithdrawService) AdminListWithdraws(ctx context.Context, params pagination.PaginationParams, sourceType, status string) ([]WithdrawRequest, *pagination.PaginationResult, error) {
	return s.agentRepo.ListWithdrawRequests(ctx, params, 0, sourceType, status)
}

// HasPendingWithdrawForSubSite 检查分站是否有 pending 提现（用于模式切换校验）。
func (s *WithdrawService) HasPendingWithdrawForSubSite(ctx context.Context, subSiteID int64) (bool, error) {
	return s.agentRepo.HasPendingWithdrawForSubSite(ctx, subSiteID)
}
