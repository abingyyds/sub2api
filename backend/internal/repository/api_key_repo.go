package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/apikey"
	"github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
	"github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

type apiKeyRepository struct {
	client *dbent.Client
	sql    sqlExecutor
}

func NewAPIKeyRepository(client *dbent.Client) service.APIKeyRepository {
	return &apiKeyRepository{client: client}
}

func NewAPIKeyRepositoryWithSQL(client *dbent.Client, sqlDB *sql.DB) service.APIKeyRepository {
	return &apiKeyRepository{client: client, sql: sqlDB}
}

func (r *apiKeyRepository) activeQuery() *dbent.APIKeyQuery {
	// 默认过滤已软删除记录，避免删除后仍被查询到。
	return r.client.APIKey.Query().Where(apikey.DeletedAtIsNil())
}

func (r *apiKeyRepository) Create(ctx context.Context, key *service.APIKey) error {
	builder := r.client.APIKey.Create().
		SetUserID(key.UserID).
		SetKey(key.Key).
		SetName(key.Name).
		SetStatus(key.Status).
		SetNillableGroupID(key.GroupID).
		SetNillableUsageLimit(key.UsageLimit)

	if len(key.IPWhitelist) > 0 {
		builder.SetIPWhitelist(key.IPWhitelist)
	}
	if len(key.IPBlacklist) > 0 {
		builder.SetIPBlacklist(key.IPBlacklist)
	}

	created, err := builder.Save(ctx)
	if err == nil {
		key.ID = created.ID
		key.CreatedAt = created.CreatedAt
		key.UpdatedAt = created.UpdatedAt
	}
	return translatePersistenceError(err, nil, service.ErrAPIKeyExists)
}

func (r *apiKeyRepository) GetByID(ctx context.Context, id int64) (*service.APIKey, error) {
	m, err := r.activeQuery().
		Where(apikey.IDEQ(id)).
		WithUser().
		WithGroup().
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrAPIKeyNotFound
		}
		return nil, err
	}
	out := apiKeyEntityToService(m)
	if err := r.hydrateAPIKeyGroupQuotaFields(ctx, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetKeyAndOwnerID 根据 API Key ID 获取其 key 与所有者（用户）ID。
// 相比 GetByID，此方法性能更优，因为：
//   - 使用 Select() 只查询必要字段，减少数据传输量
//   - 不加载完整的 API Key 实体及其关联数据（User、Group 等）
//   - 适用于删除等只需 key 与用户 ID 的场景
func (r *apiKeyRepository) GetKeyAndOwnerID(ctx context.Context, id int64) (string, int64, error) {
	m, err := r.activeQuery().
		Where(apikey.IDEQ(id)).
		Select(apikey.FieldKey, apikey.FieldUserID).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return "", 0, service.ErrAPIKeyNotFound
		}
		return "", 0, err
	}
	return m.Key, m.UserID, nil
}

func (r *apiKeyRepository) GetByKey(ctx context.Context, key string) (*service.APIKey, error) {
	m, err := r.activeQuery().
		Where(apikey.KeyEQ(key)).
		WithUser().
		WithGroup().
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrAPIKeyNotFound
		}
		return nil, err
	}
	out := apiKeyEntityToService(m)
	if err := r.hydrateAPIKeyGroupQuotaFields(ctx, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *apiKeyRepository) GetByKeyForAuth(ctx context.Context, key string) (*service.APIKey, error) {
	if r.sql != nil {
		return r.getByKeyForAuthSQL(ctx, key)
	}
	m, err := r.activeQuery().
		Where(apikey.KeyEQ(key)).
		Select(
			apikey.FieldID,
			apikey.FieldUserID,
			apikey.FieldGroupID,
			apikey.FieldStatus,
			apikey.FieldIPWhitelist,
			apikey.FieldIPBlacklist,
		).
		WithUser(func(q *dbent.UserQuery) {
			q.Select(
				user.FieldID,
				user.FieldStatus,
				user.FieldRole,
				user.FieldBalance,
				user.FieldConcurrency,
			)
		}).
		WithGroup(func(q *dbent.GroupQuery) {
			q.Select(
				group.FieldID,
				group.FieldName,
				group.FieldPlatform,
				group.FieldStatus,
				group.FieldSubscriptionType,
				group.FieldRateMultiplier,
				group.FieldDailyLimitUsd,
				group.FieldWeeklyLimitUsd,
				group.FieldMonthlyLimitUsd,
				group.FieldImagePrice1k,
				group.FieldImagePrice2k,
				group.FieldImagePrice4k,
				group.FieldClaudeCodeOnly,
				group.FieldFallbackGroupID,
				group.FieldModelRoutingEnabled,
				group.FieldModelRouting,
			)
		}).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrAPIKeyNotFound
		}
		return nil, err
	}
	out := apiKeyEntityToService(m)
	if err := r.hydrateAPIKeyGroupQuotaFields(ctx, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *apiKeyRepository) getByKeyForAuthSQL(ctx context.Context, key string) (*service.APIKey, error) {
	query := `
		SELECT
			ak.id,
			ak.user_id,
			ak.group_id,
			ak.status,
			ak.ip_whitelist,
			ak.ip_blacklist,
			ak.org_id,
			ak.org_project_id,
			u.id,
			u.status,
			u.role,
			u.balance,
			u.concurrency,
			CASE
				WHEN ula.user_id IS NOT NULL
				 AND ula.terms_version = $2
				 AND ula.privacy_version = $3
				 AND ula.api_terms_version = $4
				 AND ula.terms_accepted_at IS NOT NULL
				 AND ula.privacy_accepted_at IS NOT NULL
				 AND ula.api_terms_accepted_at IS NOT NULL
				THEN TRUE ELSE FALSE
			END,
			g.id,
			g.name,
			g.platform,
			g.status,
			g.subscription_type,
			g.rate_multiplier,
			g.daily_limit_usd,
			g.weekly_limit_usd,
			g.monthly_limit_usd,
			g.image_price_1k,
			g.image_price_2k,
			g.image_price_4k,
			g.claude_code_only,
			g.fallback_group_id,
			g.model_routing_enabled,
			g.model_routing,
			COALESCE(g.quota_package_enabled, FALSE),
			g.quota_package_quota_usd,
			COALESCE(NULLIF(g.quota_package_validity_days, 0), 30)
		FROM api_keys ak
		JOIN users u ON u.id = ak.user_id
		LEFT JOIN user_legal_agreements ula ON ula.user_id = u.id
		LEFT JOIN groups g ON g.id = ak.group_id
		WHERE ak.key = $1
		  AND ak.deleted_at IS NULL
		LIMIT 1
	`

	var keyOut service.APIKey
	var groupID sql.NullInt64
	var orgID sql.NullInt64
	var orgProjectID sql.NullInt64
	var ipWhitelistJSON sql.NullString
	var ipBlacklistJSON sql.NullString
	var userOut service.User
	var legalAccepted bool
	var groupIDValue sql.NullInt64
	var groupName sql.NullString
	var groupPlatform sql.NullString
	var groupStatus sql.NullString
	var groupSubscriptionType sql.NullString
	var groupRateMultiplier sql.NullFloat64
	var dailyLimit sql.NullFloat64
	var weeklyLimit sql.NullFloat64
	var monthlyLimit sql.NullFloat64
	var imagePrice1K sql.NullFloat64
	var imagePrice2K sql.NullFloat64
	var imagePrice4K sql.NullFloat64
	var claudeCodeOnly sql.NullBool
	var fallbackGroupID sql.NullInt64
	var modelRoutingEnabled sql.NullBool
	var modelRoutingJSON sql.NullString
	var quotaPackageEnabled sql.NullBool
	var quotaPackageQuota sql.NullFloat64
	var quotaPackageValidityDays sql.NullInt64

	err := scanSingleRow(
		ctx,
		r.sql,
		query,
		[]any{key, service.LegalTermsVersion, service.LegalPrivacyVersion, service.LegalApiTermsVersion},
		&keyOut.ID,
		&keyOut.UserID,
		&groupID,
		&keyOut.Status,
		&ipWhitelistJSON,
		&ipBlacklistJSON,
		&orgID,
		&orgProjectID,
		&userOut.ID,
		&userOut.Status,
		&userOut.Role,
		&userOut.Balance,
		&userOut.Concurrency,
		&legalAccepted,
		&groupIDValue,
		&groupName,
		&groupPlatform,
		&groupStatus,
		&groupSubscriptionType,
		&groupRateMultiplier,
		&dailyLimit,
		&weeklyLimit,
		&monthlyLimit,
		&imagePrice1K,
		&imagePrice2K,
		&imagePrice4K,
		&claudeCodeOnly,
		&fallbackGroupID,
		&modelRoutingEnabled,
		&modelRoutingJSON,
		&quotaPackageEnabled,
		&quotaPackageQuota,
		&quotaPackageValidityDays,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrAPIKeyNotFound
		}
		return nil, err
	}

	keyOut.Key = key
	var ipWhitelist []string
	var ipBlacklist []string
	if ipWhitelistJSON.Valid && ipWhitelistJSON.String != "" {
		if err := json.Unmarshal([]byte(ipWhitelistJSON.String), &ipWhitelist); err != nil {
			return nil, fmt.Errorf("decode api key ip whitelist: %w", err)
		}
	}
	if ipBlacklistJSON.Valid && ipBlacklistJSON.String != "" {
		if err := json.Unmarshal([]byte(ipBlacklistJSON.String), &ipBlacklist); err != nil {
			return nil, fmt.Errorf("decode api key ip blacklist: %w", err)
		}
	}
	keyOut.IPWhitelist = ipWhitelist
	keyOut.IPBlacklist = ipBlacklist
	if groupID.Valid {
		keyOut.GroupID = &groupID.Int64
	}
	if orgID.Valid {
		keyOut.OrgID = &orgID.Int64
	}
	if orgProjectID.Valid {
		keyOut.OrgProjectID = &orgProjectID.Int64
	}
	userOut.LegalAgreementAccepted = legalAccepted
	keyOut.User = &userOut

	if groupIDValue.Valid {
		var modelRouting map[string][]int64
		if modelRoutingJSON.Valid && modelRoutingJSON.String != "" {
			if err := json.Unmarshal([]byte(modelRoutingJSON.String), &modelRouting); err != nil {
				return nil, fmt.Errorf("decode group model routing: %w", err)
			}
		}
		groupOut := &service.Group{
			ID:                  groupIDValue.Int64,
			Name:                groupName.String,
			Platform:            groupPlatform.String,
			Status:              groupStatus.String,
			Hydrated:            true,
			SubscriptionType:    groupSubscriptionType.String,
			ModelRouting:        modelRouting,
			ModelRoutingEnabled: modelRoutingEnabled.Bool,
			ClaudeCodeOnly:      claudeCodeOnly.Bool,
		}
		if groupRateMultiplier.Valid {
			groupOut.RateMultiplier = groupRateMultiplier.Float64
		}
		if dailyLimit.Valid {
			v := dailyLimit.Float64
			groupOut.DailyLimitUSD = &v
		}
		if weeklyLimit.Valid {
			v := weeklyLimit.Float64
			groupOut.WeeklyLimitUSD = &v
		}
		if monthlyLimit.Valid {
			v := monthlyLimit.Float64
			groupOut.MonthlyLimitUSD = &v
		}
		if imagePrice1K.Valid {
			v := imagePrice1K.Float64
			groupOut.ImagePrice1K = &v
		}
		if imagePrice2K.Valid {
			v := imagePrice2K.Float64
			groupOut.ImagePrice2K = &v
		}
		if imagePrice4K.Valid {
			v := imagePrice4K.Float64
			groupOut.ImagePrice4K = &v
		}
		if fallbackGroupID.Valid {
			v := fallbackGroupID.Int64
			groupOut.FallbackGroupID = &v
		}
		groupOut.QuotaPackageEnabled = quotaPackageEnabled.Bool
		if quotaPackageQuota.Valid {
			v := quotaPackageQuota.Float64
			groupOut.QuotaPackageQuotaUSD = &v
		}
		if quotaPackageValidityDays.Valid {
			groupOut.QuotaPackageValidityDays = int(quotaPackageValidityDays.Int64)
		}
		if groupOut.QuotaPackageValidityDays <= 0 {
			groupOut.QuotaPackageValidityDays = 30
		}
		keyOut.Group = groupOut
	}

	return &keyOut, nil
}

func (r *apiKeyRepository) Update(ctx context.Context, key *service.APIKey) error {
	// 使用原子操作：将软删除检查与更新合并到同一语句，避免竞态条件。
	// 之前的实现先检查 Exist 再 UpdateOneID，若在两步之间发生软删除，
	// 则会更新已删除的记录。
	// 这里选择 Update().Where()，确保只有未软删除记录能被更新。
	// 同时显式设置 updated_at，避免二次查询带来的并发可见性问题。
	now := time.Now()
	builder := r.client.APIKey.Update().
		Where(apikey.IDEQ(key.ID), apikey.DeletedAtIsNil()).
		SetName(key.Name).
		SetStatus(key.Status).
		SetUpdatedAt(now)
	if key.GroupID != nil {
		builder.SetGroupID(*key.GroupID)
	} else {
		builder.ClearGroupID()
	}

	// IP 限制字段
	if len(key.IPWhitelist) > 0 {
		builder.SetIPWhitelist(key.IPWhitelist)
	} else {
		builder.ClearIPWhitelist()
	}
	if len(key.IPBlacklist) > 0 {
		builder.SetIPBlacklist(key.IPBlacklist)
	} else {
		builder.ClearIPBlacklist()
	}

	// 用量上限字段
	if key.UsageLimit != nil {
		builder.SetUsageLimit(*key.UsageLimit)
	} else {
		builder.ClearUsageLimit()
	}

	affected, err := builder.Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		// 更新影响行数为 0，说明记录不存在或已被软删除。
		return service.ErrAPIKeyNotFound
	}

	// 使用同一时间戳回填，避免并发删除导致二次查询失败。
	key.UpdatedAt = now
	return nil
}

func (r *apiKeyRepository) Delete(ctx context.Context, id int64) error {
	// 显式软删除：避免依赖 Hook 行为，确保 deleted_at 一定被设置。
	affected, err := r.client.APIKey.Update().
		Where(apikey.IDEQ(id), apikey.DeletedAtIsNil()).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return service.ErrAPIKeyNotFound
		}
		return err
	}
	if affected == 0 {
		exists, err := r.client.APIKey.Query().
			Where(apikey.IDEQ(id)).
			Exist(mixins.SkipSoftDelete(ctx))
		if err != nil {
			return err
		}
		if exists {
			return nil
		}
		return service.ErrAPIKeyNotFound
	}
	return nil
}

func (r *apiKeyRepository) ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams) ([]service.APIKey, *pagination.PaginationResult, error) {
	q := r.activeQuery().Where(apikey.UserIDEQ(userID))

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	keys, err := q.
		WithGroup().
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(apikey.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	outKeys := make([]service.APIKey, 0, len(keys))
	for i := range keys {
		out := apiKeyEntityToService(keys[i])
		_ = r.hydrateAPIKeyGroupQuotaFields(ctx, out)
		outKeys = append(outKeys, *out)
	}

	return outKeys, paginationResultFromTotal(int64(total), params), nil
}

func (r *apiKeyRepository) VerifyOwnership(ctx context.Context, userID int64, apiKeyIDs []int64) ([]int64, error) {
	if len(apiKeyIDs) == 0 {
		return []int64{}, nil
	}

	ids, err := r.client.APIKey.Query().
		Where(apikey.UserIDEQ(userID), apikey.IDIn(apiKeyIDs...), apikey.DeletedAtIsNil()).
		IDs(ctx)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *apiKeyRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	count, err := r.activeQuery().Where(apikey.UserIDEQ(userID)).Count(ctx)
	return int64(count), err
}

func (r *apiKeyRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	count, err := r.activeQuery().Where(apikey.KeyEQ(key)).Count(ctx)
	return count > 0, err
}

func (r *apiKeyRepository) ListByGroupID(ctx context.Context, groupID int64, params pagination.PaginationParams) ([]service.APIKey, *pagination.PaginationResult, error) {
	q := r.activeQuery().Where(apikey.GroupIDEQ(groupID))

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	keys, err := q.
		WithUser().
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(apikey.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	outKeys := make([]service.APIKey, 0, len(keys))
	for i := range keys {
		out := apiKeyEntityToService(keys[i])
		_ = r.hydrateAPIKeyGroupQuotaFields(ctx, out)
		outKeys = append(outKeys, *out)
	}

	return outKeys, paginationResultFromTotal(int64(total), params), nil
}

// SearchAPIKeys searches API keys by user ID and/or keyword (name)
func (r *apiKeyRepository) SearchAPIKeys(ctx context.Context, userID int64, keyword string, limit int) ([]service.APIKey, error) {
	q := r.activeQuery()
	if userID > 0 {
		q = q.Where(apikey.UserIDEQ(userID))
	}

	if keyword != "" {
		q = q.Where(apikey.NameContainsFold(keyword))
	}

	keys, err := q.Limit(limit).Order(dbent.Desc(apikey.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}

	outKeys := make([]service.APIKey, 0, len(keys))
	for i := range keys {
		outKeys = append(outKeys, *apiKeyEntityToService(keys[i]))
	}
	return outKeys, nil
}

// ClearGroupIDByGroupID 将指定分组的所有 API Key 的 group_id 设为 nil
func (r *apiKeyRepository) ClearGroupIDByGroupID(ctx context.Context, groupID int64) (int64, error) {
	n, err := r.client.APIKey.Update().
		Where(apikey.GroupIDEQ(groupID), apikey.DeletedAtIsNil()).
		ClearGroupID().
		Save(ctx)
	return int64(n), err
}

// CountByGroupID 获取分组的 API Key 数量
func (r *apiKeyRepository) CountByGroupID(ctx context.Context, groupID int64) (int64, error) {
	count, err := r.activeQuery().Where(apikey.GroupIDEQ(groupID)).Count(ctx)
	return int64(count), err
}

func (r *apiKeyRepository) ListKeysByUserID(ctx context.Context, userID int64) ([]string, error) {
	keys, err := r.activeQuery().
		Where(apikey.UserIDEQ(userID)).
		Select(apikey.FieldKey).
		Strings(ctx)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *apiKeyRepository) ListKeysByGroupID(ctx context.Context, groupID int64) ([]string, error) {
	keys, err := r.activeQuery().
		Where(apikey.GroupIDEQ(groupID)).
		Select(apikey.FieldKey).
		Strings(ctx)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func apiKeyEntityToService(m *dbent.APIKey) *service.APIKey {
	if m == nil {
		return nil
	}
	out := &service.APIKey{
		ID:           m.ID,
		UserID:       m.UserID,
		Key:          m.Key,
		Name:         m.Name,
		Status:       m.Status,
		IPWhitelist:  m.IPWhitelist,
		IPBlacklist:  m.IPBlacklist,
		UsageLimit:   m.UsageLimit,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		GroupID:      m.GroupID,
		OrgID:        m.OrgID,
		OrgProjectID: m.OrgProjectID,
	}
	if m.Edges.User != nil {
		out.User = userEntityToService(m.Edges.User)
	}
	if m.Edges.Group != nil {
		out.Group = groupEntityToService(m.Edges.Group)
	}
	if m.Edges.Organization != nil {
		out.Organization = organizationEntityToService(m.Edges.Organization)
	}
	return out
}

func userEntityToService(u *dbent.User) *service.User {
	if u == nil {
		return nil
	}
	return &service.User{
		ID:                      u.ID,
		Email:                   u.Email,
		Username:                u.Username,
		Notes:                   u.Notes,
		PasswordHash:            u.PasswordHash,
		Role:                    u.Role,
		Balance:                 u.Balance,
		Concurrency:             u.Concurrency,
		Status:                  u.Status,
		TotpSecretEncrypted:     u.TotpSecretEncrypted,
		TotpEnabled:             u.TotpEnabled,
		TotpEnabledAt:           u.TotpEnabledAt,
		InviteCode:              u.InviteCode,
		DiscoverySource:         u.DiscoverySource,
		InitialBalance:          u.InitialBalance,
		InitialBalanceExpiresAt: u.InitialBalanceExpiresAt,
		IsAgent:                 u.IsAgent,
		AgentStatus:             u.AgentStatus,
		AgentCommissionRate:     u.AgentCommissionRate,
		AgentNote:               u.AgentNote,
		AgentApprovedAt:         u.AgentApprovedAt,
		CreatedAt:               u.CreatedAt,
		UpdatedAt:               u.UpdatedAt,
	}
}

func groupEntityToService(g *dbent.Group) *service.Group {
	if g == nil {
		return nil
	}
	return &service.Group{
		ID:                    g.ID,
		Name:                  g.Name,
		Description:           derefString(g.Description),
		Platform:              g.Platform,
		RateMultiplier:        g.RateMultiplier,
		DisplayRateMultiplier: nil,
		IsExclusive:           g.IsExclusive,
		Status:                g.Status,
		Hydrated:              true,
		SubscriptionType:      g.SubscriptionType,
		DailyLimitUSD:         g.DailyLimitUsd,
		WeeklyLimitUSD:        g.WeeklyLimitUsd,
		MonthlyLimitUSD:       g.MonthlyLimitUsd,
		ImagePrice1K:          g.ImagePrice1k,
		ImagePrice2K:          g.ImagePrice2k,
		ImagePrice4K:          g.ImagePrice4k,
		DefaultValidityDays:   g.DefaultValidityDays,
		ClaudeCodeOnly:        g.ClaudeCodeOnly,
		FallbackGroupID:       g.FallbackGroupID,
		ModelRouting:          g.ModelRouting,
		ModelRoutingEnabled:   g.ModelRoutingEnabled,
		PriceFen:              g.PriceFen,
		Listed:                g.Listed,
		PlanFeatures:          g.PlanFeatures,
		Tags:                  g.Tags,
		ModelPlazaVisible:     g.ModelPlazaVisible,
		DisplayPrice:          g.DisplayPrice,
		DisplayDiscount:       g.DisplayDiscount,
		CreatedAt:             g.CreatedAt,
		UpdatedAt:             g.UpdatedAt,
	}
}

func (r *apiKeyRepository) hydrateAPIKeyGroupQuotaFields(ctx context.Context, key *service.APIKey) error {
	if r.sql == nil || key == nil || key.Group == nil || key.Group.ID <= 0 {
		return nil
	}
	var enabled bool
	var quota sql.NullFloat64
	var validityDays int
	if err := scanSingleRow(ctx, r.sql, `
		SELECT COALESCE(quota_package_enabled, FALSE),
		       quota_package_quota_usd,
		       COALESCE(NULLIF(quota_package_validity_days, 0), 30)
		FROM groups
		WHERE id = $1
	`, []any{key.Group.ID}, &enabled, &quota, &validityDays); err != nil {
		return err
	}
	key.Group.QuotaPackageEnabled = enabled
	if quota.Valid {
		value := quota.Float64
		key.Group.QuotaPackageQuotaUSD = &value
	} else {
		key.Group.QuotaPackageQuotaUSD = nil
	}
	if validityDays <= 0 {
		validityDays = 30
	}
	key.Group.QuotaPackageValidityDays = validityDays
	return nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
