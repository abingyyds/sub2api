package service

import (
	"context"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

type OrgProjectService struct {
	projectRepo OrgProjectRepository
	orgRepo     OrganizationRepository
}

func NewOrgProjectService(
	projectRepo OrgProjectRepository,
	orgRepo OrganizationRepository,
) *OrgProjectService {
	return &OrgProjectService{
		projectRepo: projectRepo,
		orgRepo:     orgRepo,
	}
}

type CreateOrgProjectInput struct {
	Name             string
	Description      *string
	GroupID          *int64
	AllowedModels    []string
	MonthlyBudgetUSD *float64
}

type UpdateOrgProjectInput struct {
	Name             string
	Description      *string
	GroupID          *int64
	AllowedModels    []string
	MonthlyBudgetUSD *float64
	Status           string
}

func (s *OrgProjectService) Create(ctx context.Context, orgID int64, input *CreateOrgProjectInput) (*OrgProject, error) {
	project := &OrgProject{
		OrgID:            orgID,
		Name:             input.Name,
		Description:      input.Description,
		GroupID:          input.GroupID,
		AllowedModels:    input.AllowedModels,
		MonthlyBudgetUSD: input.MonthlyBudgetUSD,
		Status:           OrgStatusActive,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *OrgProjectService) GetByID(ctx context.Context, id int64) (*OrgProject, error) {
	return s.projectRepo.GetByID(ctx, id)
}

func (s *OrgProjectService) ListByOrg(ctx context.Context, orgID int64, page, pageSize int) ([]OrgProject, *pagination.PaginationResult, error) {
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	return s.projectRepo.ListByOrg(ctx, orgID, params)
}

func (s *OrgProjectService) Update(ctx context.Context, id int64, input *UpdateOrgProjectInput) (*OrgProject, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	project.Name = input.Name
	project.Description = input.Description
	project.GroupID = input.GroupID
	project.AllowedModels = input.AllowedModels
	project.MonthlyBudgetUSD = input.MonthlyBudgetUSD
	if input.Status != "" {
		project.Status = input.Status
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *OrgProjectService) Delete(ctx context.Context, id int64) error {
	return s.projectRepo.Delete(ctx, id)
}

// IncrementUsage increments a project's monthly usage and checks budget.
func (s *OrgProjectService) IncrementUsage(ctx context.Context, projectID int64, costUSD float64) error {
	return s.projectRepo.IncrementUsage(ctx, projectID, costUSD)
}
