package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type customModelSchedulerCache struct {
	account *Account
}

func (c *customModelSchedulerCache) GetSnapshot(context.Context, SchedulerBucket) ([]*Account, bool, error) {
	if c.account == nil {
		return nil, false, nil
	}
	return []*Account{c.account}, true, nil
}

func (c *customModelSchedulerCache) SetSnapshot(_ context.Context, _ SchedulerBucket, accounts []Account) error {
	if len(accounts) == 0 {
		c.account = nil
		return nil
	}
	account := accounts[0]
	c.account = &account
	return nil
}

func (c *customModelSchedulerCache) GetAccount(context.Context, int64) (*Account, error) {
	return c.account, nil
}

func (c *customModelSchedulerCache) SetAccount(_ context.Context, account *Account) error {
	c.account = account
	return nil
}

func (c *customModelSchedulerCache) DeleteAccount(context.Context, int64) error {
	c.account = nil
	return nil
}

func (c *customModelSchedulerCache) UpdateLastUsed(context.Context, map[int64]time.Time) error {
	return nil
}

func (c *customModelSchedulerCache) TryLockBucket(context.Context, SchedulerBucket, time.Duration) (bool, error) {
	return true, nil
}

func (c *customModelSchedulerCache) ListBuckets(context.Context) ([]SchedulerBucket, error) {
	return nil, nil
}

func (c *customModelSchedulerCache) GetOutboxWatermark(context.Context) (int64, error) {
	return 0, nil
}

func (c *customModelSchedulerCache) SetOutboxWatermark(context.Context, int64) error {
	return nil
}

func TestGatewayCustomModelBecomesRoutableAfterAccountCacheRefresh(t *testing.T) {
	ctx := context.Background()
	cache := &customModelSchedulerCache{
		account: &Account{
			ID:          1,
			Name:        "CC official API",
			Platform:    PlatformAnthropic,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Credentials: map[string]any{
				"model_mapping": map[string]any{
					"claude-opus-4-7": "claude-opus-4-7",
				},
			},
		},
	}

	snapshot := &SchedulerSnapshotService{
		cache: cache,
		cfg: &config.Config{
			RunMode: config.RunModeStandard,
		},
	}
	gateway := &GatewayService{
		schedulerSnapshot: snapshot,
		cfg:               snapshot.cfg,
	}

	_, err := gateway.SelectAccountForModel(ctx, nil, "", "claude-opus-4-8")
	require.ErrorContains(t, err, "no available accounts supporting model: claude-opus-4-8")

	updated := *cache.account
	updated.Credentials = map[string]any{
		"model_mapping": map[string]any{
			"claude-opus-4-8": "claude-opus-4-8",
		},
	}
	require.NoError(t, cache.SetAccount(ctx, &updated))

	selected, err := gateway.SelectAccountForModel(ctx, nil, "", "claude-opus-4-8")
	require.NoError(t, err)
	require.Equal(t, int64(1), selected.ID)
	require.Equal(t, "claude-opus-4-8", selected.GetMappedModel("claude-opus-4-8"))
}
