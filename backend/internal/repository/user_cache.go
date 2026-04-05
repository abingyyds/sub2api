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
	userCachePrefix = "user:info:"
	userCacheTTL    = 5 * time.Minute
)

func userCacheKey(userID int64) string {
	return fmt.Sprintf("%s%d", userCachePrefix, userID)
}

type userCache struct {
	rdb *redis.Client
}

func NewUserCache(rdb *redis.Client) service.UserCache {
	return &userCache{rdb: rdb}
}

func (c *userCache) Get(ctx context.Context, userID int64) (*service.User, error) {
	key := userCacheKey(userID)
	val, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var user service.User
	if err := json.Unmarshal(val, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *userCache) Set(ctx context.Context, userID int64, user *service.User) error {
	if user == nil {
		return nil
	}

	key := userCacheKey(userID)
	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, payload, userCacheTTL).Err()
}

func (c *userCache) Delete(ctx context.Context, userID int64) error {
	key := userCacheKey(userID)
	return c.rdb.Del(ctx, key).Err()
}
