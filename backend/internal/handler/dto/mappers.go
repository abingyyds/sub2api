// Package dto provides data transfer objects for HTTP handlers.
package dto

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

func UserFromServiceShallow(u *service.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		ID:              u.ID,
		Email:           u.Email,
		Username:        u.Username,
		Role:            u.Role,
		Balance:         u.Balance,
		Concurrency:     u.Concurrency,
		Status:          u.Status,
		AllowedGroups:   u.AllowedGroups,
		IsAgent:         u.IsAgent,
		AgentStatus:     u.AgentStatus,
		DiscoverySource: u.DiscoverySource,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func UserFromService(u *service.User) *User {
	if u == nil {
		return nil
	}
	out := UserFromServiceShallow(u)
	return out
}

// UserFromServiceAdmin converts a service User to DTO for admin users.
// It includes notes - user-facing endpoints must not use this.
func UserFromServiceAdmin(u *service.User) *AdminUser {
	if u == nil {
		return nil
	}
	base := UserFromService(u)
	if base == nil {
		return nil
	}
	return &AdminUser{
		User:         *base,
		Notes:        u.Notes,
		BoundSubSite: SubSiteFromService(u.BoundSubSite),
	}
}

func SubSiteFromService(s *service.SubSite) *SubSite {
	if s == nil {
		return nil
	}
	return &SubSite{
		ID:           s.ID,
		Name:         s.Name,
		Slug:         s.Slug,
		CustomDomain: s.CustomDomain,
		Status:       s.Status,
	}
}

func APIKeyFromService(k *service.APIKey) *APIKey {
	if k == nil {
		return nil
	}
	return &APIKey{
		ID:          k.ID,
		UserID:      k.UserID,
		Key:         k.Key,
		Name:        k.Name,
		GroupID:     k.GroupID,
		Status:      k.Status,
		IPWhitelist: k.IPWhitelist,
		IPBlacklist: k.IPBlacklist,
		UsageLimit:  k.UsageLimit,
		CreatedAt:   k.CreatedAt,
		UpdatedAt:   k.UpdatedAt,
		User:        UserFromServiceShallow(k.User),
		Group:       GroupFromServiceShallow(k.Group),
	}
}

func UserAPIKeyFromService(k *service.APIKey) *UserAPIKey {
	if k == nil {
		return nil
	}
	return &UserAPIKey{
		ID:          k.ID,
		UserID:      k.UserID,
		Key:         k.Key,
		Name:        k.Name,
		GroupID:     k.GroupID,
		Status:      k.Status,
		IPWhitelist: k.IPWhitelist,
		IPBlacklist: k.IPBlacklist,
		UsageLimit:  k.UsageLimit,
		CreatedAt:   k.CreatedAt,
		UpdatedAt:   k.UpdatedAt,
		User:        UserFromServiceShallow(k.User),
		Group:       UserGroupFromService(k.Group),
	}
}

func GroupFromServiceShallow(g *service.Group) *Group {
	if g == nil {
		return nil
	}
	out := groupFromServiceBase(g)
	return &out
}

func GroupFromService(g *service.Group) *Group {
	if g == nil {
		return nil
	}
	return GroupFromServiceShallow(g)
}

func UserGroupFromService(g *service.Group) *UserGroup {
	if g == nil {
		return nil
	}
	return &UserGroup{
		ID:                       g.ID,
		Name:                     g.Name,
		Description:              g.Description,
		Platform:                 g.Platform,
		IsExclusive:              g.IsExclusive,
		Status:                   g.Status,
		SubscriptionType:         g.SubscriptionType,
		DailyLimitUSD:            g.DailyLimitUSD,
		WeeklyLimitUSD:           g.WeeklyLimitUSD,
		MonthlyLimitUSD:          g.MonthlyLimitUSD,
		ClaudeCodeOnly:           g.ClaudeCodeOnly,
		FallbackGroupID:          g.FallbackGroupID,
		PriceFen:                 g.PriceFen,
		Listed:                   g.Listed,
		DefaultValidityDays:      g.DefaultValidityDays,
		PlanFeatures:             g.PlanFeatures,
		Tags:                     g.Tags,
		ModelPlazaVisible:        g.ModelPlazaVisible,
		DisplayPrice:             g.DisplayPrice,
		DisplayDiscount:          g.DisplayDiscount,
		QuotaPackageEnabled:      g.QuotaPackageEnabled,
		QuotaPackageQuotaUSD:     g.QuotaPackageQuotaUSD,
		QuotaPackageValidityDays: g.QuotaPackageValidityDays,
		CreatedAt:                g.CreatedAt,
		UpdatedAt:                g.UpdatedAt,
	}
}

// GroupFromServiceAdmin converts a service Group to DTO for admin users.
// It includes internal fields like model_routing and account_count.
func GroupFromServiceAdmin(g *service.Group) *AdminGroup {
	if g == nil {
		return nil
	}
	out := &AdminGroup{
		Group:                 groupFromServiceBase(g),
		DisplayRateMultiplier: g.DisplayRateMultiplier,
		ModelRouting:          g.ModelRouting,
		ModelRoutingEnabled:   g.ModelRoutingEnabled,
		AccountCount:          g.AccountCount,
	}
	if len(g.AccountGroups) > 0 {
		out.AccountGroups = make([]AccountGroup, 0, len(g.AccountGroups))
		for i := range g.AccountGroups {
			ag := g.AccountGroups[i]
			out.AccountGroups = append(out.AccountGroups, *AccountGroupFromService(&ag))
		}
	}
	return out
}

func groupFromServiceBase(g *service.Group) Group {
	return Group{
		ID:                       g.ID,
		Name:                     g.Name,
		Description:              g.Description,
		Platform:                 g.Platform,
		RateMultiplier:           g.RateMultiplier,
		IsExclusive:              g.IsExclusive,
		Status:                   g.Status,
		SubscriptionType:         g.SubscriptionType,
		DailyLimitUSD:            g.DailyLimitUSD,
		WeeklyLimitUSD:           g.WeeklyLimitUSD,
		MonthlyLimitUSD:          g.MonthlyLimitUSD,
		ImagePrice1K:             g.ImagePrice1K,
		ImagePrice2K:             g.ImagePrice2K,
		ImagePrice4K:             g.ImagePrice4K,
		ClaudeCodeOnly:           g.ClaudeCodeOnly,
		FallbackGroupID:          g.FallbackGroupID,
		PriceFen:                 g.PriceFen,
		Listed:                   g.Listed,
		DefaultValidityDays:      g.DefaultValidityDays,
		PlanFeatures:             g.PlanFeatures,
		Tags:                     g.Tags,
		ModelPlazaVisible:        g.ModelPlazaVisible,
		DisplayPrice:             g.DisplayPrice,
		DisplayDiscount:          g.DisplayDiscount,
		QuotaPackageEnabled:      g.QuotaPackageEnabled,
		QuotaPackageQuotaUSD:     g.QuotaPackageQuotaUSD,
		QuotaPackageValidityDays: g.QuotaPackageValidityDays,
		CreatedAt:                g.CreatedAt,
		UpdatedAt:                g.UpdatedAt,
	}
}

func AccountFromServiceShallow(a *service.Account) *Account {
	if a == nil {
		return nil
	}
	out := &Account{
		ID:                      a.ID,
		Name:                    a.Name,
		Notes:                   a.Notes,
		Platform:                a.Platform,
		Type:                    a.Type,
		Credentials:             a.Credentials,
		Extra:                   a.Extra,
		ProxyID:                 a.ProxyID,
		Concurrency:             a.Concurrency,
		Priority:                a.Priority,
		RateMultiplier:          a.BillingRateMultiplier(),
		Status:                  a.Status,
		ErrorMessage:            a.ErrorMessage,
		LastUsedAt:              a.LastUsedAt,
		ExpiresAt:               timeToUnixSeconds(a.ExpiresAt),
		AutoPauseOnExpired:      a.AutoPauseOnExpired,
		CreatedAt:               a.CreatedAt,
		UpdatedAt:               a.UpdatedAt,
		Schedulable:             a.Schedulable,
		RateLimitedAt:           a.RateLimitedAt,
		RateLimitResetAt:        a.RateLimitResetAt,
		OverloadUntil:           a.OverloadUntil,
		TempUnschedulableUntil:  a.TempUnschedulableUntil,
		TempUnschedulableReason: a.TempUnschedulableReason,
		SessionWindowStart:      a.SessionWindowStart,
		SessionWindowEnd:        a.SessionWindowEnd,
		SessionWindowStatus:     a.SessionWindowStatus,
		RequestBodyPassthrough:  a.IsRequestBodyPassthroughEnabled(),
		GroupIDs:                a.GroupIDs,
	}

	// 提取 5h 窗口费用控制和会话数量控制配置（仅 Anthropic OAuth/SetupToken 账号有效）
	if a.IsAnthropicOAuthOrSetupToken() {
		if limit := a.GetWindowCostLimit(); limit > 0 {
			out.WindowCostLimit = &limit
		}
		if reserve := a.GetWindowCostStickyReserve(); reserve > 0 {
			out.WindowCostStickyReserve = &reserve
		}
		if maxSessions := a.GetMaxSessions(); maxSessions > 0 {
			out.MaxSessions = &maxSessions
		}
		if idleTimeout := a.GetSessionIdleTimeoutMinutes(); idleTimeout > 0 {
			out.SessionIdleTimeoutMin = &idleTimeout
		}
		// TLS指纹伪装开关
		if a.IsTLSFingerprintEnabled() {
			enabled := true
			out.EnableTLSFingerprint = &enabled
		}
		// 会话ID伪装开关
		if a.IsSessionIDMaskingEnabled() {
			enabled := true
			out.EnableSessionIDMasking = &enabled
		}
	}

	return out
}

func AccountFromService(a *service.Account) *Account {
	if a == nil {
		return nil
	}
	out := AccountFromServiceShallow(a)
	out.Proxy = ProxyFromService(a.Proxy)
	if len(a.AccountGroups) > 0 {
		out.AccountGroups = make([]AccountGroup, 0, len(a.AccountGroups))
		for i := range a.AccountGroups {
			ag := a.AccountGroups[i]
			out.AccountGroups = append(out.AccountGroups, *AccountGroupFromService(&ag))
		}
	}
	if len(a.Groups) > 0 {
		out.Groups = make([]*Group, 0, len(a.Groups))
		for _, g := range a.Groups {
			out.Groups = append(out.Groups, GroupFromServiceShallow(g))
		}
	}
	return out
}

func timeToUnixSeconds(value *time.Time) *int64 {
	if value == nil {
		return nil
	}
	ts := value.Unix()
	return &ts
}

func AccountGroupFromService(ag *service.AccountGroup) *AccountGroup {
	if ag == nil {
		return nil
	}
	return &AccountGroup{
		AccountID: ag.AccountID,
		GroupID:   ag.GroupID,
		Priority:  ag.Priority,
		CreatedAt: ag.CreatedAt,
		Account:   AccountFromServiceShallow(ag.Account),
		Group:     GroupFromServiceShallow(ag.Group),
	}
}

func ProxyFromService(p *service.Proxy) *Proxy {
	if p == nil {
		return nil
	}
	return &Proxy{
		ID:        p.ID,
		Name:      p.Name,
		Protocol:  p.Protocol,
		Host:      p.Host,
		Port:      p.Port,
		Username:  p.Username,
		Password:  p.Password,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ProxyWithAccountCountFromService(p *service.ProxyWithAccountCount) *ProxyWithAccountCount {
	if p == nil {
		return nil
	}
	return &ProxyWithAccountCount{
		Proxy:          *ProxyFromService(&p.Proxy),
		AccountCount:   p.AccountCount,
		LatencyMs:      p.LatencyMs,
		LatencyStatus:  p.LatencyStatus,
		LatencyMessage: p.LatencyMessage,
		IPAddress:      p.IPAddress,
		Country:        p.Country,
		CountryCode:    p.CountryCode,
		Region:         p.Region,
		City:           p.City,
	}
}

func ProxyAccountSummaryFromService(a *service.ProxyAccountSummary) *ProxyAccountSummary {
	if a == nil {
		return nil
	}
	return &ProxyAccountSummary{
		ID:       a.ID,
		Name:     a.Name,
		Platform: a.Platform,
		Type:     a.Type,
		Notes:    a.Notes,
	}
}

func RedeemCodeFromService(rc *service.RedeemCode) *RedeemCode {
	if rc == nil {
		return nil
	}
	out := redeemCodeFromServiceBase(rc)
	return &out
}

func UserRedeemCodeFromService(rc *service.RedeemCode) *UserRedeemCode {
	if rc == nil {
		return nil
	}
	return &UserRedeemCode{
		ID:           rc.ID,
		Code:         rc.Code,
		Type:         rc.Type,
		Value:        rc.Value,
		Status:       rc.Status,
		UsedBy:       rc.UsedBy,
		UsedAt:       rc.UsedAt,
		CreatedAt:    rc.CreatedAt,
		GroupID:      rc.GroupID,
		ValidityDays: rc.ValidityDays,
		User:         UserFromServiceShallow(rc.User),
		Group:        UserGroupFromService(rc.Group),
	}
}

// RedeemCodeFromServiceAdmin converts a service RedeemCode to DTO for admin users.
// It includes notes - user-facing endpoints must not use this.
func RedeemCodeFromServiceAdmin(rc *service.RedeemCode) *AdminRedeemCode {
	if rc == nil {
		return nil
	}
	return &AdminRedeemCode{
		RedeemCode: redeemCodeFromServiceBase(rc),
		Notes:      rc.Notes,
	}
}

func redeemCodeFromServiceBase(rc *service.RedeemCode) RedeemCode {
	return RedeemCode{
		ID:           rc.ID,
		Code:         rc.Code,
		Type:         rc.Type,
		Value:        rc.Value,
		Status:       rc.Status,
		UsedBy:       rc.UsedBy,
		UsedAt:       rc.UsedAt,
		CreatedAt:    rc.CreatedAt,
		GroupID:      rc.GroupID,
		ValidityDays: rc.ValidityDays,
		User:         UserFromServiceShallow(rc.User),
		Group:        GroupFromServiceShallow(rc.Group),
	}
}

// AccountSummaryFromService returns a minimal AccountSummary for usage log display.
// Only includes ID and Name - no sensitive fields like Credentials, Proxy, etc.
func AccountSummaryFromService(a *service.Account) *AccountSummary {
	if a == nil {
		return nil
	}
	return &AccountSummary{
		ID:   a.ID,
		Name: a.Name,
	}
}

func usageLogAccountCost(l *service.UsageLog) float64 {
	if l == nil {
		return 0
	}
	multiplier := 1.0
	if l.AccountRateMultiplier != nil {
		multiplier = *l.AccountRateMultiplier
		if multiplier < 0 {
			multiplier = 1
		}
	}
	return l.TotalCost * multiplier
}

func usageLogFromService(l *service.UsageLog) UsageLog {
	// 基础调用日志 DTO 不包含计费倍率、成本拆分、真实成本和管理员字段。
	return UsageLog{
		ID:                    l.ID,
		UserID:                l.UserID,
		APIKeyID:              l.APIKeyID,
		AccountID:             l.AccountID,
		RequestID:             l.RequestID,
		Model:                 l.Model,
		GroupID:               l.GroupID,
		SubscriptionID:        l.SubscriptionID,
		InputTokens:           l.InputTokens,
		OutputTokens:          l.OutputTokens,
		CacheCreationTokens:   l.CacheCreationTokens,
		CacheReadTokens:       l.CacheReadTokens,
		CacheCreation5mTokens: l.CacheCreation5mTokens,
		CacheCreation1hTokens: l.CacheCreation1hTokens,
		ActualCost:            l.ActualCost,
		BillingType:           l.BillingType,
		Stream:                l.Stream,
		DurationMs:            l.DurationMs,
		FirstTokenMs:          l.FirstTokenMs,
		ImageCount:            l.ImageCount,
		ImageSize:             l.ImageSize,
		UserAgent:             l.UserAgent,
		CreatedAt:             l.CreatedAt,
		User:                  UserFromServiceShallow(l.User),
		APIKey:                UserAPIKeyFromService(l.APIKey),
		Group:                 UserGroupFromService(l.Group),
		Subscription:          UserSubscriptionFromService(l.Subscription),
	}
}

func usageLogFromServiceUser(l *service.UsageLog) UsageLog {
	return usageLogFromService(l)
}

// UsageLogFromService converts a service UsageLog to DTO for regular users.
// It excludes Account details and IP address - users should not see these.
func UsageLogFromService(l *service.UsageLog) *UsageLog {
	if l == nil {
		return nil
	}
	u := usageLogFromServiceUser(l)
	return &u
}

// UsageLogFromServiceAdmin converts a service UsageLog to DTO for admin users.
// It includes minimal Account info (ID, Name only) and IP address.
func UsageLogFromServiceAdmin(l *service.UsageLog) *AdminUsageLog {
	if l == nil {
		return nil
	}
	return &AdminUsageLog{
		UsageLog:              usageLogFromService(l),
		APIKey:                APIKeyFromService(l.APIKey),
		Group:                 GroupFromServiceShallow(l.Group),
		Subscription:          UserSubscriptionFromServiceAdmin(l.Subscription),
		InputCost:             l.InputCost,
		OutputCost:            l.OutputCost,
		CacheCreationCost:     l.CacheCreationCost,
		CacheReadCost:         l.CacheReadCost,
		TotalCost:             l.TotalCost,
		RateMultiplier:        l.RateMultiplier,
		AccountCost:           usageLogAccountCost(l),
		AccountRateMultiplier: l.AccountRateMultiplier,
		IPAddress:             l.IPAddress,
		Account:               AccountSummaryFromService(l.Account),
	}
}

func UserUsageStatsFromService(stats *service.UsageStats) *UserUsageStats {
	if stats == nil {
		return nil
	}
	return &UserUsageStats{
		TotalRequests:     stats.TotalRequests,
		TotalInputTokens:  stats.TotalInputTokens,
		TotalOutputTokens: stats.TotalOutputTokens,
		TotalCacheTokens:  stats.TotalCacheTokens,
		TotalTokens:       stats.TotalTokens,
		TotalActualCost:   stats.TotalActualCost,
		AverageDurationMs: stats.AverageDurationMs,
	}
}

func UserDashboardStatsFromUsageStats(stats *usagestats.UserDashboardStats) *UserDashboardStats {
	if stats == nil {
		return nil
	}
	return &UserDashboardStats{
		TotalAPIKeys:             stats.TotalAPIKeys,
		ActiveAPIKeys:            stats.ActiveAPIKeys,
		TotalRequests:            stats.TotalRequests,
		TotalInputTokens:         stats.TotalInputTokens,
		TotalOutputTokens:        stats.TotalOutputTokens,
		TotalCacheCreationTokens: stats.TotalCacheCreationTokens,
		TotalCacheReadTokens:     stats.TotalCacheReadTokens,
		TotalTokens:              stats.TotalTokens,
		TotalActualCost:          stats.TotalActualCost,
		TodayRequests:            stats.TodayRequests,
		TodayInputTokens:         stats.TodayInputTokens,
		TodayOutputTokens:        stats.TodayOutputTokens,
		TodayCacheCreationTokens: stats.TodayCacheCreationTokens,
		TodayCacheReadTokens:     stats.TodayCacheReadTokens,
		TodayTokens:              stats.TodayTokens,
		TodayActualCost:          stats.TodayActualCost,
		AverageDurationMs:        stats.AverageDurationMs,
		Rpm:                      stats.Rpm,
		Tpm:                      stats.Tpm,
	}
}

func UserTrendDataPointFromUsageStats(points []usagestats.TrendDataPoint) []UserTrendDataPoint {
	out := make([]UserTrendDataPoint, 0, len(points))
	for i := range points {
		point := points[i]
		out = append(out, UserTrendDataPoint{
			Date:         point.Date,
			Requests:     point.Requests,
			InputTokens:  point.InputTokens,
			OutputTokens: point.OutputTokens,
			CacheTokens:  point.CacheTokens,
			TotalTokens:  point.TotalTokens,
			ActualCost:   point.ActualCost,
		})
	}
	return out
}

func UserModelStatFromUsageStats(stats []usagestats.ModelStat) []UserModelStat {
	out := make([]UserModelStat, 0, len(stats))
	for i := range stats {
		stat := stats[i]
		out = append(out, UserModelStat{
			Model:        stat.Model,
			Requests:     stat.Requests,
			InputTokens:  stat.InputTokens,
			OutputTokens: stat.OutputTokens,
			TotalTokens:  stat.TotalTokens,
			ActualCost:   stat.ActualCost,
		})
	}
	return out
}

func UsageCleanupTaskFromService(task *service.UsageCleanupTask) *UsageCleanupTask {
	if task == nil {
		return nil
	}
	return &UsageCleanupTask{
		ID:     task.ID,
		Status: task.Status,
		Filters: UsageCleanupFilters{
			StartTime:   task.Filters.StartTime,
			EndTime:     task.Filters.EndTime,
			UserID:      task.Filters.UserID,
			APIKeyID:    task.Filters.APIKeyID,
			AccountID:   task.Filters.AccountID,
			GroupID:     task.Filters.GroupID,
			Model:       task.Filters.Model,
			Stream:      task.Filters.Stream,
			BillingType: task.Filters.BillingType,
		},
		CreatedBy:    task.CreatedBy,
		DeletedRows:  task.DeletedRows,
		ErrorMessage: task.ErrorMsg,
		CanceledBy:   task.CanceledBy,
		CanceledAt:   task.CanceledAt,
		StartedAt:    task.StartedAt,
		FinishedAt:   task.FinishedAt,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
	}
}

func SettingFromService(s *service.Setting) *Setting {
	if s == nil {
		return nil
	}
	return &Setting{
		ID:        s.ID,
		Key:       s.Key,
		Value:     s.Value,
		UpdatedAt: s.UpdatedAt,
	}
}

func UserSubscriptionFromService(sub *service.UserSubscription) *UserSubscription {
	if sub == nil {
		return nil
	}
	out := userSubscriptionFromServiceBase(sub)
	return &out
}

// UserSubscriptionFromServiceAdmin converts a service UserSubscription to DTO for admin users.
// It includes assignment metadata and notes.
func UserSubscriptionFromServiceAdmin(sub *service.UserSubscription) *AdminUserSubscription {
	if sub == nil {
		return nil
	}
	return &AdminUserSubscription{
		UserSubscription: userSubscriptionFromServiceBase(sub),
		Group:            GroupFromServiceShallow(sub.Group),
		AssignedBy:       sub.AssignedBy,
		AssignedAt:       sub.AssignedAt,
		Notes:            sub.Notes,
		AssignedByUser:   UserFromServiceShallow(sub.AssignedByUser),
	}
}

func userSubscriptionFromServiceBase(sub *service.UserSubscription) UserSubscription {
	return UserSubscription{
		ID:                 sub.ID,
		UserID:             sub.UserID,
		GroupID:            sub.GroupID,
		StartsAt:           sub.StartsAt,
		ExpiresAt:          sub.ExpiresAt,
		Status:             sub.Status,
		DailyWindowStart:   sub.DailyWindowStart,
		WeeklyWindowStart:  sub.WeeklyWindowStart,
		MonthlyWindowStart: sub.MonthlyWindowStart,
		DailyUsageUSD:      sub.DailyUsageUSD,
		WeeklyUsageUSD:     sub.WeeklyUsageUSD,
		MonthlyUsageUSD:    sub.MonthlyUsageUSD,
		CreatedAt:          sub.CreatedAt,
		UpdatedAt:          sub.UpdatedAt,
		User:               UserFromServiceShallow(sub.User),
		Group:              UserGroupFromService(sub.Group),
	}
}

func BulkAssignResultFromService(r *service.BulkAssignResult) *BulkAssignResult {
	if r == nil {
		return nil
	}
	subs := make([]AdminUserSubscription, 0, len(r.Subscriptions))
	for i := range r.Subscriptions {
		subs = append(subs, *UserSubscriptionFromServiceAdmin(&r.Subscriptions[i]))
	}
	return &BulkAssignResult{
		SuccessCount:  r.SuccessCount,
		FailedCount:   r.FailedCount,
		Subscriptions: subs,
		Errors:        r.Errors,
	}
}

func PromoCodeFromService(pc *service.PromoCode) *PromoCode {
	if pc == nil {
		return nil
	}
	return &PromoCode{
		ID:             pc.ID,
		Code:           pc.Code,
		DiscountAmount: pc.DiscountAmount,
		DiscountType:   pc.DiscountType,
		MinOrderAmount: pc.MinOrderAmount,
		MaxUses:        pc.MaxUses,
		UsedCount:      pc.UsedCount,
		Status:         pc.Status,
		ExpiresAt:      pc.ExpiresAt,
		Notes:          pc.Notes,
		CreatedAt:      pc.CreatedAt,
		UpdatedAt:      pc.UpdatedAt,
	}
}

func AnnouncementFromService(a *service.Announcement) *Announcement {
	if a == nil {
		return nil
	}
	return &Announcement{
		ID:          a.ID,
		Title:       a.Title,
		Content:     a.Content,
		Status:      a.Status,
		Priority:    a.Priority,
		Version:     a.Version,
		Category:    a.Category,
		PublishedAt: a.PublishedAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func PromoCodeUsageFromService(u *service.PromoCodeUsage) *PromoCodeUsage {
	if u == nil {
		return nil
	}
	return &PromoCodeUsage{
		ID:             u.ID,
		PromoCodeID:    u.PromoCodeID,
		UserID:         u.UserID,
		DiscountAmount: u.DiscountAmount,
		OrderNo:        u.OrderNo,
		UsedAt:         u.UsedAt,
		User:           UserFromServiceShallow(u.User),
	}
}
