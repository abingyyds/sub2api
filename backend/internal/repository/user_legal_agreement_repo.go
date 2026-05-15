package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type userLegalAgreementRepository struct {
	db *sql.DB
}

func NewUserLegalAgreementRepository(db *sql.DB) service.UserLegalAgreementRepository {
	return &userLegalAgreementRepository{db: db}
}

func (r *userLegalAgreementRepository) Upsert(ctx context.Context, agreement *service.UserLegalAgreement) error {
	if agreement == nil {
		return nil
	}
	if r == nil || r.db == nil {
		return service.ErrServiceUnavailable
	}

	query := `
		INSERT INTO user_legal_agreements (
			user_id,
			terms_accepted_at,
			privacy_accepted_at,
			api_terms_accepted_at,
			terms_version,
			privacy_version,
			api_terms_version,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			terms_accepted_at = EXCLUDED.terms_accepted_at,
			privacy_accepted_at = EXCLUDED.privacy_accepted_at,
			api_terms_accepted_at = EXCLUDED.api_terms_accepted_at,
			terms_version = EXCLUDED.terms_version,
			privacy_version = EXCLUDED.privacy_version,
			api_terms_version = EXCLUDED.api_terms_version,
			updated_at = NOW()
		RETURNING created_at, updated_at
	`
	if err := r.db.QueryRowContext(
		ctx,
		query,
		agreement.UserID,
		agreement.TermsAcceptedAt,
		agreement.PrivacyAcceptedAt,
		agreement.ApiTermsAcceptedAt,
		agreement.TermsVersion,
		agreement.PrivacyVersion,
		agreement.ApiTermsVersion,
	).Scan(&agreement.CreatedAt, &agreement.UpdatedAt); err != nil {
		return fmt.Errorf("upsert user legal agreement: %w", err)
	}
	return nil
}

func (r *userLegalAgreementRepository) GetByUserID(ctx context.Context, userID int64) (*service.UserLegalAgreement, error) {
	if r == nil || r.db == nil {
		return nil, service.ErrServiceUnavailable
	}

	query := `
		SELECT
			user_id,
			terms_accepted_at,
			privacy_accepted_at,
			api_terms_accepted_at,
			terms_version,
			privacy_version,
			api_terms_version,
			created_at,
			updated_at
		FROM user_legal_agreements
		WHERE user_id = $1
	`
	var agreement service.UserLegalAgreement
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&agreement.UserID,
		&agreement.TermsAcceptedAt,
		&agreement.PrivacyAcceptedAt,
		&agreement.ApiTermsAcceptedAt,
		&agreement.TermsVersion,
		&agreement.PrivacyVersion,
		&agreement.ApiTermsVersion,
		&agreement.CreatedAt,
		&agreement.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrLegalAgreementNotFound
		}
		return nil, fmt.Errorf("get user legal agreement: %w", err)
	}
	return &agreement, nil
}
