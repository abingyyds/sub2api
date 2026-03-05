package repository

import (
	"context"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/orgmember"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type orgMemberRepository struct {
	client *dbent.Client
}

func NewOrgMemberRepository(client *dbent.Client) service.OrgMemberRepository {
	return &orgMemberRepository{client: client}
}

func (r *orgMemberRepository) Create(ctx context.Context, member *service.OrgMember) error {
	client := clientFromContext(ctx, r.client)
	builder := client.OrgMember.Create().
		SetOrgID(member.OrgID).
		SetUserID(member.UserID).
		SetRole(member.Role).
		SetMonthlyUsageUsd(member.MonthlyUsageUSD).
		SetDailyUsageUsd(member.DailyUsageUSD).
		SetStatus(member.Status)

	if member.MonthlyQuotaUSD != nil {
		builder.SetMonthlyQuotaUsd(*member.MonthlyQuotaUSD)
	}
	if member.DailyQuotaUSD != nil {
		builder.SetDailyQuotaUsd(*member.DailyQuotaUSD)
	}
	if member.Notes != nil {
		builder.SetNotes(*member.Notes)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, service.ErrOrgMemberConflict)
	}
	applyOrgMemberEntityToService(member, created)
	return nil
}

func (r *orgMemberRepository) GetByID(ctx context.Context, id int64) (*service.OrgMember, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.OrgMember.Query().
		Where(orgmember.IDEQ(id)).
		WithOrganization().
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrOrgMemberNotFound, nil)
	}
	return orgMemberEntityToService(m), nil
}

func (r *orgMemberRepository) GetByOrgAndUser(ctx context.Context, orgID, userID int64) (*service.OrgMember, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.OrgMember.Query().
		Where(
			orgmember.OrgIDEQ(orgID),
			orgmember.UserIDEQ(userID),
		).
		WithOrganization().
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrOrgMemberNotFound, nil)
	}
	return orgMemberEntityToService(m), nil
}

func (r *orgMemberRepository) Update(ctx context.Context, member *service.OrgMember) error {
	client := clientFromContext(ctx, r.client)
	builder := client.OrgMember.UpdateOneID(member.ID).
		SetRole(member.Role).
		SetStatus(member.Status)

	if member.MonthlyQuotaUSD != nil {
		builder.SetMonthlyQuotaUsd(*member.MonthlyQuotaUSD)
	} else {
		builder.ClearMonthlyQuotaUsd()
	}
	if member.DailyQuotaUSD != nil {
		builder.SetDailyQuotaUsd(*member.DailyQuotaUSD)
	} else {
		builder.ClearDailyQuotaUsd()
	}
	if member.Notes != nil {
		builder.SetNotes(*member.Notes)
	} else {
		builder.ClearNotes()
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrOrgMemberNotFound, nil)
	}
	applyOrgMemberEntityToService(member, updated)
	return nil
}

func (r *orgMemberRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	err := client.OrgMember.DeleteOneID(id).Exec(ctx)
	return translatePersistenceError(err, service.ErrOrgMemberNotFound, nil)
}

func (r *orgMemberRepository) ListByOrg(ctx context.Context, orgID int64, params pagination.PaginationParams) ([]service.OrgMember, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.OrgMember.Query().
		Where(orgmember.OrgIDEQ(orgID))

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	members, err := q.
		WithUser().
		WithOrganization().
		Order(dbent.Desc(orgmember.FieldCreatedAt)).
		Offset(params.Offset()).
		Limit(params.Limit()).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	out := make([]service.OrgMember, 0, len(members))
	for _, m := range members {
		out = append(out, *orgMemberEntityToService(m))
	}
	return out, paginationResultFromTotal(int64(total), params), nil
}

func (r *orgMemberRepository) IncrementUsage(ctx context.Context, id int64, costUSD float64) error {
	now := time.Now()
	client := clientFromContext(ctx, r.client)

	m, err := client.OrgMember.Get(ctx, id)
	if err != nil {
		return translatePersistenceError(err, service.ErrOrgMemberNotFound, nil)
	}

	builder := client.OrgMember.UpdateOneID(id)

	// Initialize or reset daily window
	if m.DailyWindowStart == nil || time.Since(*m.DailyWindowStart) >= 24*time.Hour {
		builder.SetDailyWindowStart(now).SetDailyUsageUsd(costUSD)
	} else {
		builder.AddDailyUsageUsd(costUSD)
	}

	// Initialize or reset monthly window
	if m.MonthlyWindowStart == nil || time.Since(*m.MonthlyWindowStart) >= 30*24*time.Hour {
		builder.SetMonthlyWindowStart(now).SetMonthlyUsageUsd(costUSD)
	} else {
		builder.AddMonthlyUsageUsd(costUSD)
	}

	_, err = builder.Save(ctx)
	return err
}

func (r *orgMemberRepository) ResetDailyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.OrgMember.UpdateOneID(id).
		SetDailyUsageUsd(0).
		SetDailyWindowStart(newWindowStart).
		Save(ctx)
	return translatePersistenceError(err, service.ErrOrgMemberNotFound, nil)
}

func (r *orgMemberRepository) ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.OrgMember.UpdateOneID(id).
		SetMonthlyUsageUsd(0).
		SetMonthlyWindowStart(newWindowStart).
		Save(ctx)
	return translatePersistenceError(err, service.ErrOrgMemberNotFound, nil)
}

// Entity-to-service conversion helpers

func orgMemberEntityToService(m *dbent.OrgMember) *service.OrgMember {
	if m == nil {
		return nil
	}
	out := &service.OrgMember{
		ID:                m.ID,
		OrgID:             m.OrgID,
		UserID:            m.UserID,
		Role:              m.Role,
		MonthlyQuotaUSD:   m.MonthlyQuotaUsd,
		DailyQuotaUSD:     m.DailyQuotaUsd,
		MonthlyUsageUSD:   m.MonthlyUsageUsd,
		DailyUsageUSD:     m.DailyUsageUsd,
		MonthlyWindowStart: m.MonthlyWindowStart,
		DailyWindowStart:   m.DailyWindowStart,
		Status:            m.Status,
		Notes:             m.Notes,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
	if m.Edges.Organization != nil {
		out.Organization = organizationEntityToService(m.Edges.Organization)
	}
	if m.Edges.User != nil {
		out.User = userEntityToService(m.Edges.User)
	}
	return out
}

func applyOrgMemberEntityToService(member *service.OrgMember, m *dbent.OrgMember) {
	if m == nil {
		return
	}
	member.ID = m.ID
	member.CreatedAt = m.CreatedAt
	member.UpdatedAt = m.UpdatedAt
}
