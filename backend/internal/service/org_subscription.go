package service

import (
	"context"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// OrgSubscription errors
var (
	ErrOrgSubscriptionNotFound      = infraerrors.NotFound("ORG_SUBSCRIPTION_NOT_FOUND", "organization subscription not found")
	ErrOrgSubscriptionAlreadyExists = infraerrors.Conflict("ORG_SUBSCRIPTION_CONFLICT", "organization already has an active subscription for this group")
)

type OrgSubscription struct {
	ID      int64
	OrgID   int64
	GroupID int64

	StartsAt  time.Time
	ExpiresAt time.Time
	Status    string

	DailyWindowStart   *time.Time
	WeeklyWindowStart  *time.Time
	MonthlyWindowStart *time.Time

	DailyUsageUSD   float64
	WeeklyUsageUSD  float64
	MonthlyUsageUSD float64

	AssignedBy *int64
	AssignedAt time.Time
	Notes      *string

	CreatedAt time.Time
	UpdatedAt time.Time

	Organization *Organization
	Group        *Group
}

func (s *OrgSubscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive && time.Now().Before(s.ExpiresAt)
}

func (s *OrgSubscription) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

type OrgSubscriptionRepository interface {
	Create(ctx context.Context, sub *OrgSubscription) error
	GetByID(ctx context.Context, id int64) (*OrgSubscription, error)
	GetActiveByOrgAndGroup(ctx context.Context, orgID, groupID int64) (*OrgSubscription, error)
	Update(ctx context.Context, sub *OrgSubscription) error
	Delete(ctx context.Context, id int64) error
	ListByOrg(ctx context.Context, orgID int64, params pagination.PaginationParams) ([]OrgSubscription, *pagination.PaginationResult, error)

	IncrementUsage(ctx context.Context, id int64, costUSD float64) error
	ResetDailyUsage(ctx context.Context, id int64, newWindowStart time.Time) error
	ResetWeeklyUsage(ctx context.Context, id int64, newWindowStart time.Time) error
	ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error
}
