package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

// subSiteAdminRepository 提供分站站长后台专用的 JOIN 查询。
// 所有查询都以 sub_site_users.sub_site_id 为隔离键，保证 owner 只看得到自己分站的数据。
type subSiteAdminRepository struct {
	db *sql.DB
}

// NewSubSiteAdminRepository 构造分站站长后台 repository。
func NewSubSiteAdminRepository(db *sql.DB) service.SubSiteAdminRepository {
	return &subSiteAdminRepository{db: db}
}

func (r *subSiteAdminRepository) ListUsersBySubSite(ctx context.Context, siteID int64, params pagination.PaginationParams, search string) ([]service.SubSiteAdminUser, *pagination.PaginationResult, error) {
	conds := []string{"su.sub_site_id = $1"}
	args := []any{siteID}
	if search = strings.TrimSpace(search); search != "" {
		args = append(args, "%"+search+"%")
		idx := len(args)
		conds = append(conds, fmt.Sprintf("(u.email ILIKE $%d OR u.username ILIKE $%d)", idx, idx))
	}
	where := strings.Join(conds, " AND ")

	countQuery := `SELECT COUNT(*) FROM sub_site_users su INNER JOIN users u ON u.id = su.user_id WHERE ` + where
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, nil, err
	}

	selectArgs := append([]any{}, args...)
	selectArgs = append(selectArgs, params.Limit(), params.Offset())
	query := fmt.Sprintf(`
		SELECT u.id, COALESCE(u.email, ''), COALESCE(u.username, ''), COALESCE(u.role, 'user'),
			COALESCE(u.status, 'active'), COALESCE(u.balance, 0),
			COALESCE(su.source, ''), u.created_at, su.updated_at
		FROM sub_site_users su
		INNER JOIN users u ON u.id = su.user_id
		WHERE %s
		ORDER BY u.id DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)

	rows, err := r.db.QueryContext(ctx, query, selectArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]service.SubSiteAdminUser, 0)
	for rows.Next() {
		var it service.SubSiteAdminUser
		if err := rows.Scan(&it.ID, &it.Email, &it.Username, &it.Role, &it.Status, &it.Balance, &it.BindSource, &it.CreatedAt, &it.BoundAt); err != nil {
			return nil, nil, err
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return items, paginationResultFromTotal(total, params), nil
}

func (r *subSiteAdminRepository) ListOrdersBySubSite(ctx context.Context, siteID int64, params pagination.PaginationParams, status string) ([]service.SubSiteAdminOrder, *pagination.PaginationResult, error) {
	conds := []string{"COALESCE(po.sub_site_id, su.sub_site_id) = $1"}
	args := []any{siteID}
	if status = strings.TrimSpace(strings.ToLower(status)); status != "" {
		args = append(args, status)
		conds = append(conds, fmt.Sprintf("po.status = $%d", len(args)))
	}
	where := strings.Join(conds, " AND ")

	countQuery := `
		SELECT COUNT(*) FROM payment_orders po
		LEFT JOIN sub_site_users su ON su.user_id = po.user_id
		WHERE ` + where
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, nil, err
	}

	selectArgs := append([]any{}, args...)
	selectArgs = append(selectArgs, params.Limit(), params.Offset())
	query := fmt.Sprintf(`
		SELECT po.id, po.order_no, po.user_id, COALESCE(u.email, ''),
			COALESCE(po.plan_key, ''), COALESCE(po.order_type, 'subscription'),
			po.amount_fen, COALESCE(po.status, 'pending'), COALESCE(po.pay_method, ''),
			po.paid_at, po.created_at
		FROM payment_orders po
		LEFT JOIN sub_site_users su ON su.user_id = po.user_id
		LEFT JOIN users u ON u.id = po.user_id
		WHERE %s
		ORDER BY po.id DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)

	rows, err := r.db.QueryContext(ctx, query, selectArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]service.SubSiteAdminOrder, 0)
	for rows.Next() {
		var (
			it     service.SubSiteAdminOrder
			paidAt sql.NullTime
		)
		if err := rows.Scan(&it.ID, &it.OrderNo, &it.UserID, &it.UserEmail, &it.PlanKey, &it.OrderType, &it.AmountFen, &it.Status, &it.PayMethod, &paidAt, &it.CreatedAt); err != nil {
			return nil, nil, err
		}
		if paidAt.Valid {
			it.PaidAt = &paidAt.Time
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return items, paginationResultFromTotal(total, params), nil
}

func (r *subSiteAdminRepository) ListUsageBySubSite(ctx context.Context, siteID int64, params pagination.PaginationParams, model string) ([]service.SubSiteAdminUsage, *pagination.PaginationResult, error) {
	conds := []string{"su.sub_site_id = $1"}
	args := []any{siteID}
	if model = strings.TrimSpace(model); model != "" {
		args = append(args, model)
		conds = append(conds, fmt.Sprintf("ul.model = $%d", len(args)))
	}
	where := strings.Join(conds, " AND ")

	countQuery := `
		SELECT COUNT(*) FROM usage_logs ul
		INNER JOIN sub_site_users su ON su.user_id = ul.user_id
		WHERE ` + where
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, nil, err
	}

	selectArgs := append([]any{}, args...)
	selectArgs = append(selectArgs, params.Limit(), params.Offset())
	query := fmt.Sprintf(`
		SELECT ul.id, ul.user_id, COALESCE(u.email, ''), COALESCE(ul.model, ''),
			COALESCE(ul.input_tokens, 0), COALESCE(ul.output_tokens, 0),
			COALESCE(ul.cache_creation_tokens, 0), COALESCE(ul.cache_read_tokens, 0),
			COALESCE(ul.total_cost, 0), COALESCE(ul.actual_cost, 0), ul.created_at
		FROM usage_logs ul
		INNER JOIN sub_site_users su ON su.user_id = ul.user_id
		LEFT JOIN users u ON u.id = ul.user_id
		WHERE %s
		ORDER BY ul.id DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)

	rows, err := r.db.QueryContext(ctx, query, selectArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]service.SubSiteAdminUsage, 0)
	for rows.Next() {
		var it service.SubSiteAdminUsage
		if err := rows.Scan(
			&it.ID, &it.UserID, &it.UserEmail, &it.Model,
			&it.InputTokens, &it.OutputTokens, &it.CacheCreationTokens, &it.CacheReadTokens,
			&it.TotalCost, &it.ActualCost, &it.CreatedAt,
		); err != nil {
			return nil, nil, err
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return items, paginationResultFromTotal(total, params), nil
}

func (r *subSiteAdminRepository) GetDashboardStats(ctx context.Context, siteID int64, rangeStart time.Time) (*service.SubSiteAdminDashboardStats, error) {
	stats := &service.SubSiteAdminDashboardStats{RangeStart: rangeStart}

	// 用户计数
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sub_site_users WHERE sub_site_id = $1
	`, siteID).Scan(&stats.UserCount); err != nil {
		return nil, err
	}

	// 活跃用户 —— 在 rangeStart 以后有 usage_logs 记录的用户数
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT ul.user_id)
		FROM usage_logs ul
		INNER JOIN sub_site_users su ON su.user_id = ul.user_id
		WHERE su.sub_site_id = $1 AND ul.created_at >= $2
	`, siteID, rangeStart).Scan(&stats.ActiveUsers); err != nil {
		return nil, err
	}

	// 区间请求数 + 总成本
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(ul.actual_cost), 0)
		FROM usage_logs ul
		INNER JOIN sub_site_users su ON su.user_id = ul.user_id
		WHERE su.sub_site_id = $1 AND ul.created_at >= $2
	`, siteID, rangeStart).Scan(&stats.Requests, &stats.TotalCost); err != nil {
		return nil, err
	}

	// 区间订单收入（已支付）
	if err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(po.amount_fen), 0)
		FROM payment_orders po
		LEFT JOIN sub_site_users su ON su.user_id = po.user_id
		WHERE COALESCE(po.sub_site_id, su.sub_site_id) = $1 AND po.status = 'paid' AND COALESCE(po.paid_at, po.created_at) >= $2
	`, siteID, rangeStart).Scan(&stats.RevenueFen); err != nil {
		return nil, err
	}

	// 分站池余额快照
	if err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(balance_fen, 0), COALESCE(total_topup_fen, 0), COALESCE(total_consumed_fen, 0)
		FROM sub_sites WHERE id = $1
	`, siteID).Scan(&stats.PoolBalanceFen, &stats.TotalTopupFen, &stats.TotalConsumedFen); err != nil {
		return nil, err
	}

	return stats, nil
}
