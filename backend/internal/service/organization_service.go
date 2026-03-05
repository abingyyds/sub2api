package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

type OrganizationService struct {
	orgRepo       OrganizationRepository
	memberRepo    OrgMemberRepository
	orgSubRepo    OrgSubscriptionRepository
	userRepo      UserRepository
}

func NewOrganizationService(
	orgRepo OrganizationRepository,
	memberRepo OrgMemberRepository,
	orgSubRepo OrgSubscriptionRepository,
	userRepo UserRepository,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:    orgRepo,
		memberRepo: memberRepo,
		orgSubRepo: orgSubRepo,
		userRepo:   userRepo,
	}
}

type CreateOrganizationInput struct {
	Name        string
	Slug        string
	Description *string
	OwnerUserID int64
	BillingMode string
	Balance     float64
	MonthlyBudgetUSD *float64
	MaxMembers  int
	MaxAPIKeys  int
}

type UpdateOrganizationInput struct {
	Name        string
	Slug        string
	Description *string
	BillingMode string
	MonthlyBudgetUSD *float64
	MaxMembers  int
	MaxAPIKeys  int
	Status      string
	AuditMode   string
}

var slugRegexp = regexp.MustCompile(`^[a-z0-9][a-z0-9\-]*[a-z0-9]$`)

func (s *OrganizationService) Create(ctx context.Context, input *CreateOrganizationInput) (*Organization, error) {
	if input.Slug != "" && !slugRegexp.MatchString(input.Slug) {
		return nil, fmt.Errorf("invalid slug format: must be lowercase alphanumeric with hyphens")
	}

	if input.BillingMode == "" {
		input.BillingMode = OrgBillingModeBalance
	}
	if input.MaxMembers == 0 {
		input.MaxMembers = 50
	}
	if input.MaxAPIKeys == 0 {
		input.MaxAPIKeys = 100
	}

	// Auto-generate slug from name if not provided
	if input.Slug == "" {
		input.Slug = strings.ToLower(strings.ReplaceAll(input.Name, " ", "-"))
	}

	org := &Organization{
		Name:             input.Name,
		Slug:             input.Slug,
		Description:      input.Description,
		OwnerUserID:      input.OwnerUserID,
		BillingMode:      input.BillingMode,
		Balance:          input.Balance,
		MonthlyBudgetUSD: input.MonthlyBudgetUSD,
		MaxMembers:       input.MaxMembers,
		MaxAPIKeys:       input.MaxAPIKeys,
		Status:           OrgStatusActive,
		AuditMode:        OrgAuditModeMetadata,
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	// Auto-add owner as org_admin member
	member := &OrgMember{
		OrgID:  org.ID,
		UserID: input.OwnerUserID,
		Role:   OrgMemberRoleAdmin,
		Status: StatusActive,
	}
	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	// Update owner's role to org_admin
	owner, err := s.userRepo.GetByID(ctx, input.OwnerUserID)
	if err != nil {
		return nil, err
	}
	if owner.Role == RoleUser {
		owner.Role = RoleOrgAdmin
		if err := s.userRepo.Update(ctx, owner); err != nil {
			return nil, err
		}
	}

	return org, nil
}

func (s *OrganizationService) GetByID(ctx context.Context, id int64) (*Organization, error) {
	return s.orgRepo.GetByID(ctx, id)
}

func (s *OrganizationService) GetBySlug(ctx context.Context, slug string) (*Organization, error) {
	return s.orgRepo.GetBySlug(ctx, slug)
}

func (s *OrganizationService) GetByOwnerID(ctx context.Context, ownerUserID int64) (*Organization, error) {
	return s.orgRepo.GetByOwnerID(ctx, ownerUserID)
}

func (s *OrganizationService) List(ctx context.Context, page, pageSize int, status, search string) ([]Organization, *pagination.PaginationResult, error) {
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	return s.orgRepo.List(ctx, params, status, search)
}

func (s *OrganizationService) Update(ctx context.Context, id int64, input *UpdateOrganizationInput) (*Organization, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	org.Name = input.Name
	org.Slug = input.Slug
	org.Description = input.Description
	org.BillingMode = input.BillingMode
	org.MonthlyBudgetUSD = input.MonthlyBudgetUSD
	org.MaxMembers = input.MaxMembers
	org.MaxAPIKeys = input.MaxAPIKeys
	if input.Status != "" {
		org.Status = input.Status
	}
	if input.AuditMode != "" {
		org.AuditMode = input.AuditMode
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *OrganizationService) Delete(ctx context.Context, id int64) error {
	return s.orgRepo.Delete(ctx, id)
}

type UpdateOrgBalanceInput struct {
	Amount float64
	Action string // "add" or "set"
}

func (s *OrganizationService) UpdateBalance(ctx context.Context, id int64, input *UpdateOrgBalanceInput) (*Organization, error) {
	if input.Action == "add" {
		if err := s.orgRepo.AddBalance(ctx, id, input.Amount); err != nil {
			return nil, err
		}
	} else {
		org, err := s.orgRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		diff := input.Amount - org.Balance
		if diff > 0 {
			if err := s.orgRepo.AddBalance(ctx, id, diff); err != nil {
				return nil, err
			}
		} else if diff < 0 {
			if err := s.orgRepo.DeductBalance(ctx, id, -diff); err != nil {
				return nil, err
			}
		}
	}
	return s.orgRepo.GetByID(ctx, id)
}

// SelfUpgrade allows a regular user to create their own organization
func (s *OrganizationService) SelfUpgrade(ctx context.Context, userID int64, input *CreateOrganizationInput) (*Organization, error) {
	// Check if user already owns an org
	_, err := s.orgRepo.GetByOwnerID(ctx, userID)
	if err == nil {
		return nil, fmt.Errorf("user already owns an organization")
	}

	input.OwnerUserID = userID
	return s.Create(ctx, input)
}

// GetDashboard returns org overview data
type OrgDashboard struct {
	Organization *Organization
	MemberCount  int
	APIKeyCount  int
}

func (s *OrganizationService) GetDashboard(ctx context.Context, orgID int64) (*OrgDashboard, error) {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	memberCount, err := s.orgRepo.CountMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	apiKeyCount, err := s.orgRepo.CountAPIKeys(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return &OrgDashboard{
		Organization: org,
		MemberCount:  memberCount,
		APIKeyCount:  apiKeyCount,
	}, nil
}
