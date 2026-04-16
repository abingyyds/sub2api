//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type userCacheStub struct {
	deletedIDs []int64
	deleteErr  error
}

func (s *userCacheStub) Get(ctx context.Context, userID int64) (*User, error) {
	return nil, ErrUserNotFound
}

func (s *userCacheStub) Set(ctx context.Context, userID int64, user *User) error {
	return nil
}

func (s *userCacheStub) Delete(ctx context.Context, userID int64) error {
	s.deletedIDs = append(s.deletedIDs, userID)
	return s.deleteErr
}

type updateUserRepoStub struct {
	*userRepoStub
	updateErr error
	updated   []*User
}

func (s *updateUserRepoStub) Update(ctx context.Context, user *User) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	if user == nil {
		return nil
	}
	clone := *user
	s.updated = append(s.updated, &clone)
	if s.userRepoStub != nil {
		s.userRepoStub.user = &clone
	}
	return nil
}

func TestAdminService_UpdateUser_InvalidatesUserCacheForProfileChanges(t *testing.T) {
	baseRepo := &userRepoStub{
		user: &User{ID: 7, Email: "alice@example.com", Username: "alice", Role: RoleUser, Status: StatusActive, Concurrency: 3},
	}
	repo := &updateUserRepoStub{userRepoStub: baseRepo}
	cache := &userCacheStub{}
	invalidator := &authCacheInvalidatorStub{}
	newUsername := "alice-updated"

	svc := &adminServiceImpl{
		userRepo:             repo,
		userCache:            cache,
		authCacheInvalidator: invalidator,
	}

	user, err := svc.UpdateUser(context.Background(), 7, &UpdateUserInput{
		Username: &newUsername,
	})
	require.NoError(t, err)
	require.Equal(t, "alice-updated", user.Username)
	require.Equal(t, []int64{7}, cache.deletedIDs)
	require.Empty(t, invalidator.userIDs)
}

func TestAdminService_UpdateUser_DemoteAdminInvalidatesAuthAndUserCache(t *testing.T) {
	baseRepo := &userRepoStub{
		user: &User{ID: 7, Email: "admin@example.com", Username: "admin", Role: RoleAdmin, Status: StatusActive, Concurrency: 8},
	}
	repo := &updateUserRepoStub{userRepoStub: baseRepo}
	cache := &userCacheStub{}
	invalidator := &authCacheInvalidatorStub{}

	svc := &adminServiceImpl{
		userRepo:             repo,
		userCache:            cache,
		authCacheInvalidator: invalidator,
	}

	user, err := svc.UpdateUser(context.Background(), 7, &UpdateUserInput{
		Role:       RoleUser,
		CallerRole: RoleAdmin,
		CallerID:   99,
	})
	require.NoError(t, err)
	require.Equal(t, RoleUser, user.Role)
	require.Equal(t, []int64{7}, cache.deletedIDs)
	require.Equal(t, []int64{7}, invalidator.userIDs)
}
