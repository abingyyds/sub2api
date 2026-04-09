//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type subscriptionGroupRepoStub struct {
	group *Group
}

func (s *subscriptionGroupRepoStub) Create(ctx context.Context, group *Group) error {
	panic("unexpected Create call")
}

func (s *subscriptionGroupRepoStub) GetByID(ctx context.Context, id int64) (*Group, error) {
	if s.group == nil || s.group.ID != id {
		return nil, ErrGroupNotFound
	}
	clone := *s.group
	return &clone, nil
}

func (s *subscriptionGroupRepoStub) GetByIDLite(ctx context.Context, id int64) (*Group, error) {
	return s.GetByID(ctx, id)
}

func (s *subscriptionGroupRepoStub) Update(ctx context.Context, group *Group) error {
	panic("unexpected Update call")
}

func (s *subscriptionGroupRepoStub) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
}

func (s *subscriptionGroupRepoStub) DeleteCascade(ctx context.Context, id int64) ([]int64, error) {
	panic("unexpected DeleteCascade call")
}

func (s *subscriptionGroupRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *subscriptionGroupRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, status, search string, isExclusive *bool) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
}

func (s *subscriptionGroupRepoStub) ListActive(ctx context.Context) ([]Group, error) {
	panic("unexpected ListActive call")
}

func (s *subscriptionGroupRepoStub) ListActiveByPlatform(ctx context.Context, platform string) ([]Group, error) {
	panic("unexpected ListActiveByPlatform call")
}

func (s *subscriptionGroupRepoStub) ExistsByName(ctx context.Context, name string) (bool, error) {
	panic("unexpected ExistsByName call")
}

func (s *subscriptionGroupRepoStub) GetAccountCount(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected GetAccountCount call")
}

func (s *subscriptionGroupRepoStub) DeleteAccountGroupsByGroupID(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected DeleteAccountGroupsByGroupID call")
}

type subscriptionRepoStub struct {
	created *UserSubscription
}

func (s *subscriptionRepoStub) Create(ctx context.Context, sub *UserSubscription) error {
	if sub == nil {
		return nil
	}
	clone := *sub
	clone.ID = 101
	s.created = &clone
	sub.ID = clone.ID
	return nil
}

func (s *subscriptionRepoStub) GetByID(ctx context.Context, id int64) (*UserSubscription, error) {
	if s.created == nil || s.created.ID != id {
		return nil, ErrSubscriptionNotFound
	}
	clone := *s.created
	return &clone, nil
}

func (s *subscriptionRepoStub) GetByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*UserSubscription, error) {
	return nil, ErrSubscriptionNotFound
}

func (s *subscriptionRepoStub) GetActiveByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*UserSubscription, error) {
	panic("unexpected GetActiveByUserIDAndGroupID call")
}

func (s *subscriptionRepoStub) Update(ctx context.Context, sub *UserSubscription) error {
	panic("unexpected Update call")
}

func (s *subscriptionRepoStub) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
}

func (s *subscriptionRepoStub) ListByUserID(ctx context.Context, userID int64) ([]UserSubscription, error) {
	panic("unexpected ListByUserID call")
}

func (s *subscriptionRepoStub) ListActiveByUserID(ctx context.Context, userID int64) ([]UserSubscription, error) {
	panic("unexpected ListActiveByUserID call")
}

func (s *subscriptionRepoStub) ListByGroupID(ctx context.Context, groupID int64, params pagination.PaginationParams) ([]UserSubscription, *pagination.PaginationResult, error) {
	panic("unexpected ListByGroupID call")
}

func (s *subscriptionRepoStub) List(ctx context.Context, params pagination.PaginationParams, userID, groupID *int64, status, sortBy, sortOrder string) ([]UserSubscription, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *subscriptionRepoStub) ExistsByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (bool, error) {
	return false, nil
}

func (s *subscriptionRepoStub) ExtendExpiry(ctx context.Context, subscriptionID int64, newExpiresAt time.Time) error {
	panic("unexpected ExtendExpiry call")
}

func (s *subscriptionRepoStub) UpdateStatus(ctx context.Context, subscriptionID int64, status string) error {
	panic("unexpected UpdateStatus call")
}

func (s *subscriptionRepoStub) UpdateNotes(ctx context.Context, subscriptionID int64, notes string) error {
	panic("unexpected UpdateNotes call")
}

func (s *subscriptionRepoStub) ActivateWindows(ctx context.Context, id int64, start time.Time) error {
	panic("unexpected ActivateWindows call")
}

func (s *subscriptionRepoStub) ResetDailyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	panic("unexpected ResetDailyUsage call")
}

func (s *subscriptionRepoStub) ResetWeeklyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	panic("unexpected ResetWeeklyUsage call")
}

func (s *subscriptionRepoStub) ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	panic("unexpected ResetMonthlyUsage call")
}

func (s *subscriptionRepoStub) IncrementUsage(ctx context.Context, id int64, costUSD float64) error {
	panic("unexpected IncrementUsage call")
}

func (s *subscriptionRepoStub) BatchUpdateExpiredStatus(ctx context.Context) (int64, error) {
	panic("unexpected BatchUpdateExpiredStatus call")
}

type subscriptionCacheStub struct {
	deletedUserIDs []int64
}

func (s *subscriptionCacheStub) GetActiveByUserID(ctx context.Context, userID int64) ([]UserSubscription, error) {
	panic("unexpected GetActiveByUserID call")
}

func (s *subscriptionCacheStub) SetActiveByUserID(ctx context.Context, userID int64, subscriptions []UserSubscription) error {
	panic("unexpected SetActiveByUserID call")
}

func (s *subscriptionCacheStub) DeleteByUserID(ctx context.Context, userID int64) error {
	s.deletedUserIDs = append(s.deletedUserIDs, userID)
	return nil
}

func TestSubscriptionServiceAssignOrExtendSubscription_NewSubscriptionDeletesActiveCache(t *testing.T) {
	t.Parallel()

	groupRepo := &subscriptionGroupRepoStub{
		group: &Group{
			ID:                  11,
			Status:              StatusActive,
			SubscriptionType:    SubscriptionTypeSubscription,
			DefaultValidityDays: 30,
		},
	}
	subRepo := &subscriptionRepoStub{}
	cache := &subscriptionCacheStub{}
	svc := &SubscriptionService{
		groupRepo:         groupRepo,
		userSubRepo:       subRepo,
		subscriptionCache: cache,
	}

	sub, extended, err := svc.AssignOrExtendSubscription(context.Background(), &AssignSubscriptionInput{
		UserID:       7,
		GroupID:      11,
		ValidityDays: 30,
		Notes:        "payment order test",
	})
	require.NoError(t, err)
	require.False(t, extended)
	require.NotNil(t, sub)
	require.Equal(t, int64(101), sub.ID)
	require.Equal(t, []int64{7}, cache.deletedUserIDs)
}
