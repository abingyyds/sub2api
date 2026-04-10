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
		COALESCE((SELECT COUNT(*) FROM payment_orders po2 WHERE po2.user_id = u.id AND po2.status = 'paid'), 0) AS order_count,
		ref.commission_rate,
		(u.is_agent = true AND u.agent_status = 'approved') AS is_sub_agent
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
		var commRate sql.NullFloat64
		if err := rows.Scan(&su.ID, &su.Email, &su.Username, &su.Balance, &su.Status, &su.CreatedAt,
			&su.TotalRecharge, &su.TotalConsumed, &su.OrderCount, &commRate, &su.IsAgent); err != nil {
			return nil, nil, err
		}
		if commRate.Valid {
			v := commRate.Float64
			su.CommissionRate = &v
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

func (r *agentRepository) GetDashboardStats(ctx context.Context, agentID int64, siteBalance float64) (*service.AgentDashboardStats, error) {
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

	stats.SiteBalance = siteBalance
	r.db.QueryRowContext(ctx,
		`SELECT COALESCE(frozen_balance, 0), COALESCE(withdrawable_balance, 0), COALESCE(total_withdrawn, 0)
		 FROM agent_profiles WHERE user_id = $1`, agentID).Scan(&stats.FrozenBalance, &stats.WithdrawableBalance, &stats.TotalWithdrawn)

	return stats, nil
}

// --- Commission CRUD ---

func (r *agentRepository) CreateCommission(ctx context.Context, c *service.AgentCommission) error {
	query := `INSERT INTO agent_commissions (agent_id, user_id, order_id, source_type, source_amount, commission_rate, commission_amount, status, settled_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		c.AgentID, c.UserID, c.OrderID, c.SourceType, c.SourceAmount,
		c.CommissionRate, c.CommissionAmount, c.Status, c.SettledAt,
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
		COALESCE(ap.real_name, ''), COALESCE(ap.phone, ''), COALESCE(ap.identity_status, 'unsubmitted'),
		COALESCE(ap.contract_status, 'unsigned'), ap.activation_fee_paid_at, COALESCE(ap.is_frozen, false),
		COALESCE(ap.frozen_reason, ''), COALESCE(ap.frozen_balance, 0), COALESCE(ap.withdrawable_balance, 0), COALESCE(ap.total_withdrawn, 0),
		(SELECT COUNT(*) FROM referrals WHERE inviter_id = u.id) AS sub_user_count,
		COALESCE((SELECT SUM(commission_amount) FROM agent_commissions WHERE agent_id = u.id), 0) AS total_commission,
		COALESCE((SELECT SUM(commission_amount) FROM agent_commissions WHERE agent_id = u.id AND status = 'pending'), 0) AS pending_commission
		FROM users u
		LEFT JOIN agent_profiles ap ON ap.user_id = u.id
		WHERE (u.is_agent = true OR u.agent_status != '') AND u.deleted_at IS NULL`

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
		var activationFeePaidAt sql.NullTime
		if err := rows.Scan(&a.ID, &a.Email, &a.Username, &a.IsAgent, &a.AgentStatus, &a.CommissionRate,
			&a.AgentNote, &approvedAt, &a.InviteCode, &a.CreatedAt,
			&a.RealName, &a.Phone, &a.IdentityStatus, &a.ContractStatus, &activationFeePaidAt, &a.IsFrozen,
			&a.FrozenReason, &a.FrozenBalance, &a.WithdrawableBalance, &a.TotalWithdrawn,
			&a.SubUserCount, &a.TotalCommission, &a.PendingCommission); err != nil {
			return nil, nil, err
		}
		if approvedAt.Valid {
			a.ApprovedAt = &approvedAt.Time
		}
		if activationFeePaidAt.Valid {
			a.ActivationFeePaidAt = &activationFeePaidAt.Time
		}
		results = append(results, a)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return results, paginationResultFromTotal(total, params), nil
}

// GetAgentByUserID finds the agent (inviter) for a given user via referrals table.
// Returns agentID, per-user commission rate (may be nil), and error.
func (r *agentRepository) GetAgentByUserID(ctx context.Context, userID int64) (int64, *float64, error) {
	var agentID int64
	var perUserRate sql.NullFloat64
	err := r.db.QueryRowContext(ctx,
		`SELECT ref.inviter_id, ref.commission_rate FROM referrals ref
		 JOIN users u ON u.id = ref.inviter_id
		 WHERE ref.invitee_id = $1 AND u.is_agent = true AND u.agent_status = 'approved'`,
		userID).Scan(&agentID, &perUserRate)
	if err == sql.ErrNoRows {
		return 0, nil, nil
	}
	if err != nil {
		return 0, nil, err
	}
	if perUserRate.Valid {
		v := perUserRate.Float64
		return agentID, &v, nil
	}
	return agentID, nil, nil
}

func (r *agentRepository) GetProfile(ctx context.Context, userID int64) (*service.AgentProfile, error) {
	query := `SELECT user_id, real_name, id_card_no, phone, identity_status, identity_submitted_at,
		contract_status, contract_version, contract_signed_at, contract_ip, contract_signature_data, activation_order_id,
		activation_fee_paid_at, frozen_balance, withdrawable_balance, total_withdrawn,
		is_frozen, frozen_reason, created_at, updated_at
		FROM agent_profiles WHERE user_id = $1`

	var profile service.AgentProfile
	var identitySubmittedAt sql.NullTime
	var contractSignedAt sql.NullTime
	var activationOrderID sql.NullInt64
	var activationFeePaidAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.UserID, &profile.RealName, &profile.IDCardNo, &profile.Phone,
		&profile.IdentityStatus, &identitySubmittedAt, &profile.ContractStatus, &profile.ContractVersion,
		&contractSignedAt, &profile.ContractIP, &profile.ContractSignatureData, &activationOrderID, &activationFeePaidAt,
		&profile.FrozenBalance, &profile.WithdrawableBalance, &profile.TotalWithdrawn,
		&profile.IsFrozen, &profile.FrozenReason, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return &service.AgentProfile{
			UserID:          userID,
			IdentityStatus:  service.AgentIdentityStatusUnsubmitted,
			ContractStatus:  service.AgentContractStatusUnsigned,
			ContractVersion: "v1",
		}, nil
	}
	if err != nil {
		return nil, err
	}
	if identitySubmittedAt.Valid {
		profile.IdentitySubmittedAt = &identitySubmittedAt.Time
	}
	if contractSignedAt.Valid {
		profile.ContractSignedAt = &contractSignedAt.Time
	}
	if activationOrderID.Valid {
		profile.ActivationOrderID = &activationOrderID.Int64
	}
	if activationFeePaidAt.Valid {
		profile.ActivationFeePaidAt = &activationFeePaidAt.Time
	}
	return &profile, nil
}

func (r *agentRepository) UpsertProfile(ctx context.Context, profile *service.AgentProfile) error {
	query := `INSERT INTO agent_profiles (
			user_id, real_name, id_card_no, phone, identity_status, identity_submitted_at,
			contract_status, contract_version, contract_signed_at, contract_ip, contract_signature_data,
			frozen_balance, withdrawable_balance, total_withdrawn, is_frozen, frozen_reason, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, NOW(), NOW()
		)
		ON CONFLICT (user_id) DO UPDATE SET
			real_name = EXCLUDED.real_name,
			id_card_no = EXCLUDED.id_card_no,
			phone = EXCLUDED.phone,
			identity_status = EXCLUDED.identity_status,
			identity_submitted_at = EXCLUDED.identity_submitted_at,
			contract_status = EXCLUDED.contract_status,
			contract_version = EXCLUDED.contract_version,
			contract_signed_at = EXCLUDED.contract_signed_at,
			contract_ip = EXCLUDED.contract_ip,
			contract_signature_data = EXCLUDED.contract_signature_data,
			updated_at = NOW()`
	_, err := r.db.ExecContext(ctx, query,
		profile.UserID, profile.RealName, profile.IDCardNo, profile.Phone,
		profile.IdentityStatus, profile.IdentitySubmittedAt,
		profile.ContractStatus, profile.ContractVersion, profile.ContractSignedAt, profile.ContractIP, profile.ContractSignatureData,
		profile.FrozenBalance, profile.WithdrawableBalance, profile.TotalWithdrawn,
		profile.IsFrozen, profile.FrozenReason,
	)
	return err
}

func (r *agentRepository) MarkActivationFeePaid(ctx context.Context, userID int64, orderID int64) error {
	query := `INSERT INTO agent_profiles (user_id, activation_order_id, activation_fee_paid_at, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE SET activation_order_id = EXCLUDED.activation_order_id, activation_fee_paid_at = NOW(), updated_at = NOW()`
	_, err := r.db.ExecContext(ctx, query, userID, orderID)
	return err
}

func (r *agentRepository) SetAgentFrozen(ctx context.Context, userID int64, frozen bool, reason string) error {
	query := `INSERT INTO agent_profiles (user_id, is_frozen, frozen_reason, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE SET is_frozen = EXCLUDED.is_frozen, frozen_reason = EXCLUDED.frozen_reason, updated_at = NOW()`
	_, err := r.db.ExecContext(ctx, query, userID, frozen, reason)
	return err
}

func (r *agentRepository) AddWalletLog(ctx context.Context, userID int64, balanceType, changeType string, amount float64, relatedUserID *int64, relatedOrderID *int64, remark string, unlockAt *time.Time) error {
	query := `INSERT INTO agent_wallet_logs (user_id, balance_type, change_type, amount, related_user_id, related_order_id, unlock_at, remark, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`
	_, err := r.db.ExecContext(ctx, query, userID, balanceType, changeType, amount, relatedUserID, relatedOrderID, unlockAt, remark)
	return err
}

// GetParentAgent finds the parent agent of a given agent (the agent's inviter who is also an approved agent).
// Returns parentAgentID, per-user commission rate, and error.
func (r *agentRepository) GetParentAgent(ctx context.Context, agentID int64) (int64, *float64, error) {
	var parentID int64
	var perUserRate sql.NullFloat64
	err := r.db.QueryRowContext(ctx,
		`SELECT ref.inviter_id, ref.commission_rate FROM referrals ref
		 JOIN users u ON u.id = ref.inviter_id
		 WHERE ref.invitee_id = $1 AND u.is_agent = true AND u.agent_status = 'approved'`,
		agentID).Scan(&parentID, &perUserRate)
	if err == sql.ErrNoRows {
		return 0, nil, nil
	}
	if err != nil {
		return 0, nil, err
	}
	if perUserRate.Valid {
		v := perUserRate.Float64
		return parentID, &v, nil
	}
	return parentID, nil, nil
}

// UpdateReferralCommissionRate sets the per-user commission rate for a referral relationship.
func (r *agentRepository) UpdateReferralCommissionRate(ctx context.Context, inviterID, inviteeID int64, rate float64) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE referrals SET commission_rate = $1 WHERE inviter_id = $2 AND invitee_id = $3`,
		rate, inviterID, inviteeID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// UpdateReferralInviter changes the inviter (parent) of a user. Upserts the referral record.
func (r *agentRepository) UpdateReferralInviter(ctx context.Context, inviteeID, newInviterID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO referrals (inviter_id, invitee_id, commission_rate, created_at)
		 VALUES ($1, $2, NULL, NOW())
		 ON CONFLICT (invitee_id) DO UPDATE SET inviter_id = $1, commission_rate = NULL`,
		newInviterID, inviteeID)
	return err
}
