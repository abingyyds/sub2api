package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type agentRepository struct {
	db *sql.DB
}

func NewAgentRepository(sqlDB *sql.DB) service.AgentRepository {
	return &agentRepository{db: sqlDB}
}

// --- Sub-user queries (via referrals table) ---

func (r *agentRepository) CountSubUsers(ctx context.Context, agentID int64) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM referrals ref JOIN users u ON u.id = ref.invitee_id
		 WHERE ref.inviter_id = $1 AND u.deleted_at IS NULL`, agentID).Scan(&count)
	return count, err
}

func (r *agentRepository) ListSubUsers(ctx context.Context, agentID int64, params pagination.PaginationParams, search string) ([]service.AgentSubUser, *pagination.PaginationResult, error) {
	// Count
	var total int64
	countBase := `SELECT COUNT(*) FROM referrals ref JOIN users u ON u.id = ref.invitee_id
		WHERE ref.inviter_id = $1 AND u.deleted_at IS NULL`
	if search != "" {
		p := "%" + search + "%"
		err := r.db.QueryRowContext(ctx, countBase+` AND u.email ILIKE $2`, agentID, p).Scan(&total)
		if err != nil {
			return nil, nil, err
		}
	} else {
		if err := r.db.QueryRowContext(ctx, countBase, agentID).Scan(&total); err != nil {
			return nil, nil, err
		}
	}

	// Query
	queryBase := `SELECT u.id, u.email, u.username, u.balance, u.status, u.created_at,
		COALESCE((SELECT SUM(po.balance_amount) FROM payment_orders po WHERE po.user_id = u.id AND po.status = 'paid'), 0) AS total_recharge,
		COALESCE((SELECT SUM(ul.total_cost) FROM usage_logs ul WHERE ul.user_id = u.id), 0) AS total_consumed,
		COALESCE((SELECT COUNT(*) FROM payment_orders po2 WHERE po2.user_id = u.id AND po2.status = 'paid'), 0) AS order_count
		FROM referrals ref JOIN users u ON u.id = ref.invitee_id
		WHERE ref.inviter_id = $1 AND u.deleted_at IS NULL`

	var rows *sql.Rows
	var err error
	if search != "" {
		p := "%" + search + "%"
		rows, err = r.db.QueryContext(ctx,
			queryBase+` AND u.email ILIKE $2 ORDER BY u.created_at DESC LIMIT $3 OFFSET $4`,
			agentID, p, params.Limit(), params.Offset())
	} else {
		rows, err = r.db.QueryContext(ctx,
			queryBase+` ORDER BY u.created_at DESC LIMIT $2 OFFSET $3`,
			agentID, params.Limit(), params.Offset())
	}
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []service.AgentSubUser
	for rows.Next() {
		var su service.AgentSubUser
		if err := rows.Scan(&su.ID, &su.Email, &su.Username, &su.Balance, &su.Status, &su.CreatedAt,
			&su.TotalRecharge, &su.TotalConsumed, &su.OrderCount); err != nil {
			return nil, nil, err
		}
		results = append(results, su)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return results, paginationResultFromTotal(total, params), nil
}

// --- Sub-user financial logs ---

func (r *agentRepository) ListSubUserPaymentOrders(ctx context.Context, agentID int64, params pagination.PaginationParams, search string) ([]service.AgentFinancialLog, *pagination.PaginationResult, error) {
	// Count
	var total int64
	countQuery := `SELECT COUNT(*)
		FROM payment_orders po
		JOIN referrals ref ON ref.invitee_id = po.user_id AND ref.inviter_id = $1
		WHERE po.status = 'paid'`
	if search != "" {
		p := "%" + search + "%"
		countQuery += ` AND EXISTS (SELECT 1 FROM users u WHERE u.id = po.user_id AND u.email ILIKE $2)`
		err := r.db.QueryRowContext(ctx, countQuery, agentID, p).Scan(&total)
		if err != nil {
			return nil, nil, err
		}
	} else {
		if err := r.db.QueryRowContext(ctx, countQuery, agentID).Scan(&total); err != nil {
			return nil, nil, err
		}
	}

	// Query
	queryBase := `SELECT po.id, po.user_id, u.email, po.balance_amount, po.plan_key, po.order_type, po.amount_fen, po.created_at
		FROM payment_orders po
		JOIN referrals ref ON ref.invitee_id = po.user_id AND ref.inviter_id = $1
		JOIN users u ON u.id = po.user_id
		WHERE po.status = 'paid'`

	var rows *sql.Rows
	var err error
	if search != "" {
		p := "%" + search + "%"
		rows, err = r.db.QueryContext(ctx,
			queryBase+` AND u.email ILIKE $2 ORDER BY po.created_at DESC LIMIT $3 OFFSET $4`,
			agentID, p, params.Limit(), params.Offset())
	} else {
		rows, err = r.db.QueryContext(ctx,
			queryBase+` ORDER BY po.created_at DESC LIMIT $2 OFFSET $3`,
			agentID, params.Limit(), params.Offset())
	}
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []service.AgentFinancialLog
	for rows.Next() {
		var log service.AgentFinancialLog
		var planKey, orderType string
		var amountFen int
		if err := rows.Scan(&log.ID, &log.UserID, &log.UserEmail, &log.Amount, &planKey, &orderType, &amountFen, &log.CreatedAt); err != nil {
			return nil, nil, err
		}
		log.Type = "payment"
		log.Detail = fmt.Sprintf("%s (%s) ¥%.2f", planKey, orderType, float64(amountFen)/100)
		results = append(results, log)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return results, paginationResultFromTotal(total, params), nil
}

// --- Dashboard stats ---

func (r *agentRepository) GetDashboardStats(ctx context.Context, agentID int64) (*service.AgentDashboardStats, error) {
	stats := &service.AgentDashboardStats{}

	// Total sub-users
	r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM referrals WHERE inviter_id = $1`, agentID).Scan(&stats.TotalSubUsers)

	// Total recharge & consumed
	r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(po.balance_amount), 0)
		 FROM payment_orders po
		 JOIN referrals ref ON ref.invitee_id = po.user_id AND ref.inviter_id = $1
		 WHERE po.status = 'paid'`, agentID).Scan(&stats.TotalRecharge)

	r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(ul.total_cost), 0)
		 FROM usage_logs ul
		 JOIN referrals ref ON ref.invitee_id = ul.user_id AND ref.inviter_id = $1`, agentID).Scan(&stats.TotalConsumed)

	// Commission stats
	r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(commission_amount), 0),
			COALESCE(SUM(commission_amount) FILTER (WHERE status = 'pending'), 0),
			COALESCE(SUM(commission_amount) FILTER (WHERE status = 'settled'), 0)
		 FROM agent_commissions WHERE agent_id = $1`, agentID).Scan(
		&stats.TotalCommission, &stats.PendingCommission, &stats.SettledCommission)

	// Today stats
	today := time.Now().Truncate(24 * time.Hour)
	r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM referrals ref
		 JOIN users u ON u.id = ref.invitee_id
		 WHERE ref.inviter_id = $1 AND u.created_at >= $2`, agentID, today).Scan(&stats.TodayNewUsers)

	r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(po.balance_amount), 0)
		 FROM payment_orders po
		 JOIN referrals ref ON ref.invitee_id = po.user_id AND ref.inviter_id = $1
		 WHERE po.status = 'paid' AND po.paid_at >= $2`, agentID, today).Scan(&stats.TodayRecharge)

	return stats, nil
}

// --- Commission CRUD ---

func (r *agentRepository) CreateCommission(ctx context.Context, c *service.AgentCommission) error {
	query := `INSERT INTO agent_commissions (agent_id, user_id, order_id, source_type, source_amount, commission_rate, commission_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		c.AgentID, c.UserID, c.OrderID, c.SourceType, c.SourceAmount,
		c.CommissionRate, c.CommissionAmount, c.Status,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *agentRepository) ListCommissions(ctx context.Context, agentID int64, params pagination.PaginationParams, status string) ([]service.AgentCommission, *pagination.PaginationResult, error) {
	// Count
	var total int64
	countQuery := `SELECT COUNT(*) FROM agent_commissions WHERE agent_id = $1`
	args := []any{agentID}
	if status != "" {
		countQuery += ` AND status = $2`
		args = append(args, status)
	}
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, nil, err
	}

	// Query
	selectBase := `SELECT ac.id, ac.agent_id, ac.user_id, ac.order_id, ac.source_type, ac.source_amount,
		ac.commission_rate, ac.commission_amount, ac.status, ac.settled_at, ac.created_at, ac.updated_at,
		u.email AS user_email, COALESCE(po.order_no, '') AS order_no
		FROM agent_commissions ac
		JOIN users u ON u.id = ac.user_id
		LEFT JOIN payment_orders po ON po.id = ac.order_id
		WHERE ac.agent_id = $1`

	queryArgs := []any{agentID}
	if status != "" {
		selectBase += ` AND ac.status = $2`
		queryArgs = append(queryArgs, status)
	}
	selectBase += fmt.Sprintf(` ORDER BY ac.created_at DESC LIMIT $%d OFFSET $%d`, len(queryArgs)+1, len(queryArgs)+2)
	queryArgs = append(queryArgs, params.Limit(), params.Offset())

	rows, err := r.db.QueryContext(ctx, selectBase, queryArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []service.AgentCommission
	for rows.Next() {
		var c service.AgentCommission
		var settledAt sql.NullTime
		var orderID sql.NullInt64
		if err := rows.Scan(&c.ID, &c.AgentID, &c.UserID, &orderID, &c.SourceType, &c.SourceAmount,
			&c.CommissionRate, &c.CommissionAmount, &c.Status, &settledAt, &c.CreatedAt, &c.UpdatedAt,
			&c.UserEmail, &c.OrderNo); err != nil {
			return nil, nil, err
		}
		if settledAt.Valid {
			c.SettledAt = &settledAt.Time
		}
		if orderID.Valid {
			c.OrderID = &orderID.Int64
		}
		results = append(results, c)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return results, paginationResultFromTotal(total, params), nil
}

func (r *agentRepository) SettlePendingCommissions(ctx context.Context, agentID int64) (float64, error) {
	var totalAmount float64
	err := r.db.QueryRowContext(ctx,
		`WITH settled AS (
			UPDATE agent_commissions SET status = 'settled', settled_at = NOW(), updated_at = NOW()
			WHERE agent_id = $1 AND status = 'pending'
			RETURNING commission_amount
		) SELECT COALESCE(SUM(commission_amount), 0) FROM settled`, agentID).Scan(&totalAmount)
	return totalAmount, err
}

// --- Admin queries ---

func (r *agentRepository) ListAgents(ctx context.Context, params pagination.PaginationParams, status string, search string) ([]service.AgentInfo, *pagination.PaginationResult, error) {
	// Build count query
	countQuery := `SELECT COUNT(*) FROM users WHERE (is_agent = true OR agent_status != '') AND deleted_at IS NULL`
	var countArgs []any
	argIdx := 1

	if status != "" {
		countQuery += fmt.Sprintf(` AND agent_status = $%d`, argIdx)
		countArgs = append(countArgs, status)
		argIdx++
	}
	if search != "" {
		countQuery += fmt.Sprintf(` AND email ILIKE $%d`, argIdx)
		countArgs = append(countArgs, "%"+search+"%")
		argIdx++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, nil, err
	}

	// Build select query
	selectQuery := `SELECT u.id, u.email, u.username, u.is_agent, u.agent_status, u.agent_commission_rate,
		u.agent_note, u.agent_approved_at, COALESCE(u.invite_code, ''), u.created_at,
		(SELECT COUNT(*) FROM referrals WHERE inviter_id = u.id) AS sub_user_count,
		COALESCE((SELECT SUM(commission_amount) FROM agent_commissions WHERE agent_id = u.id), 0) AS total_commission
		FROM users u WHERE (u.is_agent = true OR u.agent_status != '') AND u.deleted_at IS NULL`

	var selectArgs []any
	argIdx = 1

	if status != "" {
		selectQuery += fmt.Sprintf(` AND u.agent_status = $%d`, argIdx)
		selectArgs = append(selectArgs, status)
		argIdx++
	}
	if search != "" {
		selectQuery += fmt.Sprintf(` AND u.email ILIKE $%d`, argIdx)
		selectArgs = append(selectArgs, "%"+search+"%")
		argIdx++
	}

	selectQuery += fmt.Sprintf(` ORDER BY u.id DESC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	selectArgs = append(selectArgs, params.Limit(), params.Offset())

	rows, err := r.db.QueryContext(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []service.AgentInfo
	for rows.Next() {
		var a service.AgentInfo
		var approvedAt sql.NullTime
		if err := rows.Scan(&a.ID, &a.Email, &a.Username, &a.IsAgent, &a.AgentStatus, &a.CommissionRate,
			&a.AgentNote, &approvedAt, &a.InviteCode, &a.CreatedAt,
			&a.SubUserCount, &a.TotalCommission); err != nil {
			return nil, nil, err
		}
		if approvedAt.Valid {
			a.ApprovedAt = &approvedAt.Time
		}
		results = append(results, a)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return results, paginationResultFromTotal(total, params), nil
}

// GetAgentByUserID finds the agent (inviter) for a given user via referrals table.
func (r *agentRepository) GetAgentByUserID(ctx context.Context, userID int64) (int64, error) {
	var agentID int64
	err := r.db.QueryRowContext(ctx,
		`SELECT ref.inviter_id FROM referrals ref
		 JOIN users u ON u.id = ref.inviter_id
		 WHERE ref.invitee_id = $1 AND u.is_agent = true AND u.agent_status = 'approved'`,
		userID).Scan(&agentID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return agentID, err
}
