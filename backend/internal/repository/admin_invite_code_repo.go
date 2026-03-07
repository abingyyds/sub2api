package repository

import (
	"context"
	"fmt"

	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/admininvitecode"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type adminInviteCodeRepo struct {
	client *ent.Client
}

func NewAdminInviteCodeRepo(client *ent.Client) service.AdminInviteCodeRepository {
	return &adminInviteCodeRepo{client: client}
}

func (r *adminInviteCodeRepo) Create(ctx context.Context, code *service.AdminInviteCode) error {
	created, err := r.client.AdminInviteCode.Create().
		SetCode(code.Code).
		SetSourceName(code.SourceName).
		SetCreatedBy(code.CreatedBy).
		SetNillableMaxUses(code.MaxUses).
		SetEnabled(code.Enabled).
		SetNillableNotes(&code.Notes).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return service.ErrAdminInviteCodeExists
		}
		return fmt.Errorf("create admin invite code: %w", err)
	}
	code.ID = created.ID
	code.CreatedAt = created.CreatedAt
	code.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *adminInviteCodeRepo) GetByID(ctx context.Context, id int64) (*service.AdminInviteCode, error) {
	code, err := r.client.AdminInviteCode.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, service.ErrAdminInviteCodeNotFound
		}
		return nil, fmt.Errorf("get admin invite code: %w", err)
	}
	return toServiceAdminInviteCode(code), nil
}

func (r *adminInviteCodeRepo) GetByCode(ctx context.Context, code string) (*service.AdminInviteCode, error) {
	result, err := r.client.AdminInviteCode.Query().
		Where(admininvitecode.CodeEQ(code)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, service.ErrAdminInviteCodeNotFound
		}
		return nil, fmt.Errorf("get admin invite code by code: %w", err)
	}
	return toServiceAdminInviteCode(result), nil
}

func (r *adminInviteCodeRepo) Update(ctx context.Context, code *service.AdminInviteCode) error {
	_, err := r.client.AdminInviteCode.UpdateOneID(code.ID).
		SetSourceName(code.SourceName).
		SetNillableMaxUses(code.MaxUses).
		SetEnabled(code.Enabled).
		SetNillableNotes(&code.Notes).
		Save(ctx)
	return err
}

func (r *adminInviteCodeRepo) Delete(ctx context.Context, id int64) error {
	return r.client.AdminInviteCode.DeleteOneID(id).Exec(ctx)
}

func (r *adminInviteCodeRepo) List(ctx context.Context, params pagination.PaginationParams) ([]service.AdminInviteCode, *pagination.PaginationResult, error) {
	query := r.client.AdminInviteCode.Query()
	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("count admin invite codes: %w", err)
	}
	codes, err := query.
		Limit(params.PageSize).
		Offset((params.Page - 1) * params.PageSize).
		Order(ent.Desc(admininvitecode.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list admin invite codes: %w", err)
	}
	result := make([]service.AdminInviteCode, len(codes))
	for i, c := range codes {
		result[i] = *toServiceAdminInviteCode(c)
	}
	return result, &pagination.PaginationResult{
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		Pages:    (total + params.PageSize - 1) / params.PageSize,
	}, nil
}

func (r *adminInviteCodeRepo) IncrementUsedCount(ctx context.Context, id int64) error {
	_, err := r.client.AdminInviteCode.UpdateOneID(id).
		AddUsedCount(1).
		Save(ctx)
	return err
}

func toServiceAdminInviteCode(e *ent.AdminInviteCode) *service.AdminInviteCode {
	return &service.AdminInviteCode{
		ID:         e.ID,
		Code:       e.Code,
		SourceName: e.SourceName,
		CreatedBy:  e.CreatedBy,
		UsedCount:  e.UsedCount,
		MaxUses:    e.MaxUses,
		Enabled:    e.Enabled,
		Notes:      e.Notes,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}
