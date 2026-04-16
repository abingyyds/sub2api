package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type subSiteRepository struct {
	db *sql.DB
}

func NewSubSiteRepository(sqlDB *sql.DB) service.SubSiteRepository {
	return &subSiteRepository{db: sqlDB}
}

const subSiteBaseSelect = `
	SELECT
		s.id,
		s.owner_user_id,
		COALESCE(u.email, ''),
		s.parent_sub_site_id,
		COALESCE(parent.name, ''),
		COALESCE(s.level, 1),
		s.name,
		s.slug,
		COALESCE(s.custom_domain, ''),
		s.status,
		COALESCE(s.site_logo, ''),
		COALESCE(s.site_favicon, ''),
		COALESCE(s.site_subtitle, ''),
		COALESCE(s.announcement, ''),
		COALESCE(s.contact_info, ''),
		COALESCE(s.doc_url, ''),
		COALESCE(s.home_content, ''),
		COALESCE(s.theme_template, 'starter'),
		COALESCE(s.theme_config, ''),
		COALESCE(s.custom_config, ''),
		COALESCE(s.registration_mode, 'open'),
		COALESCE(s.enable_topup, TRUE),
		COALESCE(s.allow_sub_site, FALSE),
		COALESCE(s.sub_site_price_fen, 0),
		COALESCE(s.consume_rate_multiplier, 1.0),
		s.subscription_expired_at,
		s.created_at,
		s.updated_at,
		COALESCE((SELECT COUNT(*) FROM sub_site_users su WHERE su.sub_site_id = s.id), 0) AS user_count,
		COALESCE((SELECT COUNT(*) FROM sub_sites child WHERE child.parent_sub_site_id = s.id), 0) AS child_site_count
	FROM sub_sites s
	LEFT JOIN users u ON u.id = s.owner_user_id
	LEFT JOIN sub_sites parent ON parent.id = s.parent_sub_site_id
`

func (r *subSiteRepository) List(ctx context.Context, params pagination.PaginationParams, search, status string) ([]service.SubSite, *pagination.PaginationResult, error) {
	var (
		conds []string
		args  []any
	)
	conds = append(conds, "1=1")
	if search = strings.TrimSpace(search); search != "" {
		args = append(args, "%"+search+"%")
		idx := len(args)
		conds = append(conds, fmt.Sprintf(
			"(s.name ILIKE $%d OR s.slug ILIKE $%d OR COALESCE(s.custom_domain, '') ILIKE $%d OR COALESCE(u.email, '') ILIKE $%d)",
			idx, idx, idx, idx,
		))
	}
	if status = strings.TrimSpace(strings.ToLower(status)); status != "" {
		args = append(args, status)
		conds = append(conds, fmt.Sprintf("s.status = $%d", len(args)))
	}
	where := strings.Join(conds, " AND ")

	countQuery := `SELECT COUNT(*) FROM sub_sites s LEFT JOIN users u ON u.id = s.owner_user_id WHERE ` + where
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, nil, err
	}

	selectArgs := append([]any{}, args...)
	selectArgs = append(selectArgs, params.Limit(), params.Offset())
	query := subSiteBaseSelect + fmt.Sprintf(`
		WHERE %s
		ORDER BY s.id DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)
	rows, err := r.db.QueryContext(ctx, query, selectArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	var items []service.SubSite
	for rows.Next() {
		site, err := scanSubSite(rows)
		if err != nil {
			return nil, nil, err
		}
		items = append(items, *site)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return items, paginationResultFromTotal(total, params), nil
}

func (r *subSiteRepository) ListByOwner(ctx context.Context, ownerUserID int64) ([]service.SubSite, error) {
	rows, err := r.db.QueryContext(ctx, subSiteBaseSelect+`
		WHERE s.owner_user_id = $1
		ORDER BY s.created_at DESC, s.id DESC
	`, ownerUserID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var items []service.SubSite
	for rows.Next() {
		site, err := scanSubSite(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *site)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *subSiteRepository) GetByID(ctx context.Context, id int64) (*service.SubSite, error) {
	row := r.db.QueryRowContext(ctx, subSiteBaseSelect+` WHERE s.id = $1`, id)
	return scanSubSite(row)
}

func (r *subSiteRepository) GetByDomain(ctx context.Context, domain string) (*service.SubSite, error) {
	row := r.db.QueryRowContext(ctx, subSiteBaseSelect+` WHERE LOWER(COALESCE(s.custom_domain, '')) = LOWER($1)`, domain)
	return scanSubSite(row)
}

func (r *subSiteRepository) GetBySlug(ctx context.Context, slug string) (*service.SubSite, error) {
	row := r.db.QueryRowContext(ctx, subSiteBaseSelect+` WHERE LOWER(s.slug) = LOWER($1)`, slug)
	return scanSubSite(row)
}

func (r *subSiteRepository) ExistsBySlug(ctx context.Context, slug string, excludeID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sub_sites WHERE LOWER(slug) = LOWER($1))`
	args := []any{slug}
	if excludeID > 0 {
		query = `SELECT EXISTS(SELECT 1 FROM sub_sites WHERE LOWER(slug) = LOWER($1) AND id <> $2)`
		args = append(args, excludeID)
	}
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *subSiteRepository) ExistsByDomain(ctx context.Context, domain string, excludeID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sub_sites WHERE LOWER(COALESCE(custom_domain, '')) = LOWER($1) AND COALESCE(custom_domain, '') <> '')`
	args := []any{domain}
	if excludeID > 0 {
		query = `SELECT EXISTS(SELECT 1 FROM sub_sites WHERE LOWER(COALESCE(custom_domain, '')) = LOWER($1) AND COALESCE(custom_domain, '') <> '' AND id <> $2)`
		args = append(args, excludeID)
	}
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *subSiteRepository) Create(ctx context.Context, site *service.SubSite) error {
	query := `
		INSERT INTO sub_sites (
			owner_user_id, parent_sub_site_id, level, name, slug, custom_domain, status,
			site_logo, site_favicon, site_subtitle, announcement,
			contact_info, doc_url, home_content, theme_template, theme_config, custom_config,
			registration_mode, enable_topup, allow_sub_site, sub_site_price_fen,
			consume_rate_multiplier, subscription_expired_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, NULLIF($6, ''), $7,
			$8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22,
			$23, NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		site.OwnerUserID, site.ParentSubSiteID, site.Level, site.Name, site.Slug, site.CustomDomain, site.Status,
		site.SiteLogo, site.SiteFavicon, site.SiteSubtitle, site.Announcement,
		site.ContactInfo, site.DocURL, site.HomeContent, site.ThemeTemplate, site.ThemeConfig, site.CustomConfig,
		site.RegistrationMode, site.EnableTopup, site.AllowSubSite, site.SubSitePriceFen,
		site.ConsumeRateMultiplier, site.SubscriptionExpiredAt,
	).Scan(&site.ID, &site.CreatedAt, &site.UpdatedAt)
}

func (r *subSiteRepository) Update(ctx context.Context, site *service.SubSite) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE sub_sites
		SET owner_user_id = $2,
			parent_sub_site_id = $3,
			level = $4,
			name = $5,
			slug = $6,
			custom_domain = NULLIF($7, ''),
			status = $8,
			site_logo = $9,
			site_favicon = $10,
			site_subtitle = $11,
			announcement = $12,
			contact_info = $13,
			doc_url = $14,
			home_content = $15,
			theme_template = $16,
			theme_config = $17,
			custom_config = $18,
			registration_mode = $19,
			enable_topup = $20,
			allow_sub_site = $21,
			sub_site_price_fen = $22,
			consume_rate_multiplier = $23,
			subscription_expired_at = $24,
			updated_at = NOW()
		WHERE id = $1
	`,
		site.ID, site.OwnerUserID, site.ParentSubSiteID, site.Level, site.Name, site.Slug, site.CustomDomain, site.Status,
		site.SiteLogo, site.SiteFavicon, site.SiteSubtitle, site.Announcement,
		site.ContactInfo, site.DocURL, site.HomeContent, site.ThemeTemplate, site.ThemeConfig, site.CustomConfig,
		site.RegistrationMode, site.EnableTopup, site.AllowSubSite, site.SubSitePriceFen, site.ConsumeRateMultiplier,
		site.SubscriptionExpiredAt,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return service.ErrSubSiteNotFound
	}
	return nil
}

func (r *subSiteRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM sub_sites WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return service.ErrSubSiteNotFound
	}
	return nil
}

func (r *subSiteRepository) BindUser(ctx context.Context, siteID int64, userID int64, source string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO sub_site_users (sub_site_id, user_id, source, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET sub_site_id = EXCLUDED.sub_site_id, source = EXCLUDED.source, updated_at = NOW()
	`, siteID, userID, source)
	return err
}

func (r *subSiteRepository) ReplaceGroupPriceOverrides(ctx context.Context, siteID int64, items []service.SubSiteGroupPriceOverride) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM sub_site_group_prices WHERE sub_site_id = $1`, siteID); err != nil {
		return err
	}
	for _, item := range items {
		if item.GroupID <= 0 || item.PriceFen <= 0 {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO sub_site_group_prices (sub_site_id, group_id, price_fen, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
		`, siteID, item.GroupID, item.PriceFen); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *subSiteRepository) ListGroupPriceOverrides(ctx context.Context, siteID int64) ([]service.SubSiteGroupPriceOverride, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.group_id, COALESCE(g.name, ''), p.price_fen, COALESCE(g.price_fen, 0)
		FROM sub_site_group_prices p
		LEFT JOIN groups g ON g.id = p.group_id
		WHERE p.sub_site_id = $1
		ORDER BY p.group_id ASC
	`, siteID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var items []service.SubSiteGroupPriceOverride
	for rows.Next() {
		var item service.SubSiteGroupPriceOverride
		if err := rows.Scan(&item.GroupID, &item.GroupName, &item.PriceFen, &item.BasePriceFen); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *subSiteRepository) ReplaceRechargePriceOverrides(ctx context.Context, siteID int64, items []service.SubSiteRechargePriceOverride) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM sub_site_recharge_prices WHERE sub_site_id = $1`, siteID); err != nil {
		return err
	}
	for _, item := range items {
		if strings.TrimSpace(item.PlanKey) == "" || item.PayAmountFen <= 0 {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO sub_site_recharge_prices (sub_site_id, plan_key, pay_amount_fen, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
		`, siteID, strings.TrimSpace(item.PlanKey), item.PayAmountFen); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *subSiteRepository) ListRechargePriceOverrides(ctx context.Context, siteID int64) ([]service.SubSiteRechargePriceOverride, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT plan_key, pay_amount_fen
		FROM sub_site_recharge_prices
		WHERE sub_site_id = $1
		ORDER BY plan_key ASC
	`, siteID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var items []service.SubSiteRechargePriceOverride
	for rows.Next() {
		var item service.SubSiteRechargePriceOverride
		if err := rows.Scan(&item.PlanKey, &item.PayAmountFen); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *subSiteRepository) CreateActivationRequest(ctx context.Context, request *service.SubSiteActivationRequest) error {
	payload, err := json.Marshal(request.Site)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO sub_site_activation_orders (
			payment_order_id, user_id, parent_sub_site_id, level, validity_days, payload_json, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6::jsonb, NOW(), NOW()
		)
		ON CONFLICT (payment_order_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			parent_sub_site_id = EXCLUDED.parent_sub_site_id,
			level = EXCLUDED.level,
			validity_days = EXCLUDED.validity_days,
			payload_json = EXCLUDED.payload_json,
			updated_at = NOW()
	`,
		request.PaymentOrderID, request.UserID, request.ParentSubSiteID, request.Level, request.ValidityDays, string(payload),
	)
	return err
}

func (r *subSiteRepository) GetActivationRequestByOrderID(ctx context.Context, orderID int64) (*service.SubSiteActivationRequest, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT payment_order_id, user_id, parent_sub_site_id, level, validity_days, payload_json,
			sub_site_id, activated_at, created_at, updated_at
		FROM sub_site_activation_orders
		WHERE payment_order_id = $1
	`, orderID)

	var (
		request     service.SubSiteActivationRequest
		parentID    sql.NullInt64
		payloadJSON []byte
		subSiteID   sql.NullInt64
		activatedAt sql.NullTime
	)
	if err := row.Scan(
		&request.PaymentOrderID,
		&request.UserID,
		&parentID,
		&request.Level,
		&request.ValidityDays,
		&payloadJSON,
		&subSiteID,
		&activatedAt,
		&request.CreatedAt,
		&request.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrSubSiteActivationNotFound
		}
		return nil, err
	}
	if parentID.Valid {
		request.ParentSubSiteID = &parentID.Int64
	}
	if len(payloadJSON) > 0 {
		if err := json.Unmarshal(payloadJSON, &request.Site); err != nil {
			return nil, err
		}
	}
	if subSiteID.Valid {
		request.ActivatedSubSiteID = &subSiteID.Int64
	}
	if activatedAt.Valid {
		request.ActivatedAt = &activatedAt.Time
	}
	return &request, nil
}

func (r *subSiteRepository) MarkActivationRequestCompleted(ctx context.Context, orderID int64, subSiteID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE sub_site_activation_orders
		SET sub_site_id = $2, activated_at = NOW(), updated_at = NOW()
		WHERE payment_order_id = $1
	`, orderID, subSiteID)
	return err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanSubSite(row scanner) (*service.SubSite, error) {
	var (
		site                  service.SubSite
		parentID              sql.NullInt64
		subscriptionExpiredAt sql.NullTime
	)
	if err := row.Scan(
		&site.ID,
		&site.OwnerUserID,
		&site.OwnerEmail,
		&parentID,
		&site.ParentSubSiteName,
		&site.Level,
		&site.Name,
		&site.Slug,
		&site.CustomDomain,
		&site.Status,
		&site.SiteLogo,
		&site.SiteFavicon,
		&site.SiteSubtitle,
		&site.Announcement,
		&site.ContactInfo,
		&site.DocURL,
		&site.HomeContent,
		&site.ThemeTemplate,
		&site.ThemeConfig,
		&site.CustomConfig,
		&site.RegistrationMode,
		&site.EnableTopup,
		&site.AllowSubSite,
		&site.SubSitePriceFen,
		&site.ConsumeRateMultiplier,
		&subscriptionExpiredAt,
		&site.CreatedAt,
		&site.UpdatedAt,
		&site.UserCount,
		&site.ChildSiteCount,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrSubSiteNotFound
		}
		return nil, err
	}
	if parentID.Valid {
		site.ParentSubSiteID = &parentID.Int64
	}
	if subscriptionExpiredAt.Valid {
		site.SubscriptionExpiredAt = &subscriptionExpiredAt.Time
	}
	return &site, nil
}
