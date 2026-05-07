package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
		COALESCE(s.mode, 'pool'),
		COALESCE(s.site_logo, ''),
		COALESCE(s.site_favicon, ''),
		COALESCE(s.site_subtitle, ''),
		COALESCE(s.announcement, ''),
		COALESCE(s.contact_info, ''),
		COALESCE(s.doc_url, ''),
		COALESCE(s.home_content, ''),
		COALESCE(s.pending_home_content, ''),
		COALESCE(s.home_content_review_status, 'none'),
		COALESCE(s.home_content_review_note, ''),
		s.home_content_submitted_at,
		s.home_content_reviewed_at,
		s.home_content_reviewed_by,
		COALESCE(s.theme_template, 'starter'),
		COALESCE(s.registration_mode, 'open'),
		COALESCE(s.enable_topup, TRUE),
		COALESCE(s.allow_sub_site, FALSE),
		COALESCE(s.sub_site_price_fen, 0),
		COALESCE(s.consume_rate_multiplier, 1.0),
		COALESCE(s.balance_fen, 0),
		COALESCE(s.total_topup_fen, 0),
		COALESCE(s.total_consumed_fen, 0),
		COALESCE(s.total_withdrawn_fen, 0),
		COALESCE(s.allow_online_topup, TRUE),
		COALESCE(s.allow_offline_topup, TRUE),
		s.owner_payment_config,
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
	ownerPaymentJSON, err := marshalOwnerPaymentConfig(site.OwnerPaymentConfig)
	if err != nil {
		return err
	}
	mode := normalizeSubSiteMode(site.Mode)
	site.Mode = mode
	query := `
		INSERT INTO sub_sites (
			owner_user_id, parent_sub_site_id, level, name, slug, custom_domain, status, mode,
			site_logo, site_favicon, site_subtitle, announcement,
			contact_info, doc_url, home_content, theme_template,
			registration_mode, enable_topup, allow_sub_site, sub_site_price_fen,
			consume_rate_multiplier, allow_online_topup, allow_offline_topup,
			owner_payment_config, subscription_expired_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, NULLIF($6, ''), $7, $8,
			$9, $10, $11, $12,
			$13, $14, $15, $16,
			$17, $18, $19, $20,
			$21, $22, $23,
			$24::jsonb, $25, NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		site.OwnerUserID, site.ParentSubSiteID, site.Level, site.Name, site.Slug, site.CustomDomain, site.Status, mode,
		site.SiteLogo, site.SiteFavicon, site.SiteSubtitle, site.Announcement,
		site.ContactInfo, site.DocURL, site.HomeContent, site.ThemeTemplate,
		site.RegistrationMode, site.EnableTopup, site.AllowSubSite, site.SubSitePriceFen,
		site.ConsumeRateMultiplier, site.AllowOnlineTopup, site.AllowOfflineTopup,
		ownerPaymentJSON, site.SubscriptionExpiredAt,
	).Scan(&site.ID, &site.CreatedAt, &site.UpdatedAt)
}

func (r *subSiteRepository) Update(ctx context.Context, site *service.SubSite) error {
	ownerPaymentJSON, err := marshalOwnerPaymentConfig(site.OwnerPaymentConfig)
	if err != nil {
		return err
	}
	mode := normalizeSubSiteMode(site.Mode)
	site.Mode = mode
	res, err := r.db.ExecContext(ctx, `
		UPDATE sub_sites
		SET owner_user_id = $2,
			parent_sub_site_id = $3,
			level = $4,
			name = $5,
			slug = $6,
			custom_domain = NULLIF($7, ''),
			status = $8,
			mode = $9,
			site_logo = $10,
			site_favicon = $11,
			site_subtitle = $12,
			announcement = $13,
			contact_info = $14,
			doc_url = $15,
			home_content = $16,
			theme_template = $17,
			registration_mode = $18,
			enable_topup = $19,
			allow_sub_site = $20,
			sub_site_price_fen = $21,
			consume_rate_multiplier = $22,
			allow_online_topup = $23,
			allow_offline_topup = $24,
			owner_payment_config = $25::jsonb,
			subscription_expired_at = $26,
			updated_at = NOW()
		WHERE id = $1
	`,
		site.ID, site.OwnerUserID, site.ParentSubSiteID, site.Level, site.Name, site.Slug, site.CustomDomain, site.Status, mode,
		site.SiteLogo, site.SiteFavicon, site.SiteSubtitle, site.Announcement,
		site.ContactInfo, site.DocURL, site.HomeContent, site.ThemeTemplate,
		site.RegistrationMode, site.EnableTopup, site.AllowSubSite, site.SubSitePriceFen, site.ConsumeRateMultiplier,
		site.AllowOnlineTopup, site.AllowOfflineTopup,
		ownerPaymentJSON, site.SubscriptionExpiredAt,
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

// UpdateMode 仅更新 sub_sites.mode（由 admin 授权后使用），不触动其他字段。
func (r *subSiteRepository) UpdateMode(ctx context.Context, siteID int64, newMode string) error {
	mode := normalizeSubSiteMode(newMode)
	res, err := r.db.ExecContext(ctx, `
		UPDATE sub_sites SET mode = $2, updated_at = NOW() WHERE id = $1
	`, siteID, mode)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return service.ErrSubSiteNotFound
	}
	return nil
}

func (r *subSiteRepository) SubmitHomeContentReview(ctx context.Context, siteID int64, pendingContent string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE sub_sites
		SET pending_home_content = $2,
			home_content_review_status = 'pending',
			home_content_review_note = '',
			home_content_submitted_at = NOW(),
			home_content_reviewed_at = NULL,
			home_content_reviewed_by = NULL,
			updated_at = NOW()
		WHERE id = $1
	`, siteID, pendingContent)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return service.ErrSubSiteNotFound
	}
	return nil
}

func (r *subSiteRepository) ReviewHomeContent(ctx context.Context, siteID int64, approved bool, reviewerID int64, reviewNote string) error {
	status := service.SubSiteHomeContentReviewRejected
	setHomeContent := `home_content`
	if approved {
		status = service.SubSiteHomeContentReviewApproved
		setHomeContent = `COALESCE(pending_home_content, '')`
	}
	res, err := r.db.ExecContext(ctx, fmt.Sprintf(`
		UPDATE sub_sites
		SET home_content = %s,
			pending_home_content = NULL,
			home_content_review_status = $2,
			home_content_review_note = $3,
			home_content_reviewed_at = NOW(),
			home_content_reviewed_by = NULLIF($4, 0),
			updated_at = NOW()
		WHERE id = $1 AND home_content_review_status = 'pending'
	`, setHomeContent), siteID, status, reviewNote, reviewerID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return service.ErrSubSiteNotFound
	}
	return nil
}

// IncrementTotalWithdrawnFen 提现打款完成时累加分站已提现总额。
func (r *subSiteRepository) IncrementTotalWithdrawnFen(ctx context.Context, siteID int64, amountFen int64) error {
	if amountFen <= 0 {
		return nil
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE sub_sites SET total_withdrawn_fen = total_withdrawn_fen + $2, updated_at = NOW() WHERE id = $1
	`, siteID, amountFen)
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

// CascadeUpdateStatus 使用递归 CTE 找到 rootID 的所有后代，更新它们的 status。
// 不包含 rootID 本身（调用方通过 Update 单独处理 root）。返回被更新的后代 id 列表。
func (r *subSiteRepository) CascadeUpdateStatus(ctx context.Context, rootID int64, newStatus string) ([]int64, error) {
	if rootID <= 0 {
		return nil, nil
	}
	rows, err := r.db.QueryContext(ctx, `
		WITH RECURSIVE descendants AS (
			SELECT id FROM sub_sites WHERE parent_sub_site_id = $1
			UNION ALL
			SELECT s.id FROM sub_sites s
			INNER JOIN descendants d ON s.parent_sub_site_id = d.id
		)
		UPDATE sub_sites SET status = $2, updated_at = NOW()
		WHERE id IN (SELECT id FROM descendants) AND status <> $2
		RETURNING id
	`, rootID, newStatus)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
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

func (r *subSiteRepository) UnbindUser(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sub_site_users WHERE user_id = $1`, userID)
	return err
}

func (r *subSiteRepository) GetBoundSubSiteByUserID(ctx context.Context, userID int64) (*service.SubSite, error) {
	row := r.db.QueryRowContext(ctx, subSiteBaseSelect+`
		INNER JOIN sub_site_users su ON su.sub_site_id = s.id
		WHERE su.user_id = $1
		LIMIT 1
	`, userID)
	site, err := scanSubSite(row)
	if err != nil {
		if errors.Is(err, service.ErrSubSiteNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return site, nil
}

// AdjustBalance 原子调整分站余额池并写入流水。deltaFen 为正表示入账，负表示出账；
// 写入逻辑：事务内先 UPDATE sub_sites 返回新余额，再 INSERT sub_site_ledger。
// total_topup_fen / total_consumed_fen 根据 tx_type 分类累加，便于面板展示。
// 允许扣成负余额（透支），上层业务可根据需要记 warn 或拦截。
func (r *subSiteRepository) AdjustBalance(ctx context.Context, siteID int64, deltaFen int64, entry service.SubSiteLedgerEntry) (int64, error) {
	if siteID <= 0 {
		return 0, fmt.Errorf("invalid sub_site_id")
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	newBalance, err := r.adjustBalanceInTx(ctx, tx, siteID, deltaFen, entry, false)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return newBalance, nil
}

// ApplyUserBalanceAndPoolLedger 在单个 SQL 事务内完成用户余额变动和资金池流水。
// 与 AdjustBalance 不同，这条路径禁止任何分站池扣成负数，避免“用户已加余额但池扣款失败”的账实不一致。
func (r *subSiteRepository) ApplyUserBalanceAndPoolLedger(ctx context.Context, userID int64, balanceDelta float64, entries []service.SubSiteLedgerEntry) error {
	if userID <= 0 {
		return service.ErrUserNotFound
	}
	if len(entries) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `UPDATE users SET balance = balance + $2, updated_at = NOW() WHERE id = $1`, userID, balanceDelta)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return service.ErrUserNotFound
	}
	for _, entry := range entries {
		if entry.SubSiteID <= 0 {
			return service.ErrSubSiteNotFound
		}
		if _, err := r.adjustBalanceInTx(ctx, tx, entry.SubSiteID, entry.DeltaFen, entry, true); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *subSiteRepository) adjustBalanceInTx(ctx context.Context, tx *sql.Tx, siteID int64, deltaFen int64, entry service.SubSiteLedgerEntry, preventNegative bool) (int64, error) {
	// 累计入账：topup 常规 + 分站利润
	isTopup := deltaFen > 0 && (entry.TxType == service.SubSiteLedgerTopupOnline ||
		entry.TxType == service.SubSiteLedgerTopupAdmin ||
		entry.TxType == service.SubSiteLedgerManualCredit ||
		entry.TxType == service.SubSiteLedgerProfit)
	// 累计出账：pool 模式消费扣池
	isConsume := deltaFen < 0 && entry.TxType == service.SubSiteLedgerConsume

	var newBalance int64
	updateQuery := `
		UPDATE sub_sites
		SET balance_fen = balance_fen + $2,
			total_topup_fen = total_topup_fen + CASE WHEN $3 THEN $2 ELSE 0 END,
			total_consumed_fen = total_consumed_fen + CASE WHEN $4 THEN -$2 ELSE 0 END,
			updated_at = NOW()
		WHERE id = $1 AND (NOT $5 OR balance_fen + $2 >= 0)
		RETURNING balance_fen
	`
	if err := tx.QueryRowContext(ctx, updateQuery, siteID, deltaFen, isTopup, isConsume, preventNegative).Scan(&newBalance); err != nil {
		if err == sql.ErrNoRows {
			var exists bool
			if existsErr := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sub_sites WHERE id = $1)`, siteID).Scan(&exists); existsErr != nil {
				return 0, existsErr
			}
			if !exists {
				return 0, service.ErrSubSiteNotFound
			}
			return 0, service.ErrSubSitePoolInsufficient
		}
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO sub_site_ledger (
			sub_site_id, tx_type, delta_fen, balance_after_fen,
			related_user_id, related_usage_log_id, related_order_id, operator_id, note, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
	`,
		siteID, entry.TxType, deltaFen, newBalance,
		entry.RelatedUserID, entry.RelatedUsageLogID, entry.RelatedOrderID, entry.OperatorID, entry.Note,
	); err != nil {
		return 0, err
	}
	return newBalance, nil
}

func (r *subSiteRepository) ListLedger(ctx context.Context, siteID int64, params pagination.PaginationParams, txType string) ([]service.SubSiteLedgerEntry, *pagination.PaginationResult, error) {
	var (
		conds = []string{"sub_site_id = $1"}
		args  = []any{siteID}
	)
	if txType = strings.TrimSpace(txType); txType != "" {
		args = append(args, txType)
		conds = append(conds, fmt.Sprintf("tx_type = $%d", len(args)))
	}
	where := strings.Join(conds, " AND ")

	countQuery := `SELECT COUNT(*) FROM sub_site_ledger WHERE ` + where
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, nil, err
	}

	selectArgs := append([]any{}, args...)
	selectArgs = append(selectArgs, params.Limit(), params.Offset())
	query := fmt.Sprintf(`
		SELECT id, sub_site_id, tx_type, delta_fen, balance_after_fen,
			related_user_id, related_usage_log_id, related_order_id, operator_id, COALESCE(note, ''), created_at
		FROM sub_site_ledger
		WHERE %s
		ORDER BY id DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)
	rows, err := r.db.QueryContext(ctx, query, selectArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	var items []service.SubSiteLedgerEntry
	for rows.Next() {
		var (
			entry         service.SubSiteLedgerEntry
			relUserID     sql.NullInt64
			relUsageLogID sql.NullInt64
			relOrderID    sql.NullInt64
			operatorID    sql.NullInt64
		)
		if err := rows.Scan(
			&entry.ID, &entry.SubSiteID, &entry.TxType, &entry.DeltaFen, &entry.BalanceAfterFen,
			&relUserID, &relUsageLogID, &relOrderID, &operatorID, &entry.Note, &entry.CreatedAt,
		); err != nil {
			return nil, nil, err
		}
		if relUserID.Valid {
			entry.RelatedUserID = &relUserID.Int64
		}
		if relUsageLogID.Valid {
			entry.RelatedUsageLogID = &relUsageLogID.Int64
		}
		if relOrderID.Valid {
			entry.RelatedOrderID = &relOrderID.Int64
		}
		if operatorID.Valid {
			entry.OperatorID = &operatorID.Int64
		}
		items = append(items, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return items, paginationResultFromTotal(total, params), nil
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
		site                   service.SubSite
		parentID               sql.NullInt64
		ownerPaymentRaw        []byte
		subscriptionExpiredAt  sql.NullTime
		homeContentSubmittedAt sql.NullTime
		homeContentReviewedAt  sql.NullTime
		homeContentReviewedBy  sql.NullInt64
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
		&site.Mode,
		&site.SiteLogo,
		&site.SiteFavicon,
		&site.SiteSubtitle,
		&site.Announcement,
		&site.ContactInfo,
		&site.DocURL,
		&site.HomeContent,
		&site.PendingHomeContent,
		&site.HomeContentReviewStatus,
		&site.HomeContentReviewNote,
		&homeContentSubmittedAt,
		&homeContentReviewedAt,
		&homeContentReviewedBy,
		&site.ThemeTemplate,
		&site.RegistrationMode,
		&site.EnableTopup,
		&site.AllowSubSite,
		&site.SubSitePriceFen,
		&site.ConsumeRateMultiplier,
		&site.BalanceFen,
		&site.TotalTopupFen,
		&site.TotalConsumedFen,
		&site.TotalWithdrawnFen,
		&site.AllowOnlineTopup,
		&site.AllowOfflineTopup,
		&ownerPaymentRaw,
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
	if homeContentSubmittedAt.Valid {
		site.HomeContentSubmittedAt = &homeContentSubmittedAt.Time
	}
	if homeContentReviewedAt.Valid {
		site.HomeContentReviewedAt = &homeContentReviewedAt.Time
	}
	if homeContentReviewedBy.Valid {
		site.HomeContentReviewedBy = &homeContentReviewedBy.Int64
	}
	site.Mode = normalizeSubSiteMode(site.Mode)
	if cfg, err := unmarshalOwnerPaymentConfig(ownerPaymentRaw); err != nil {
		return nil, err
	} else {
		site.OwnerPaymentConfig = cfg
	}
	return &site, nil
}

// normalizeSubSiteMode 把任意输入对齐到允许的枚举值；未知/空字符串回退 pool。
func normalizeSubSiteMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case service.SubSiteModeRate:
		return service.SubSiteModeRate
	case service.SubSiteModePool, "":
		return service.SubSiteModePool
	default:
		return service.SubSiteModePool
	}
}

func marshalOwnerPaymentConfig(cfg *service.OwnerPaymentConfig) (any, error) {
	if cfg == nil {
		return nil, nil
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func unmarshalOwnerPaymentConfig(raw []byte) (*service.OwnerPaymentConfig, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var cfg service.OwnerPaymentConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
