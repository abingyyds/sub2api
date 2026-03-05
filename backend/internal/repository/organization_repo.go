package repository

import (
	"context"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/organization"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type organizationRepository struct {
	client *dbent.Client
}

func NewOrganizationRepository(client *dbent.Client) service.OrganizationRepository {
	return &organizationRepository{client: client}
}

func (r *organizationRepository) Create(ctx context.Context, org *service.Organization) error {
	client := clientFromContext(ctx, r.client)
	builder := client.Organization.Create().
		SetName(org.Name).
		SetSlug(org.Slug).
		SetOwnerUserID(org.OwnerUserID).
		SetBillingMode(org.BillingMode).
		SetBalance(org.Balance).
		SetMaxMembers(org.MaxMembers).
		SetMaxAPIKeys(org.MaxAPIKeys).
		SetStatus(org.Status).
		SetAuditMode(org.AuditMode)

	if org.Description != nil {
		builder.SetDescription(*org.Description)
	}
	if org.MonthlyBudgetUSD != nil {
		builder.SetMonthlyBudgetUsd(*org.MonthlyBudgetUSD)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, service.ErrOrganizationConflict)
	}
	applyOrganizationEntityToService(org, created)
	return nil
}

func (r *organizationRepository) GetByID(ctx context.Context, id int64) (*service.Organization, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.Organization.Query().
		Where(organization.IDEQ(id)).
		WithOwner().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrOrganizationNotFound, nil)
	}
	return organizationEntityToService(m), nil
}

func (r *organizationRepository) GetBySlug(ctx context.Context, slug string) (*service.Organization, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.Organization.Query().
		Where(organization.SlugEQ(slug)).
		WithOwner().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrOrganizationNotFound, nil)
	}
	return organizationEntityToService(m), nil
}

func (r *organizationRepository) GetByOwnerID(ctx context.Context, ownerUserID int64) (*service.Organization, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.Organization.Query().
		Where(organization.OwnerUserIDEQ(ownerUserID)).
		WithOwner().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrOrganizationNotFound, nil)
	}
	return organizationEntityToService(m), nil
}

func (r *organizationRepository) Update(ctx context.Context, org *service.Organization) error {
	client := clientFromContext(ctx, r.client)
	builder := client.Organization.UpdateOneID(org.ID).
		SetName(org.Name).
		SetSlug(org.Slug).
		SetBillingMode(org.BillingMode).
		SetMaxMembers(org.MaxMembers).
		SetMaxAPIKeys(org.MaxAPIKeys).
		SetStatus(org.Status).
		SetAuditMode(org.AuditMode)

	if org.Description != nil {
		builder.SetDescription(*org.Description)
	} else {
		builder.ClearDescription()
	}
	if org.MonthlyBudgetUSD != nil {
		builder.SetMonthlyBudgetUsd(*org.MonthlyBudgetUSD)
	} else {
		builder.ClearMonthlyBudgetUsd()
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrOrganizationNotFound, service.ErrOrganizationConflict)
	}
	applyOrganizationEntityToService(org, updated)
	return nil
}

func (r *organizationRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	err := client.Organization.DeleteOneID(id).Exec(ctx)
	return translatePersistenceError(err, service.ErrOrganizationNotFound, nil)
}

func (r *organizationRepository) List(ctx context.Context, params pagination.PaginationParams, status, search string) ([]service.Organization, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.Organization.Query()

	if status != "" {
		q = q.Where(organization.StatusEQ(status))
	}
	if search != "" {
		q = q.Where(
			organization.Or(
				organization.NameContains(search),
				organization.SlugContains(search),
			),
		)
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	orgs, err := q.
		WithOwner().
		Order(dbent.Desc(organization.FieldCreatedAt)).
		Offset(params.Offset()).
		Limit(params.Limit()).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	out := make([]service.Organization, 0, len(orgs))
	for _, m := range orgs {
		out = append(out, *organizationEntityToService(m))
	}
	return out, paginationResultFromTotal(int64(total), params), nil
}

func (r *organizationRepository) DeductBalance(ctx context.Context, orgID int64, amount float64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.Organization.UpdateOneID(orgID).
		AddBalance(-amount).
		Save(ctx)
	return translatePersistenceError(err, service.ErrOrganizationNotFound, nil)
}

func (r *organizationRepository) AddBalance(ctx context.Context, orgID int64, amount float64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.Organization.UpdateOneID(orgID).
		AddBalance(amount).
		Save(ctx)
	return translatePersistenceError(err, service.ErrOrganizationNotFound, nil)
}

func (r *organizationRepository) CountMembers(ctx context.Context, orgID int64) (int, error) {
	client := clientFromContext(ctx, r.client)
	org, err := client.Organization.Query().
		Where(organization.IDEQ(orgID)).
		WithMembers().
		Only(ctx)
	if err != nil {
		return 0, translatePersistenceError(err, service.ErrOrganizationNotFound, nil)
	}
	return len(org.Edges.Members), nil
}

func (r *organizationRepository) CountAPIKeys(ctx context.Context, orgID int64) (int, error) {
	client := clientFromContext(ctx, r.client)
	org, err := client.Organization.Query().
		Where(organization.IDEQ(orgID)).
		WithAPIKeys().
		Only(ctx)
	if err != nil {
		return 0, translatePersistenceError(err, service.ErrOrganizationNotFound, nil)
	}
	return len(org.Edges.APIKeys), nil
}

// Entity-to-service conversion helpers

func organizationEntityToService(m *dbent.Organization) *service.Organization {
	if m == nil {
		return nil
	}
	out := &service.Organization{
		ID:             m.ID,
		Name:           m.Name,
		Slug:           m.Slug,
		Description:    m.Description,
		OwnerUserID:    m.OwnerUserID,
		BillingMode:    m.BillingMode,
		Balance:        m.Balance,
		MonthlyBudgetUSD: m.MonthlyBudgetUsd,
		MaxMembers:     m.MaxMembers,
		MaxAPIKeys:     m.MaxAPIKeys,
		Status:         m.Status,
		AuditMode:      m.AuditMode,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	if m.Edges.Owner != nil {
		out.Owner = userEntityToService(m.Edges.Owner)
	}
	return out
}

func applyOrganizationEntityToService(org *service.Organization, m *dbent.Organization) {
	if m == nil {
		return
	}
	org.ID = m.ID
	org.CreatedAt = m.CreatedAt
	org.UpdatedAt = m.UpdatedAt
}
