package dto

import "github.com/Wei-Shaw/sub2api/internal/service"

func OrganizationFromService(o *service.Organization) *Organization {
	if o == nil {
		return nil
	}
	out := &Organization{
		ID:               o.ID,
		Name:             o.Name,
		Slug:             o.Slug,
		Description:      o.Description,
		OwnerUserID:      o.OwnerUserID,
		BillingMode:      o.BillingMode,
		Balance:          o.Balance,
		MonthlyBudgetUSD: o.MonthlyBudgetUSD,
		MaxMembers:       o.MaxMembers,
		MaxAPIKeys:       o.MaxAPIKeys,
		Status:           o.Status,
		AuditMode:        o.AuditMode,
		CreatedAt:        o.CreatedAt,
		UpdatedAt:        o.UpdatedAt,
	}
	if o.Owner != nil {
		out.Owner = UserFromServiceShallow(o.Owner)
	}
	return out
}

func OrgMemberFromService(m *service.OrgMember) *OrgMember {
	if m == nil {
		return nil
	}
	out := &OrgMember{
		ID:              m.ID,
		OrgID:           m.OrgID,
		UserID:          m.UserID,
		Role:            m.Role,
		MonthlyQuotaUSD: m.MonthlyQuotaUSD,
		DailyQuotaUSD:   m.DailyQuotaUSD,
		MonthlyUsageUSD: m.MonthlyUsageUSD,
		DailyUsageUSD:   m.DailyUsageUSD,
		Status:          m.Status,
		Notes:           m.Notes,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
	if m.User != nil {
		out.User = UserFromServiceShallow(m.User)
	}
	return out
}

func OrgSubscriptionFromService(s *service.OrgSubscription) *OrgSubscription {
	if s == nil {
		return nil
	}
	out := &OrgSubscription{
		ID:              s.ID,
		OrgID:           s.OrgID,
		GroupID:         s.GroupID,
		StartsAt:        s.StartsAt,
		ExpiresAt:       s.ExpiresAt,
		Status:          s.Status,
		DailyUsageUSD:   s.DailyUsageUSD,
		WeeklyUsageUSD:  s.WeeklyUsageUSD,
		MonthlyUsageUSD: s.MonthlyUsageUSD,
		AssignedBy:      s.AssignedBy,
		AssignedAt:      s.AssignedAt,
		Notes:           s.Notes,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
	if s.Group != nil {
		out.Group = GroupFromServiceShallow(s.Group)
	}
	return out
}

func OrgDashboardFromService(d *service.OrgDashboard) *OrgDashboard {
	return &OrgDashboard{
		Organization: OrganizationFromService(d.Organization),
		MemberCount:  d.MemberCount,
		APIKeyCount:  d.APIKeyCount,
	}
}

func OrgProjectFromService(p *service.OrgProject) *OrgProject {
	if p == nil {
		return nil
	}
	return &OrgProject{
		ID:                 p.ID,
		OrgID:              p.OrgID,
		Name:               p.Name,
		Description:        p.Description,
		GroupID:             p.GroupID,
		AllowedModels:      p.AllowedModels,
		MonthlyBudgetUSD:   p.MonthlyBudgetUSD,
		MonthlyUsageUSD:    p.MonthlyUsageUSD,
		MonthlyWindowStart: p.MonthlyWindowStart,
		Status:             p.Status,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}

func OrgAuditLogFromService(l *service.OrgAuditLog) *OrgAuditLog {
	if l == nil {
		return nil
	}
	return &OrgAuditLog{
		ID:              l.ID,
		OrgID:           l.OrgID,
		UserID:          l.UserID,
		MemberID:        l.MemberID,
		ProjectID:       l.ProjectID,
		UsageLogID:      l.UsageLogID,
		Action:          l.Action,
		Model:           l.Model,
		AuditMode:       l.AuditMode,
		RequestSummary:  l.RequestSummary,
		RequestContent:  l.RequestContent,
		ResponseSummary: l.ResponseSummary,
		Keywords:        l.Keywords,
		Flagged:         l.Flagged,
		FlagReason:      l.FlagReason,
		InputTokens:     l.InputTokens,
		OutputTokens:    l.OutputTokens,
		CostUSD:         l.CostUSD,
		IPAddress:       l.IPAddress,
		UserAgent:       l.UserAgent,
		Detail:          l.Detail,
		CreatedAt:       l.CreatedAt,
	}
}
