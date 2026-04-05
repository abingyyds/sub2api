package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	subscriptionCachePrefix = "subscription:active:"
	subscriptionCacheTTL    = 5 * time.Minute
)

func subscriptionCacheKey(userID int64) string {
	return fmt.Sprintf("%s%d", subscriptionCachePrefix, userID)
}

type subscriptionCache struct {
	rdb *redis.Client
}

func NewSubscriptionCache(rdb *redis.Client) service.SubscriptionCache {
	return &subscriptionCache{rdb: rdb}
}

func (c *subscriptionCache) GetActiveByUserID(ctx context.Context, userID int64) ([]service.UserSubscription, error) {
	key := subscriptionCacheKey(userID)
	val, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var subscriptions []service.UserSubscription
	if err := json.Unmarshal(val, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (c *subscriptionCache) SetActiveByUserID(ctx context.Context, userID int64, subscriptions []service.UserSubscription) error {
	if subscriptions == nil {
		return nil
	}

	key := subscriptionCacheKey(userID)
	payload, err := json.Marshal(subscriptions)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, payload, subscriptionCacheTTL).Err()
}

func (c *subscriptionCache) DeleteByUserID(ctx context.Context, userID int64) error {
	key := subscriptionCacheKey(userID)
	return c.rdb.Del(ctx, key).Err()
}
