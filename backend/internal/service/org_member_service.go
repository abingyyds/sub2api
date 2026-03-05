package service

import (
	"context"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

type OrgMemberService struct {
	memberRepo OrgMemberRepository
	orgRepo    OrganizationRepository
	userRepo   UserRepository
}

func NewOrgMemberService(
	memberRepo OrgMemberRepository,
	orgRepo OrganizationRepository,
	userRepo UserRepository,
) *OrgMemberService {
	return &OrgMemberService{
		memberRepo: memberRepo,
		orgRepo:    orgRepo,
		userRepo:   userRepo,
	}
}

type CreateEmployeeInput struct {
	Email          string
	Password       string
	Username       string
	DailyQuotaUSD  *float64
	MonthlyQuotaUSD *float64
	Notes          *string
}

// CreateEmployee creates a new user account and adds them as a member to the organization
func (s *OrgMemberService) CreateEmployee(ctx context.Context, orgID int64, input *CreateEmployeeInput) (*OrgMember, error) {
	// Check org exists and is active
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if !org.IsActive() {
		return nil, ErrOrganizationDisabled
	}

	// Check member limit
	memberCount, err := s.orgRepo.CountMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if memberCount >= org.MaxMembers {
		return nil, ErrOrgMaxMembersReached
	}

	// Create user account
	user := &User{
		Email:  input.Email,
		Role:   RoleUser,
		Status: StatusActive,
	}
	if input.Username != "" {
		user.Username = input.Username
	}
	if err := user.SetPassword(input.Password); err != nil {
		return nil, err
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create org member
	member := &OrgMember{
		OrgID:           orgID,
		UserID:          user.ID,
		Role:            OrgMemberRoleMember,
		DailyQuotaUSD:   input.DailyQuotaUSD,
		MonthlyQuotaUSD: input.MonthlyQuotaUSD,
		Status:          StatusActive,
		Notes:           input.Notes,
	}
	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	member.User = user
	return member, nil
}

func (s *OrgMemberService) GetByID(ctx context.Context, id int64) (*OrgMember, error) {
	return s.memberRepo.GetByID(ctx, id)
}

func (s *OrgMemberService) GetByOrgAndUser(ctx context.Context, orgID, userID int64) (*OrgMember, error) {
	return s.memberRepo.GetByOrgAndUser(ctx, orgID, userID)
}

func (s *OrgMemberService) ListByOrg(ctx context.Context, orgID int64, page, pageSize int) ([]OrgMember, *pagination.PaginationResult, error) {
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	return s.memberRepo.ListByOrg(ctx, orgID, params)
}

type UpdateMemberInput struct {
	Role            string
	DailyQuotaUSD   *float64
	MonthlyQuotaUSD *float64
	Notes           *string
}

func (s *OrgMemberService) Update(ctx context.Context, id int64, input *UpdateMemberInput) (*OrgMember, error) {
	member, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Role != "" {
		member.Role = input.Role
	}
	member.DailyQuotaUSD = input.DailyQuotaUSD
	member.MonthlyQuotaUSD = input.MonthlyQuotaUSD
	if input.Notes != nil {
		member.Notes = input.Notes
	}

	if err := s.memberRepo.Update(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *OrgMemberService) Suspend(ctx context.Context, id int64) (*OrgMember, error) {
	member, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	member.Status = StatusDisabled
	if err := s.memberRepo.Update(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *OrgMemberService) Remove(ctx context.Context, id int64) error {
	return s.memberRepo.Delete(ctx, id)
}

// CheckMemberQuota checks if a member has enough quota for additional cost.
// Returns nil if OK, or an error if quota exceeded.
func (s *OrgMemberService) CheckMemberQuota(member *OrgMember, additionalCost float64) error {
	if !member.CheckDailyQuota(additionalCost) {
		return ErrOrgMemberDailyQuotaExceeded
	}
	if !member.CheckMonthlyQuota(additionalCost) {
		return ErrOrgMemberMonthlyQuotaExceeded
	}
	return nil
}

// IncrementMemberUsage atomically increments a member's daily and monthly usage.
func (s *OrgMemberService) IncrementMemberUsage(ctx context.Context, memberID int64, cost float64) error {
	return s.memberRepo.IncrementUsage(ctx, memberID, cost)
}
