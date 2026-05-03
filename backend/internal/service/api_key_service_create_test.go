package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type apiKeyCreateRepoStub struct {
	created *APIKey
}

func (s *apiKeyCreateRepoStub) Create(ctx context.Context, key *APIKey) error {
	clone := *key
	s.created = &clone
	key.ID = 1
	return nil
}

func (s *apiKeyCreateRepoStub) GetByID(ctx context.Context, id int64) (*APIKey, error) {
	panic("unexpected GetByID call")
}

func (s *apiKeyCreateRepoStub) GetKeyAndOwnerID(ctx context.Context, id int64) (string, int64, error) {
	panic("unexpected GetKeyAndOwnerID call")
}

func (s *apiKeyCreateRepoStub) GetByKey(ctx context.Context, key string) (*APIKey, error) {
	panic("unexpected GetByKey call")
}

func (s *apiKeyCreateRepoStub) GetByKeyForAuth(ctx context.Context, key string) (*APIKey, error) {
	panic("unexpected GetByKeyForAuth call")
}

func (s *apiKeyCreateRepoStub) Update(ctx context.Context, key *APIKey) error {
	panic("unexpected Update call")
}

func (s *apiKeyCreateRepoStub) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
}

func (s *apiKeyCreateRepoStub) ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams) ([]APIKey, *pagination.PaginationResult, error) {
	panic("unexpected ListByUserID call")
}

func (s *apiKeyCreateRepoStub) VerifyOwnership(ctx context.Context, userID int64, apiKeyIDs []int64) ([]int64, error) {
	panic("unexpected VerifyOwnership call")
}

func (s *apiKeyCreateRepoStub) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	panic("unexpected CountByUserID call")
}

func (s *apiKeyCreateRepoStub) ExistsByKey(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (s *apiKeyCreateRepoStub) ListByGroupID(ctx context.Context, groupID int64, params pagination.PaginationParams) ([]APIKey, *pagination.PaginationResult, error) {
	panic("unexpected ListByGroupID call")
}

func (s *apiKeyCreateRepoStub) SearchAPIKeys(ctx context.Context, userID int64, keyword string, limit int) ([]APIKey, error) {
	panic("unexpected SearchAPIKeys call")
}

func (s *apiKeyCreateRepoStub) ClearGroupIDByGroupID(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected ClearGroupIDByGroupID call")
}

func (s *apiKeyCreateRepoStub) CountByGroupID(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected CountByGroupID call")
}

func (s *apiKeyCreateRepoStub) ListKeysByUserID(ctx context.Context, userID int64) ([]string, error) {
	panic("unexpected ListKeysByUserID call")
}

func (s *apiKeyCreateRepoStub) ListKeysByGroupID(ctx context.Context, groupID int64) ([]string, error) {
	panic("unexpected ListKeysByGroupID call")
}

type apiKeyCreateUserRepoStub struct{}

func (s apiKeyCreateUserRepoStub) Create(ctx context.Context, user *User) error {
	panic("unexpected Create call")
}

func (s apiKeyCreateUserRepoStub) GetByID(ctx context.Context, id int64) (*User, error) {
	return &User{ID: id, Status: StatusActive}, nil
}

func (s apiKeyCreateUserRepoStub) GetByEmail(ctx context.Context, email string) (*User, error) {
	panic("unexpected GetByEmail call")
}

func (s apiKeyCreateUserRepoStub) GetFirstAdmin(ctx context.Context) (*User, error) {
	panic("unexpected GetFirstAdmin call")
}

func (s apiKeyCreateUserRepoStub) Update(ctx context.Context, user *User) error {
	panic("unexpected Update call")
}

func (s apiKeyCreateUserRepoStub) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
}

func (s apiKeyCreateUserRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s apiKeyCreateUserRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters UserListFilters) ([]User, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
}

func (s apiKeyCreateUserRepoStub) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	panic("unexpected UpdateBalance call")
}

func (s apiKeyCreateUserRepoStub) DeductBalance(ctx context.Context, id int64, amount float64) error {
	panic("unexpected DeductBalance call")
}

func (s apiKeyCreateUserRepoStub) UpdateConcurrency(ctx context.Context, id int64, amount int) error {
	panic("unexpected UpdateConcurrency call")
}

func (s apiKeyCreateUserRepoStub) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	panic("unexpected ExistsByEmail call")
}

func (s apiKeyCreateUserRepoStub) RemoveGroupFromAllowedGroups(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected RemoveGroupFromAllowedGroups call")
}

func (s apiKeyCreateUserRepoStub) UpdateTotpSecret(ctx context.Context, userID int64, encryptedSecret *string) error {
	panic("unexpected UpdateTotpSecret call")
}

func (s apiKeyCreateUserRepoStub) EnableTotp(ctx context.Context, userID int64) error {
	panic("unexpected EnableTotp call")
}

func (s apiKeyCreateUserRepoStub) DisableTotp(ctx context.Context, userID int64) error {
	panic("unexpected DisableTotp call")
}

func (s apiKeyCreateUserRepoStub) GetByInviteCode(ctx context.Context, code string) (*User, error) {
	panic("unexpected GetByInviteCode call")
}

func (s apiKeyCreateUserRepoStub) GetDiscoverySourceStats(ctx context.Context, startTime, endTime time.Time) ([]DiscoverySourceStat, error) {
	panic("unexpected GetDiscoverySourceStats call")
}

func (s apiKeyCreateUserRepoStub) ClearExpiredInitialBalances(ctx context.Context) ([]int64, error) {
	panic("unexpected ClearExpiredInitialBalances call")
}

func TestAPIKeyService_Create_DefaultsBlankNameToRandomShortName(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
	}{
		{name: "empty", inputName: ""},
		{name: "spaces", inputName: "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &apiKeyCreateRepoStub{}
			svc := NewAPIKeyService(repo, apiKeyCreateUserRepoStub{}, nil, nil, nil, &config.Config{})

			key, err := svc.Create(context.Background(), 7, CreateAPIKeyRequest{Name: tt.inputName})
			require.NoError(t, err)
			require.NotNil(t, repo.created)
			require.Len(t, repo.created.Name, apiKeyDefaultNameLen)
			require.Equal(t, repo.created.Name, key.Name)
		})
	}
}

func TestAPIKeyService_Create_TrimsProvidedName(t *testing.T) {
	repo := &apiKeyCreateRepoStub{}
	svc := NewAPIKeyService(repo, apiKeyCreateUserRepoStub{}, nil, nil, nil, &config.Config{})

	key, err := svc.Create(context.Background(), 7, CreateAPIKeyRequest{Name: "  production  "})
	require.NoError(t, err)
	require.Equal(t, "production", key.Name)
}

var _ APIKeyRepository = (*apiKeyCreateRepoStub)(nil)
var _ UserRepository = (*apiKeyCreateUserRepoStub)(nil)
