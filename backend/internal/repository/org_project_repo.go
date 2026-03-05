package repository

import (
	"context"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/orgproject"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type orgProjectRepository struct {
	client *dbent.Client
}

func NewOrgProjectRepository(client *dbent.Client) service.OrgProjectRepository {
	return &orgProjectRepository{client: client}
}

func (r *orgProjectRepository) Create(ctx context.Context, project *service.OrgProject) error {
	client := clientFromContext(ctx, r.client)
	builder := client.OrgProject.Create().
		SetOrgID(project.OrgID).
		SetName(project.Name).
		SetStatus(project.Status)

	if project.Description != nil {
		builder.SetNillableDescription(project.Description)
	}
	if project.GroupID != nil {
		builder.SetNillableGroupID(project.GroupID)
	}
	if project.AllowedModels != nil {
		builder.SetAllowedModels(project.AllowedModels)
	}
	if project.MonthlyBudgetUSD != nil {
		builder.SetNillableMonthlyBudgetUsd(project.MonthlyBudgetUSD)
	}

	entity, err := builder.Save(ctx)
	if err != nil {
		return err
	}
	project.ID = entity.ID
	project.CreatedAt = entity.CreatedAt
	project.UpdatedAt = entity.UpdatedAt
	return nil
}

func (r *orgProjectRepository) GetByID(ctx context.Context, id int64) (*service.OrgProject, error) {
	entity, err := r.client.OrgProject.Get(ctx, id)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrOrgProjectNotFound
		}
		return nil, err
	}
	return orgProjectEntityToService(entity), nil
}

func (r *orgProjectRepository) Update(ctx context.Context, project *service.OrgProject) error {
	client := clientFromContext(ctx, r.client)
	builder := client.OrgProject.UpdateOneID(project.ID).
		SetName(project.Name).
		SetStatus(project.Status).
		SetUpdatedAt(time.Now())

	if project.Description != nil {
		builder.SetNillableDescription(project.Description)
	} else {
		builder.ClearDescription()
	}
	if project.GroupID != nil {
		builder.SetNillableGroupID(project.GroupID)
	} else {
		builder.ClearGroupID()
	}
	if project.AllowedModels != nil {
		builder.SetAllowedModels(project.AllowedModels)
	} else {
		builder.ClearAllowedModels()
	}
	if project.MonthlyBudgetUSD != nil {
		builder.SetNillableMonthlyBudgetUsd(project.MonthlyBudgetUSD)
	} else {
		builder.ClearMonthlyBudgetUsd()
	}

	_, err := builder.Save(ctx)
	return err
}

func (r *orgProjectRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	return client.OrgProject.DeleteOneID(id).Exec(ctx)
}

func (r *orgProjectRepository) ListByOrg(ctx context.Context, orgID int64, params pagination.PaginationParams) ([]service.OrgProject, *pagination.PaginationResult, error) {
	query := r.client.OrgProject.Query().
		Where(orgproject.OrgID(orgID)).
		Order(dbent.Desc(orgproject.FieldCreatedAt))

	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	paginationResult := paginationResultFromTotal(total, params)
	entities, err := query.
		Offset(params.Offset()).
		Limit(params.Limit()).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	projects := make([]service.OrgProject, len(entities))
	for i, e := range entities {
		projects[i] = *orgProjectEntityToService(e)
	}
	return projects, paginationResult, nil
}

func (r *orgProjectRepository) IncrementUsage(ctx context.Context, id int64, costUSD float64) error {
	client := clientFromContext(ctx, r.client)
	project, err := client.OrgProject.Get(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	builder := client.OrgProject.UpdateOneID(id)

	// Reset monthly window if needed
	if project.MonthlyWindowStart == nil || now.Year() != project.MonthlyWindowStart.Year() || now.Month() != project.MonthlyWindowStart.Month() {
		builder.SetMonthlyUsageUsd(costUSD).SetMonthlyWindowStart(now)
	} else {
		builder.SetMonthlyUsageUsd(project.MonthlyUsageUsd + costUSD)
	}

	_, err = builder.Save(ctx)
	return err
}

func (r *orgProjectRepository) ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.OrgProject.UpdateOneID(id).
		SetMonthlyUsageUsd(0).
		SetMonthlyWindowStart(newWindowStart).
		Save(ctx)
	return err
}

func orgProjectEntityToService(e *dbent.OrgProject) *service.OrgProject {
	p := &service.OrgProject{
		ID:                 e.ID,
		OrgID:              e.OrgID,
		Name:               e.Name,
		Description:        e.Description,
		GroupID:             e.GroupID,
		AllowedModels:      e.AllowedModels,
		MonthlyBudgetUSD:   e.MonthlyBudgetUsd,
		MonthlyUsageUSD:    e.MonthlyUsageUsd,
		MonthlyWindowStart: e.MonthlyWindowStart,
		Status:             e.Status,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
	return p
}
