package repository

import (
	"context"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/orgsubscription"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type orgSubscriptionRepository struct {
	client *dbent.Client
}

func NewOrgSubscriptionRepository(client *dbent.Client) service.OrgSubscriptionRepository {
	return &orgSubscriptionRepository{client: client}
}

func (r *orgSubscriptionRepository) Create(ctx context.Context, sub *service.OrgSubscription) error {
	client := clientFromContext(ctx, r.client)
	builder := client.OrgSubscription.Create().
		SetOrgID(sub.OrgID).
		SetGroupID(sub.GroupID).
		SetExpiresAt(sub.ExpiresAt).
		SetDailyUsageUsd(sub.DailyUsageUSD).
		SetWeeklyUsageUsd(sub.WeeklyUsageUSD).
		SetMonthlyUsageUsd(sub.MonthlyUsageUSD)

	if sub.StartsAt.IsZero() {
		builder.SetStartsAt(time.Now())
	} else {
		builder.SetStartsAt(sub.StartsAt)
	}
	if sub.Status != "" {
		builder.SetStatus(sub.Status)
	}
	if sub.AssignedBy != nil {
		builder.SetAssignedBy(*sub.AssignedBy)
	}
	if !sub.AssignedAt.IsZero() {
		builder.SetAssignedAt(sub.AssignedAt)
	}
	if sub.Notes != nil {
		builder.SetNotes(*sub.Notes)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, service.ErrOrgSubscriptionAlreadyExists)
	}
	applyOrgSubscriptionEntityToService(sub, created)
	return nil
}

func (r *orgSubscriptionRepository) GetByID(ctx context.Context, id int64) (*service.OrgSubscription, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.OrgSubscription.Query().
		Where(orgsubscription.IDEQ(id)).
		WithOrganization().
		WithGroup().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrOrgSubscriptionNotFound, nil)
	}
	return orgSubscriptionEntityToService(m), nil
}

func (r *orgSubscriptionRepository) GetActiveByOrgAndGroup(ctx context.Context, orgID, groupID int64) (*service.OrgSubscription, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.OrgSubscription.Query().
		Where(
			orgsubscription.OrgIDEQ(orgID),
			orgsubscription.GroupIDEQ(groupID),
			orgsubscription.StatusEQ(service.SubscriptionStatusActive),
			orgsubscription.ExpiresAtGT(time.Now()),
		).
		WithGroup().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrOrgSubscriptionNotFound, nil)
	}
	return orgSubscriptionEntityToService(m), nil
}

func (r *orgSubscriptionRepository) Update(ctx context.Context, sub *service.OrgSubscription) error {
	client := clientFromContext(ctx, r.client)
	builder := client.OrgSubscription.UpdateOneID(sub.ID).
		SetStatus(sub.Status).
		SetExpiresAt(sub.ExpiresAt)

	if sub.Notes != nil {
		builder.SetNotes(*sub.Notes)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrOrgSubscriptionNotFound, nil)
	}
	applyOrgSubscriptionEntityToService(sub, updated)
	return nil
}

func (r *orgSubscriptionRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	err := client.OrgSubscription.DeleteOneID(id).Exec(ctx)
	return translatePersistenceError(err, service.ErrOrgSubscriptionNotFound, nil)
}

func (r *orgSubscriptionRepository) ListByOrg(ctx context.Context, orgID int64, params pagination.PaginationParams) ([]service.OrgSubscription, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.OrgSubscription.Query().
		Where(orgsubscription.OrgIDEQ(orgID))

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	subs, err := q.
		WithOrganization().
		WithGroup().
		Order(dbent.Desc(orgsubscription.FieldCreatedAt)).
		Offset(params.Offset()).
		Limit(params.Limit()).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	out := make([]service.OrgSubscription, 0, len(subs))
	for _, m := range subs {
		out = append(out, *orgSubscriptionEntityToService(m))
	}
	return out, paginationResultFromTotal(int64(total), params), nil
}

func (r *orgSubscriptionRepository) IncrementUsage(ctx context.Context, id int64, costUSD float64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.OrgSubscription.UpdateOneID(id).
		AddDailyUsageUsd(costUSD).
		AddWeeklyUsageUsd(costUSD).
		AddMonthlyUsageUsd(costUSD).
		Save(ctx)
	return translatePersistenceError(err, service.ErrOrgSubscriptionNotFound, nil)
}

func (r *orgSubscriptionRepository) ResetDailyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.OrgSubscription.UpdateOneID(id).
		SetDailyUsageUsd(0).
		SetDailyWindowStart(newWindowStart).
		Save(ctx)
	return translatePersistenceError(err, service.ErrOrgSubscriptionNotFound, nil)
}

func (r *orgSubscriptionRepository) ResetWeeklyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.OrgSubscription.UpdateOneID(id).
		SetWeeklyUsageUsd(0).
		SetWeeklyWindowStart(newWindowStart).
		Save(ctx)
	return translatePersistenceError(err, service.ErrOrgSubscriptionNotFound, nil)
}

func (r *orgSubscriptionRepository) ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.OrgSubscription.UpdateOneID(id).
		SetMonthlyUsageUsd(0).
		SetMonthlyWindowStart(newWindowStart).
		Save(ctx)
	return translatePersistenceError(err, service.ErrOrgSubscriptionNotFound, nil)
}

// Entity-to-service conversion helpers

func orgSubscriptionEntityToService(m *dbent.OrgSubscription) *service.OrgSubscription {
	if m == nil {
		return nil
	}
	out := &service.OrgSubscription{
		ID:                 m.ID,
		OrgID:              m.OrgID,
		GroupID:            m.GroupID,
		StartsAt:           m.StartsAt,
		ExpiresAt:          m.ExpiresAt,
		Status:             m.Status,
		DailyWindowStart:   m.DailyWindowStart,
		WeeklyWindowStart:  m.WeeklyWindowStart,
		MonthlyWindowStart: m.MonthlyWindowStart,
		DailyUsageUSD:      m.DailyUsageUsd,
		WeeklyUsageUSD:     m.WeeklyUsageUsd,
		MonthlyUsageUSD:    m.MonthlyUsageUsd,
		AssignedBy:         m.AssignedBy,
		AssignedAt:         m.AssignedAt,
		Notes:              m.Notes,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
	if m.Edges.Organization != nil {
		out.Organization = organizationEntityToService(m.Edges.Organization)
	}
	if m.Edges.Group != nil {
		out.Group = groupEntityToService(m.Edges.Group)
	}
	return out
}

func applyOrgSubscriptionEntityToService(sub *service.OrgSubscription, m *dbent.OrgSubscription) {
	if m == nil {
		return
	}
	sub.ID = m.ID
	sub.CreatedAt = m.CreatedAt
	sub.UpdatedAt = m.UpdatedAt
}
