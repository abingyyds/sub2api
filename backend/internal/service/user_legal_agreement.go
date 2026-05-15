package service

import (
	"context"
	"errors"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	LegalTermsVersion    = "2026-05-15"
	LegalPrivacyVersion  = "2026-05-15"
	LegalApiTermsVersion = "2026-05-15"
)

var (
	ErrLegalAgreementNotFound = infraerrors.NotFound("LEGAL_AGREEMENT_NOT_FOUND", "legal agreement not found")
	ErrLegalAgreementRequired = infraerrors.Forbidden("LEGAL_AGREEMENT_REQUIRED", "You must read and accept the User Agreement and Privacy Policy before registering, creating API keys, or using the API.")
)

// UserLegalAgreement stores the latest legal confirmations accepted by a user.
type UserLegalAgreement struct {
	UserID             int64
	TermsAcceptedAt    time.Time
	PrivacyAcceptedAt  time.Time
	ApiTermsAcceptedAt time.Time
	TermsVersion       string
	PrivacyVersion     string
	ApiTermsVersion    string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type UserLegalAgreementRepository interface {
	Upsert(ctx context.Context, agreement *UserLegalAgreement) error
	GetByUserID(ctx context.Context, userID int64) (*UserLegalAgreement, error)
}

func NewCurrentLegalAgreement(userID int64) *UserLegalAgreement {
	now := time.Now()
	return &UserLegalAgreement{
		UserID:             userID,
		TermsAcceptedAt:    now,
		PrivacyAcceptedAt:  now,
		ApiTermsAcceptedAt: now,
		TermsVersion:       LegalTermsVersion,
		PrivacyVersion:     LegalPrivacyVersion,
		ApiTermsVersion:    LegalApiTermsVersion,
	}
}

func (a *UserLegalAgreement) IsCurrent() bool {
	if a == nil {
		return false
	}
	return a.TermsVersion == LegalTermsVersion &&
		a.PrivacyVersion == LegalPrivacyVersion &&
		a.ApiTermsVersion == LegalApiTermsVersion &&
		!a.TermsAcceptedAt.IsZero() &&
		!a.PrivacyAcceptedAt.IsZero() &&
		!a.ApiTermsAcceptedAt.IsZero()
}

func IsLegalAgreementMissing(err error) bool {
	return errors.Is(err, ErrLegalAgreementNotFound)
}
