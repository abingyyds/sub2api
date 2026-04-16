package service

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// SubSiteAdminUser 分站后台用户列表行。
type SubSiteAdminUser struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username,omitempty"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	Balance    float64   `json:"balance"`
	BindSource string    `json:"bind_source,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	BoundAt    time.Time `json:"bound_at"`
}

// SubSiteAdminOrder 分站后台订单行。
type SubSiteAdminOrder struct {
	ID        int64      `json:"id"`
	OrderNo   string     `json:"order_no"`
	UserID    int64      `json:"user_id"`
	UserEmail string     `json:"user_email,omitempty"`
	PlanKey   string     `json:"plan_key,omitempty"`
	OrderType string     `json:"order_type"`
	AmountFen int        `json:"amount_fen"`
	Status    string     `json:"status"`
	PayMethod string     `json:"pay_method,omitempty"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// SubSiteAdminUsage 分站后台用量行。
type SubSiteAdminUsage struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"user_id"`
	UserEmail           string    `json:"user_email,omitempty"`
	Model               string    `json:"model"`
	InputTokens         int       `json:"input_tokens"`
	OutputTokens        int       `json:"output_tokens"`
	CacheCreationTokens int       `json:"cache_creation_tokens"`
	CacheReadTokens     int       `json:"cache_read_tokens"`
	TotalCost           float64   `json:"total_cost"`
	ActualCost          float64   `json:"actual_cost"`
	CreatedAt           time.Time `json:"created_at"`
}

// SubSiteAdminDashboardStats 分站后台概览卡片数据。RangeStart 标记区间起点（一般为 30 天前）。
type SubSiteAdminDashboardStats struct {
	UserCount        int64     `json:"user_count"`
	ActiveUsers      int64     `json:"active_users"`
	Requests         int64     `json:"requests"`
	TotalCost        float64   `json:"total_cost"`
	RevenueFen       int64     `json:"revenue_fen"`
	PoolBalanceFen   int64     `json:"pool_balance_fen"`
	TotalTopupFen    int64     `json:"total_topup_fen"`
	TotalConsumedFen int64     `json:"total_consumed_fen"`
	RangeStart       time.Time `json:"range_start"`
}

// SubSiteAdminRepository 分站站长后台专用查询。
type SubSiteAdminRepository interface {
	ListUsersBySubSite(ctx context.Context, siteID int64, params pagination.PaginationParams, search string) ([]SubSiteAdminUser, *pagination.PaginationResult, error)
	ListOrdersBySubSite(ctx context.Context, siteID int64, params pagination.PaginationParams, status string) ([]SubSiteAdminOrder, *pagination.PaginationResult, error)
	ListUsageBySubSite(ctx context.Context, siteID int64, params pagination.PaginationParams, model string) ([]SubSiteAdminUsage, *pagination.PaginationResult, error)
	GetDashboardStats(ctx context.Context, siteID int64, rangeStart time.Time) (*SubSiteAdminDashboardStats, error)
}

// SubSiteAdminService 封装「先 AuthorizeOwner，再查数据」的骨架。
// 中间件是第一道闸门；这里是第二道，保证任何绕过中间件直接拿到 Service 的调用也安全。
type SubSiteAdminService struct {
	subSiteSvc *SubSiteService
	repo       SubSiteAdminRepository
}

// NewSubSiteAdminService 构造分站后台服务。
func NewSubSiteAdminService(subSiteSvc *SubSiteService, repo SubSiteAdminRepository) *SubSiteAdminService {
	return &SubSiteAdminService{subSiteSvc: subSiteSvc, repo: repo}
}

func (s *SubSiteAdminService) authorize(ctx context.Context, ownerUserID, siteID int64) error {
	_, err := s.subSiteSvc.AuthorizeOwner(ctx, ownerUserID, siteID)
	return err
}

// Dashboard 返回分站概览卡片数据（近 30 天）。
func (s *SubSiteAdminService) Dashboard(ctx context.Context, ownerUserID, siteID int64) (*SubSiteAdminDashboardStats, error) {
	if err := s.authorize(ctx, ownerUserID, siteID); err != nil {
		return nil, err
	}
	return s.repo.GetDashboardStats(ctx, siteID, time.Now().AddDate(0, 0, -30))
}

// ListUsers 返回绑定到本分站的用户列表。
func (s *SubSiteAdminService) ListUsers(ctx context.Context, ownerUserID, siteID int64, params pagination.PaginationParams, search string) ([]SubSiteAdminUser, *pagination.PaginationResult, error) {
	if err := s.authorize(ctx, ownerUserID, siteID); err != nil {
		return nil, nil, err
	}
	return s.repo.ListUsersBySubSite(ctx, siteID, params, search)
}

// ListOrders 返回本分站用户的订单列表。
func (s *SubSiteAdminService) ListOrders(ctx context.Context, ownerUserID, siteID int64, params pagination.PaginationParams, status string) ([]SubSiteAdminOrder, *pagination.PaginationResult, error) {
	if err := s.authorize(ctx, ownerUserID, siteID); err != nil {
		return nil, nil, err
	}
	return s.repo.ListOrdersBySubSite(ctx, siteID, params, status)
}

// ListUsage 返回本分站用户的用量明细。
func (s *SubSiteAdminService) ListUsage(ctx context.Context, ownerUserID, siteID int64, params pagination.PaginationParams, model string) ([]SubSiteAdminUsage, *pagination.PaginationResult, error) {
	if err := s.authorize(ctx, ownerUserID, siteID); err != nil {
		return nil, nil, err
	}
	return s.repo.ListUsageBySubSite(ctx, siteID, params, model)
}
