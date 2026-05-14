package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type wechatNotificationRepository struct {
	db *sql.DB
}

func NewWechatNotificationRepository(db *sql.DB) service.WechatOfficialRepository {
	return &wechatNotificationRepository{db: db}
}

func (r *wechatNotificationRepository) GetBinding(ctx context.Context, userID int64) (*service.WechatBinding, error) {
	if r == nil || r.db == nil {
		return nil, service.ErrWechatOfficialUnbound
	}
	binding := &service.WechatBinding{}
	err := scanSingleRow(ctx, r.db, `
		SELECT id, user_id, openid, enabled, bound_at, unbound_at, created_at, updated_at
		FROM user_wechat_official_bindings
		WHERE user_id = $1
		  AND unbound_at IS NULL
		ORDER BY id DESC
		LIMIT 1
	`, []any{userID},
		&binding.ID,
		&binding.UserID,
		&binding.OpenID,
		&binding.Enabled,
		&binding.BoundAt,
		&binding.UnboundAt,
		&binding.CreatedAt,
		&binding.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrWechatOfficialUnbound
		}
		return nil, fmt.Errorf("get wechat binding: %w", err)
	}
	return binding, nil
}

func (r *wechatNotificationRepository) Bind(ctx context.Context, userID int64, openID string) error {
	if r == nil || r.db == nil {
		return service.ErrWechatOfficialNotReady
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start wechat bind transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `
		UPDATE user_wechat_official_bindings
		SET enabled = FALSE,
		    unbound_at = NOW(),
		    updated_at = NOW()
		WHERE unbound_at IS NULL
		  AND (user_id = $1 OR openid = $2)
	`, userID, openID); err != nil {
		return fmt.Errorf("disable previous wechat bindings: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO user_wechat_official_bindings (
			user_id,
			openid,
			enabled,
			bound_at,
			created_at,
			updated_at
		) VALUES ($1, $2, TRUE, NOW(), NOW(), NOW())
	`, userID, openID); err != nil {
		return fmt.Errorf("bind wechat openid: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit wechat bind transaction: %w", err)
	}
	return nil
}

func (r *wechatNotificationRepository) Unbind(ctx context.Context, userID int64) error {
	if r == nil || r.db == nil {
		return service.ErrWechatOfficialUnbound
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE user_wechat_official_bindings
		SET enabled = FALSE,
		    unbound_at = NOW(),
		    updated_at = NOW()
		WHERE user_id = $1
		  AND unbound_at IS NULL
	`, userID)
	if err != nil {
		return fmt.Errorf("unbind wechat: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return service.ErrWechatOfficialUnbound
	}
	return nil
}

func (r *wechatNotificationRepository) ShouldSend(ctx context.Context, userID int64, channel, eventType, resourceKey string, cooldown time.Duration) (bool, error) {
	if r == nil || r.db == nil {
		return false, nil
	}
	if cooldown < 0 {
		cooldown = 0
	}
	cooldownSeconds := cooldown.Seconds()
	var id int64
	err := scanSingleRow(ctx, r.db, `
		INSERT INTO user_notification_events (
			user_id,
			channel,
			event_type,
			resource_key,
			last_sent_at,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, NOW(), NOW(), NOW())
		ON CONFLICT (user_id, channel, event_type, resource_key)
		DO UPDATE SET
			last_sent_at = NOW(),
			updated_at = NOW()
		WHERE user_notification_events.last_sent_at <= NOW() - ($5::double precision * INTERVAL '1 second')
		RETURNING id
	`, []any{userID, channel, eventType, resourceKey, cooldownSeconds}, &id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("reserve notification delivery: %w", err)
	}
	return true, nil
}
