package service

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// OrgProject errors
var (
	ErrOrgProjectNotFound      = infraerrors.NotFound("ORG_PROJECT_NOT_FOUND", "organization project not found")
	ErrOrgProjectConflict      = infraerrors.Conflict("ORG_PROJECT_CONFLICT", "project name already exists in this organization")
	ErrOrgProjectBudgetExceeded = infraerrors.Forbidden("ORG_PROJECT_BUDGET_EXCEEDED", "project monthly budget exceeded")
	ErrOrgProjectModelNotAllowed = infraerrors.Forbidden("ORG_PROJECT_MODEL_NOT_ALLOWED", "model is not allowed for this project")
)

type OrgProject struct {
	ID               int64
	OrgID            int64
	Name             string
	Description      *string
	GroupID          *int64
	AllowedModels    []string
	MonthlyBudgetUSD *float64
	MonthlyUsageUSD  float64
	MonthlyWindowStart *time.Time
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time

	Organization *Organization
}

// IsActive returns true if the project is active.
func (p *OrgProject) IsActive() bool {
	return p.Status == OrgStatusActive
}

// IsModelAllowed checks if the given model matches the project's allowed models list.
// Returns true if the allow list is empty (all models allowed).
// Supports wildcard matching (e.g., "claude-sonnet-4-*" matches "claude-sonnet-4-20250514").
func (p *OrgProject) IsModelAllowed(model string) bool {
	if len(p.AllowedModels) == 0 {
		return true
	}
	for _, pattern := range p.AllowedModels {
		matched, _ := filepath.Match(strings.ToLower(pattern), strings.ToLower(model))
		if matched {
			return true
		}
		// Exact match fallback
		if strings.EqualFold(pattern, model) {
			return true
		}
	}
	return false
}

// CheckBudget checks if the project has exceeded its monthly budget.
func (p *OrgProject) CheckBudget(additionalCost float64) bool {
	if p.MonthlyBudgetUSD == nil {
		return true
	}
	return (p.MonthlyUsageUSD + additionalCost) <= *p.MonthlyBudgetUSD
}

// NeedsMonthlyReset returns true if the monthly usage window needs resetting.
func (p *OrgProject) NeedsMonthlyReset() bool {
	if p.MonthlyWindowStart == nil {
		return true
	}
	now := time.Now()
	return now.Year() != p.MonthlyWindowStart.Year() || now.Month() != p.MonthlyWindowStart.Month()
}

type OrgProjectRepository interface {
	Create(ctx context.Context, project *OrgProject) error
	GetByID(ctx context.Context, id int64) (*OrgProject, error)
	Update(ctx context.Context, project *OrgProject) error
	Delete(ctx context.Context, id int64) error
	ListByOrg(ctx context.Context, orgID int64, params pagination.PaginationParams) ([]OrgProject, *pagination.PaginationResult, error)

	IncrementUsage(ctx context.Context, id int64, costUSD float64) error
	ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error
}
