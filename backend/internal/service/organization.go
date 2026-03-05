package service

import (
	"context"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// Organization errors
var (
	ErrOrganizationNotFound = infraerrors.NotFound("ORGANIZATION_NOT_FOUND", "organization not found")
	ErrOrganizationConflict = infraerrors.Conflict("ORGANIZATION_CONFLICT", "organization slug already exists")
	ErrOrganizationDisabled = infraerrors.Forbidden("ORGANIZATION_DISABLED", "organization is disabled or suspended")
	ErrOrgInsufficientBalance = infraerrors.Forbidden("ORG_INSUFFICIENT_BALANCE", "organization has insufficient balance")
	ErrOrgMaxMembersReached = infraerrors.Forbidden("ORG_MAX_MEMBERS", "organization has reached maximum member limit")
	ErrOrgMaxAPIKeysReached = infraerrors.Forbidden("ORG_MAX_API_KEYS", "organization has reached maximum API key limit")
)

type Organization struct {
	ID             int64
	Name           string
	Slug           string
	Description    *string
	OwnerUserID    int64
	BillingMode    string
	Balance        float64
	MonthlyBudgetUSD *float64
	MaxMembers     int
	MaxAPIKeys     int
	Status         string
	AuditMode      string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Owner   *User
	Members []OrgMember
}

func (o *Organization) IsActive() bool {
	return o.Status == OrgStatusActive
}

func (o *Organization) IsSuspended() bool {
	return o.Status == OrgStatusSuspended
}

type OrganizationRepository interface {
	Create(ctx context.Context, org *Organization) error
	GetByID(ctx context.Context, id int64) (*Organization, error)
	GetBySlug(ctx context.Context, slug string) (*Organization, error)
	GetByOwnerID(ctx context.Context, ownerUserID int64) (*Organization, error)
	Update(ctx context.Context, org *Organization) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, params pagination.PaginationParams, status, search string) ([]Organization, *pagination.PaginationResult, error)

	DeductBalance(ctx context.Context, orgID int64, amount float64) error
	AddBalance(ctx context.Context, orgID int64, amount float64) error
	CountMembers(ctx context.Context, orgID int64) (int, error)
	CountAPIKeys(ctx context.Context, orgID int64) (int, error)
}
