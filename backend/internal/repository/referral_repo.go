package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type referralRepository struct {
	db *sql.DB
}

func NewReferralRepository(sqlDB *sql.DB) service.ReferralRepository {
	return &referralRepository{db: sqlDB}
}

func (r *referralRepository) Create(ctx context.Context, referral *service.Referral) error {
	query := `INSERT INTO referrals (inviter_id, invitee_id, reward_status, reward_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		referral.InviterID, referral.InviteeID, referral.RewardStatus, referral.RewardAmount,
	).Scan(&referral.ID, &referral.CreatedAt, &referral.UpdatedAt)
}

func (r *referralRepository) GetByInviteeID(ctx context.Context, inviteeID int64) (*service.Referral, error) {
	query := `SELECT id, inviter_id, invitee_id, reward_status, reward_amount, rewarded_at, created_at, updated_at
		FROM referrals WHERE invitee_id = $1`
	return r.scanReferral(r.db.QueryRowContext(ctx, query, inviteeID))
}

func (r *referralRepository) UpdateRewardStatus(ctx context.Context, id int64, status string, amount float64) error {
	var query string
	var args []any
	if status == service.ReferralRewardRewarded {
		now := time.Now()
		query = `UPDATE referrals SET reward_status = $1, reward_amount = $2, rewarded_at = $3, updated_at = $3 WHERE id = $4`
		args = []any{status, amount, now, id}
	} else {
		query = `UPDATE referrals SET reward_status = $1, reward_amount = $2, updated_at = NOW() WHERE id = $3`
		args = []any{status, amount, id}
	}
	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return service.ErrReferralNotFound
	}
	return nil
}

func (r *referralRepository) ListByInviterID(ctx context.Context, inviterID int64, params pagination.PaginationParams) ([]service.Referral, *pagination.PaginationResult, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM referrals WHERE inviter_id = $1`, inviterID).Scan(&total); err != nil {
		return nil, nil, err
	}
	query := `SELECT r.id, r.inviter_id, r.invitee_id, r.reward_status, r.reward_amount, r.rewarded_at, r.created_at, r.updated_at,
		u.email AS invitee_email
		FROM referrals r LEFT JOIN users u ON u.id = r.invitee_id
		WHERE r.inviter_id = $1 ORDER BY r.id DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, inviterID, params.Limit(), params.Offset())
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()
	var results []service.Referral
	for rows.Next() {
		ref, err := r.scanReferralWithJoin(rows)
		if err != nil {
			return nil, nil, err
		}
		results = append(results, *ref)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return results, paginationResultFromTotal(total, params), nil
}

func (r *referralRepository) ListAll(ctx context.Context, params pagination.PaginationParams, search string) ([]service.Referral, *pagination.PaginationResult, error) {
	var total int64
	if search != "" {
		p := "%" + search + "%"
		err := r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM referrals r
			LEFT JOIN users inviter ON inviter.id = r.inviter_id
			LEFT JOIN users invitee ON invitee.id = r.invitee_id
			WHERE inviter.email ILIKE $1 OR invitee.email ILIKE $2`, p, p).Scan(&total)
		if err != nil {
			return nil, nil, err
		}
	} else {
		if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM referrals`).Scan(&total); err != nil {
			return nil, nil, err
		}
	}

	base := `SELECT r.id, r.inviter_id, r.invitee_id, r.reward_status, r.reward_amount, r.rewarded_at, r.created_at, r.updated_at,
		inviter.email AS inviter_email, invitee.email AS invitee_email
		FROM referrals r
		LEFT JOIN users inviter ON inviter.id = r.inviter_id
		LEFT JOIN users invitee ON invitee.id = r.invitee_id`
	var rows *sql.Rows
	var err error
	if search != "" {
		p := "%" + search + "%"
		rows, err = r.db.QueryContext(ctx,
			fmt.Sprintf(`%s WHERE inviter.email ILIKE $1 OR invitee.email ILIKE $2 ORDER BY r.id DESC LIMIT $3 OFFSET $4`, base),
			p, p, params.Limit(), params.Offset())
	} else {
		rows, err = r.db.QueryContext(ctx,
			fmt.Sprintf(`%s ORDER BY r.id DESC LIMIT $1 OFFSET $2`, base),
			params.Limit(), params.Offset())
	}
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []service.Referral
	for rows.Next() {
		ref, err := r.scanReferralWithFullJoin(rows)
		if err != nil {
			return nil, nil, err
		}
		results = append(results, *ref)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return results, paginationResultFromTotal(total, params), nil
}

func (r *referralRepository) GetStatsByInviterID(ctx context.Context, inviterID int64) (*service.ReferralStats, error) {
	query := `SELECT COUNT(*),
		COUNT(*) FILTER (WHERE reward_status = 'rewarded'),
		COUNT(*) FILTER (WHERE reward_status = 'pending'),
		COALESCE(SUM(reward_amount) FILTER (WHERE reward_status = 'rewarded'), 0)
		FROM referrals WHERE inviter_id = $1`
	s := &service.ReferralStats{}
	if err := r.db.QueryRowContext(ctx, query, inviterID).Scan(&s.TotalInvites, &s.RewardedCount, &s.PendingCount, &s.TotalRewarded); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *referralRepository) GetGlobalStats(ctx context.Context) (*service.ReferralStats, error) {
	query := `SELECT COUNT(*),
		COUNT(*) FILTER (WHERE reward_status = 'rewarded'),
		COUNT(*) FILTER (WHERE reward_status = 'pending'),
		COALESCE(SUM(reward_amount) FILTER (WHERE reward_status = 'rewarded'), 0)
		FROM referrals`
	s := &service.ReferralStats{}
	if err := r.db.QueryRowContext(ctx, query).Scan(&s.TotalInvites, &s.RewardedCount, &s.PendingCount, &s.TotalRewarded); err != nil {
		return nil, err
	}
	return s, nil
}

// --- scan helpers ---

func (r *referralRepository) scanReferral(row *sql.Row) (*service.Referral, error) {
	ref := &service.Referral{}
	var rewardedAt sql.NullTime
	err := row.Scan(&ref.ID, &ref.InviterID, &ref.InviteeID, &ref.RewardStatus,
		&ref.RewardAmount, &rewardedAt, &ref.CreatedAt, &ref.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrReferralNotFound
		}
		return nil, err
	}
	if rewardedAt.Valid {
		ref.RewardedAt = &rewardedAt.Time
	}
	return ref, nil
}

func (r *referralRepository) scanReferralWithJoin(rows *sql.Rows) (*service.Referral, error) {
	ref := &service.Referral{}
	var rewardedAt sql.NullTime
	var inviteeEmail sql.NullString
	err := rows.Scan(&ref.ID, &ref.InviterID, &ref.InviteeID, &ref.RewardStatus,
		&ref.RewardAmount, &rewardedAt, &ref.CreatedAt, &ref.UpdatedAt, &inviteeEmail)
	if err != nil {
		return nil, err
	}
	if rewardedAt.Valid {
		ref.RewardedAt = &rewardedAt.Time
	}
	if inviteeEmail.Valid {
		ref.InviteeEmail = inviteeEmail.String
	}
	return ref, nil
}

func (r *referralRepository) scanReferralWithFullJoin(rows *sql.Rows) (*service.Referral, error) {
	ref := &service.Referral{}
	var rewardedAt sql.NullTime
	var inviterEmail, inviteeEmail sql.NullString
	err := rows.Scan(&ref.ID, &ref.InviterID, &ref.InviteeID, &ref.RewardStatus,
		&ref.RewardAmount, &rewardedAt, &ref.CreatedAt, &ref.UpdatedAt,
		&inviterEmail, &inviteeEmail)
	if err != nil {
		return nil, err
	}
	if rewardedAt.Valid {
		ref.RewardedAt = &rewardedAt.Time
	}
	if inviterEmail.Valid {
		ref.InviterEmail = inviterEmail.String
	}
	if inviteeEmail.Valid {
		ref.InviteeEmail = inviteeEmail.String
	}
	return ref, nil
}
