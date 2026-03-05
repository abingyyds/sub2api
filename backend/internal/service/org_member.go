package service

import (
	"context"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// OrgMember errors
var (
	ErrOrgMemberNotFound = infraerrors.NotFound("ORG_MEMBER_NOT_FOUND", "organization member not found")
	ErrOrgMemberConflict = infraerrors.Conflict("ORG_MEMBER_CONFLICT", "user is already a member of this organization")
	ErrOrgMemberDailyQuotaExceeded = infraerrors.New(429, "ORG_MEMBER_DAILY_QUOTA", "daily quota exceeded for organization member")
	ErrOrgMemberMonthlyQuotaExceeded = infraerrors.New(429, "ORG_MEMBER_MONTHLY_QUOTA", "monthly quota exceeded for organization member")
)

type OrgMember struct {
	ID               int64
	OrgID            int64
	UserID           int64
	Role             string
	MonthlyQuotaUSD  *float64
	DailyQuotaUSD    *float64
	MonthlyUsageUSD  float64
	DailyUsageUSD    float64
	MonthlyWindowStart *time.Time
	DailyWindowStart   *time.Time
	Status           string
	Notes            *string
	CreatedAt        time.Time
	UpdatedAt        time.Time

	Organization *Organization
	User         *User
}

func (m *OrgMember) IsActive() bool {
	return m.Status == StatusActive
}

func (m *OrgMember) IsAdmin() bool {
	return m.Role == OrgMemberRoleAdmin
}

func (m *OrgMember) NeedsDailyReset() bool {
	if m.DailyWindowStart == nil {
		return false
	}
	return time.Since(*m.DailyWindowStart) >= 24*time.Hour
}

func (m *OrgMember) NeedsMonthlyReset() bool {
	if m.MonthlyWindowStart == nil {
		return false
	}
	return time.Since(*m.MonthlyWindowStart) >= 30*24*time.Hour
}

func (m *OrgMember) CheckDailyQuota(additionalCost float64) bool {
	if m.DailyQuotaUSD == nil {
		return true
	}
	return m.DailyUsageUSD+additionalCost <= *m.DailyQuotaUSD
}

func (m *OrgMember) CheckMonthlyQuota(additionalCost float64) bool {
	if m.MonthlyQuotaUSD == nil {
		return true
	}
	return m.MonthlyUsageUSD+additionalCost <= *m.MonthlyQuotaUSD
}

type OrgMemberRepository interface {
	Create(ctx context.Context, member *OrgMember) error
	GetByID(ctx context.Context, id int64) (*OrgMember, error)
	GetByOrgAndUser(ctx context.Context, orgID, userID int64) (*OrgMember, error)
	Update(ctx context.Context, member *OrgMember) error
	Delete(ctx context.Context, id int64) error
	ListByOrg(ctx context.Context, orgID int64, params pagination.PaginationParams) ([]OrgMember, *pagination.PaginationResult, error)

	IncrementUsage(ctx context.Context, id int64, costUSD float64) error
	ResetDailyUsage(ctx context.Context, id int64, newWindowStart time.Time) error
	ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error
}
