package service

import (
	"context"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var (
	ErrQuotaPackageInsufficient = infraerrors.Forbidden("QUOTA_PACKAGE_INSUFFICIENT", "quota package balance is insufficient")
	ErrQuotaPackageInvalid      = infraerrors.BadRequest("QUOTA_PACKAGE_INVALID", "quota package configuration is invalid")
)

type QuotaPackageRepository interface {
	CreateFromOrder(ctx context.Context, userID, groupID, orderID int64, quotaUSD float64, expiresAt time.Time) error
	GetAvailableTotal(ctx context.Context, userID, groupID int64) (float64, error)
	Deduct(ctx context.Context, userID, groupID int64, amount float64) error
}
