package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type quotaPackageRepository struct {
	db *sql.DB
}

func NewQuotaPackageRepository(db *sql.DB) service.QuotaPackageRepository {
	return &quotaPackageRepository{db: db}
}

func (r *quotaPackageRepository) execQuerier() sqlExecutor {
	return r.db
}

func (r *quotaPackageRepository) CreateFromOrder(ctx context.Context, userID, groupID, orderID int64, quotaUSD float64, expiresAt time.Time) error {
	if quotaUSD <= 0 || expiresAt.IsZero() {
		return service.ErrQuotaPackageInvalid
	}

	_, err := r.execQuerier().ExecContext(ctx, `
		INSERT INTO user_quota_packages (
			user_id,
			group_id,
			order_id,
			total_quota_usd,
			remaining_quota_usd,
			expires_at,
			status,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $4, $5, 'active', NOW(), NOW())
		ON CONFLICT (order_id) WHERE order_id IS NOT NULL DO NOTHING
	`, userID, groupID, orderID, quotaUSD, expiresAt)
	if err != nil {
		return fmt.Errorf("create quota package: %w", err)
	}
	return nil
}

func (r *quotaPackageRepository) GetAvailableTotal(ctx context.Context, userID, groupID int64) (float64, error) {
	var total float64
	if err := scanSingleRow(ctx, r.execQuerier(), `
		SELECT COALESCE(SUM(remaining_quota_usd), 0)
		FROM user_quota_packages
		WHERE user_id = $1
		  AND group_id = $2
		  AND status = 'active'
		  AND expires_at > NOW()
		  AND remaining_quota_usd > 0
	`, []any{userID, groupID}, &total); err != nil {
		return 0, fmt.Errorf("get quota package total: %w", err)
	}
	return total, nil
}

func (r *quotaPackageRepository) Deduct(ctx context.Context, userID, groupID int64, amount float64) error {
	if amount <= 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start quota package deduct transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	rows, err := tx.QueryContext(ctx, `
		SELECT id, remaining_quota_usd
		FROM user_quota_packages
		WHERE user_id = $1
		  AND group_id = $2
		  AND status = 'active'
		  AND expires_at > NOW()
		  AND remaining_quota_usd > 0
		ORDER BY expires_at ASC, id ASC
		FOR UPDATE
	`, userID, groupID)
	if err != nil {
		return fmt.Errorf("query quota packages for deduct: %w", err)
	}
	defer rows.Close()

	type quotaRow struct {
		id        int64
		remaining float64
	}
	var packages []quotaRow
	for rows.Next() {
		var row quotaRow
		if err := rows.Scan(&row.id, &row.remaining); err != nil {
			return fmt.Errorf("scan quota package for deduct: %w", err)
		}
		packages = append(packages, row)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate quota packages for deduct: %w", err)
	}
	if len(packages) == 0 {
		return service.ErrQuotaPackageInsufficient
	}

	totalAvailable := 0.0
	for _, pkg := range packages {
		totalAvailable += pkg.remaining
	}
	if totalAvailable+1e-9 < amount {
		return service.ErrQuotaPackageInsufficient
	}

	remainingToDeduct := amount
	for _, pkg := range packages {
		newRemaining := pkg.remaining - remainingToDeduct
		status := "active"
		if newRemaining <= 0 {
			status = "depleted"
		}
		if newRemaining < 0 {
			newRemaining = 0
		}

		if _, err := tx.ExecContext(ctx, `
			UPDATE user_quota_packages
			SET remaining_quota_usd = $1,
			    status = $2,
			    updated_at = NOW()
			WHERE id = $3
		`, newRemaining, status, pkg.id); err != nil {
			return fmt.Errorf("deduct quota package: %w", err)
		}

		if pkg.remaining >= remainingToDeduct {
			remainingToDeduct = 0
			break
		}
		remainingToDeduct -= pkg.remaining
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit quota package deduct transaction: %w", err)
	}
	return nil
}
