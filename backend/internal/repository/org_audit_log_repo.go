package repository

import (
	"context"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/orgauditlog"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type orgAuditLogRepository struct {
	client *dbent.Client
}

func NewOrgAuditLogRepository(client *dbent.Client) service.OrgAuditLogRepository {
	return &orgAuditLogRepository{client: client}
}

func (r *orgAuditLogRepository) Create(ctx context.Context, auditLog *service.OrgAuditLog) error {
	builder := r.client.OrgAuditLog.Create().
		SetOrgID(auditLog.OrgID).
		SetUserID(auditLog.UserID).
		SetAction(auditLog.Action).
		SetAuditMode(auditLog.AuditMode).
		SetFlagged(auditLog.Flagged)

	if auditLog.MemberID != nil {
		builder.SetNillableMemberID(auditLog.MemberID)
	}
	if auditLog.ProjectID != nil {
		builder.SetNillableProjectID(auditLog.ProjectID)
	}
	if auditLog.UsageLogID != nil {
		builder.SetNillableUsageLogID(auditLog.UsageLogID)
	}
	if auditLog.Model != nil {
		builder.SetNillableModel(auditLog.Model)
	}
	if auditLog.RequestSummary != nil {
		builder.SetNillableRequestSummary(auditLog.RequestSummary)
	}
	if auditLog.RequestContent != nil {
		builder.SetNillableRequestContent(auditLog.RequestContent)
	}
	if auditLog.ResponseSummary != nil {
		builder.SetNillableResponseSummary(auditLog.ResponseSummary)
	}
	if auditLog.Keywords != nil {
		builder.SetKeywords(auditLog.Keywords)
	}
	if auditLog.FlagReason != nil {
		builder.SetNillableFlagReason(auditLog.FlagReason)
	}
	if auditLog.InputTokens != nil {
		builder.SetNillableInputTokens(auditLog.InputTokens)
	}
	if auditLog.OutputTokens != nil {
		builder.SetNillableOutputTokens(auditLog.OutputTokens)
	}
	if auditLog.CostUSD != nil {
		builder.SetNillableCostUsd(auditLog.CostUSD)
	}
	if auditLog.IPAddress != nil {
		builder.SetNillableIPAddress(auditLog.IPAddress)
	}
	if auditLog.UserAgent != nil {
		builder.SetNillableUserAgent(auditLog.UserAgent)
	}
	if auditLog.Detail != nil {
		builder.SetDetail(auditLog.Detail)
	}

	entity, err := builder.Save(ctx)
	if err != nil {
		return err
	}
	auditLog.ID = entity.ID
	auditLog.CreatedAt = entity.CreatedAt
	return nil
}

func (r *orgAuditLogRepository) GetByID(ctx context.Context, id int64) (*service.OrgAuditLog, error) {
	entity, err := r.client.OrgAuditLog.Get(ctx, id)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return orgAuditLogEntityToService(entity), nil
}

func (r *orgAuditLogRepository) List(ctx context.Context, orgID int64, params pagination.PaginationParams, filters service.AuditLogFilters) ([]service.OrgAuditLog, *pagination.PaginationResult, error) {
	query := r.client.OrgAuditLog.Query().
		Where(orgauditlog.OrgID(orgID)).
		Order(dbent.Desc(orgauditlog.FieldCreatedAt))

	if filters.MemberID != nil {
		query.Where(orgauditlog.MemberID(*filters.MemberID))
	}
	if filters.ProjectID != nil {
		query.Where(orgauditlog.ProjectID(*filters.ProjectID))
	}
	if filters.Action != "" {
		query.Where(orgauditlog.Action(filters.Action))
	}
	if filters.Model != "" {
		query.Where(orgauditlog.Model(filters.Model))
	}
	if filters.Flagged != nil {
		query.Where(orgauditlog.Flagged(*filters.Flagged))
	}
	if filters.StartDate != nil {
		query.Where(orgauditlog.CreatedAtGTE(*filters.StartDate))
	}
	if filters.EndDate != nil {
		query.Where(orgauditlog.CreatedAtLTE(*filters.EndDate))
	}

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

	logs := make([]service.OrgAuditLog, len(entities))
	for i, e := range entities {
		logs[i] = *orgAuditLogEntityToService(e)
	}
	return logs, paginationResult, nil
}

func (r *orgAuditLogRepository) CountFlagged(ctx context.Context, orgID int64) (int, error) {
	return r.client.OrgAuditLog.Query().
		Where(orgauditlog.OrgID(orgID), orgauditlog.Flagged(true)).
		Count(ctx)
}

func orgAuditLogEntityToService(e *dbent.OrgAuditLog) *service.OrgAuditLog {
	return &service.OrgAuditLog{
		ID:              e.ID,
		OrgID:           e.OrgID,
		UserID:          e.UserID,
		MemberID:        e.MemberID,
		ProjectID:       e.ProjectID,
		UsageLogID:      e.UsageLogID,
		Action:          e.Action,
		Model:           e.Model,
		AuditMode:       e.AuditMode,
		RequestSummary:  e.RequestSummary,
		RequestContent:  e.RequestContent,
		ResponseSummary: e.ResponseSummary,
		Keywords:        e.Keywords,
		Flagged:         e.Flagged,
		FlagReason:      e.FlagReason,
		InputTokens:     e.InputTokens,
		OutputTokens:    e.OutputTokens,
		CostUSD:         e.CostUsd,
		IPAddress:       e.IPAddress,
		UserAgent:       e.UserAgent,
		Detail:          e.Detail,
		CreatedAt:       e.CreatedAt,
	}
}
