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
	return nil, nil
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

func (s apiKeyCreateUserRepoStub) GetBalance(ctx context.Context, id int64) (float64, error) {
	panic("unexpected GetBalance call")
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

type apiKeyCreateGroupRepoStub struct {
	group *Group
}

func (s *apiKeyCreateGroupRepoStub) Create(ctx context.Context, group *Group) error {
	panic("unexpected Create call")
}

func (s *apiKeyCreateGroupRepoStub) GetByID(ctx context.Context, id int64) (*Group, error) {
	if s.group == nil || s.group.ID != id {
		return nil, ErrGroupNotFound
	}
	clone := *s.group
	return &clone, nil
}

func (s *apiKeyCreateGroupRepoStub) GetByIDLite(ctx context.Context, id int64) (*Group, error) {
	panic("unexpected GetByIDLite call")
}

func (s *apiKeyCreateGroupRepoStub) Update(ctx context.Context, group *Group) error {
	panic("unexpected Update call")
}

func (s *apiKeyCreateGroupRepoStub) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
}

func (s *apiKeyCreateGroupRepoStub) DeleteCascade(ctx context.Context, id int64) ([]int64, error) {
	panic("unexpected DeleteCascade call")
}

func (s *apiKeyCreateGroupRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *apiKeyCreateGroupRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, status, search string, isExclusive *bool) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
}

func (s *apiKeyCreateGroupRepoStub) ListActive(ctx context.Context) ([]Group, error) {
	if s.group == nil {
		return nil, nil
	}
	clone := *s.group
	return []Group{clone}, nil
}

func (s *apiKeyCreateGroupRepoStub) ListActiveByPlatform(ctx context.Context, platform string) ([]Group, error) {
	panic("unexpected ListActiveByPlatform call")
}

func (s *apiKeyCreateGroupRepoStub) ExistsByName(ctx context.Context, name string) (bool, error) {
	panic("unexpected ExistsByName call")
}

func (s *apiKeyCreateGroupRepoStub) GetAccountCount(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected GetAccountCount call")
}

func (s *apiKeyCreateGroupRepoStub) DeleteAccountGroupsByGroupID(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected DeleteAccountGroupsByGroupID call")
}

type apiKeyCreateSubscriptionRepoStub struct {
	activeSubscriptions []UserSubscription
}

func (s *apiKeyCreateSubscriptionRepoStub) Create(ctx context.Context, sub *UserSubscription) error {
	panic("unexpected Create call")
}

func (s *apiKeyCreateSubscriptionRepoStub) GetByID(ctx context.Context, id int64) (*UserSubscription, error) {
	panic("unexpected GetByID call")
}

func (s *apiKeyCreateSubscriptionRepoStub) GetByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*UserSubscription, error) {
	panic("unexpected GetByUserIDAndGroupID call")
}

func (s *apiKeyCreateSubscriptionRepoStub) GetActiveByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*UserSubscription, error) {
	for _, sub := range s.activeSubscriptions {
		if sub.UserID == userID && sub.GroupID == groupID {
			clone := sub
			return &clone, nil
		}
	}
	return nil, ErrSubscriptionNotFound
}

func (s *apiKeyCreateSubscriptionRepoStub) Update(ctx context.Context, sub *UserSubscription) error {
	panic("unexpected Update call")
}

func (s *apiKeyCreateSubscriptionRepoStub) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
}

func (s *apiKeyCreateSubscriptionRepoStub) ListByUserID(ctx context.Context, userID int64) ([]UserSubscription, error) {
	panic("unexpected ListByUserID call")
}

func (s *apiKeyCreateSubscriptionRepoStub) ListActiveByUserID(ctx context.Context, userID int64) ([]UserSubscription, error) {
	out := make([]UserSubscription, 0, len(s.activeSubscriptions))
	for _, sub := range s.activeSubscriptions {
		if sub.UserID == userID {
			out = append(out, sub)
		}
	}
	return out, nil
}

func (s *apiKeyCreateSubscriptionRepoStub) ListByGroupID(ctx context.Context, groupID int64, params pagination.PaginationParams) ([]UserSubscription, *pagination.PaginationResult, error) {
	panic("unexpected ListByGroupID call")
}

func (s *apiKeyCreateSubscriptionRepoStub) List(ctx context.Context, params pagination.PaginationParams, userID, groupID *int64, status, sortBy, sortOrder string) ([]UserSubscription, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *apiKeyCreateSubscriptionRepoStub) ExistsByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (bool, error) {
	panic("unexpected ExistsByUserIDAndGroupID call")
}

func (s *apiKeyCreateSubscriptionRepoStub) ExtendExpiry(ctx context.Context, subscriptionID int64, newExpiresAt time.Time) error {
	panic("unexpected ExtendExpiry call")
}

func (s *apiKeyCreateSubscriptionRepoStub) UpdateStatus(ctx context.Context, subscriptionID int64, status string) error {
	panic("unexpected UpdateStatus call")
}

func (s *apiKeyCreateSubscriptionRepoStub) UpdateNotes(ctx context.Context, subscriptionID int64, notes string) error {
	panic("unexpected UpdateNotes call")
}

func (s *apiKeyCreateSubscriptionRepoStub) ActivateWindows(ctx context.Context, id int64, start time.Time) error {
	panic("unexpected ActivateWindows call")
}

func (s *apiKeyCreateSubscriptionRepoStub) ResetDailyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	panic("unexpected ResetDailyUsage call")
}

func (s *apiKeyCreateSubscriptionRepoStub) ResetWeeklyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	panic("unexpected ResetWeeklyUsage call")
}

func (s *apiKeyCreateSubscriptionRepoStub) ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	panic("unexpected ResetMonthlyUsage call")
}

func (s *apiKeyCreateSubscriptionRepoStub) IncrementUsage(ctx context.Context, id int64, costUSD float64) error {
	panic("unexpected IncrementUsage call")
}

func (s *apiKeyCreateSubscriptionRepoStub) BatchUpdateExpiredStatus(ctx context.Context) (int64, error) {
	panic("unexpected BatchUpdateExpiredStatus call")
}

type apiKeyCreateQuotaPackageRepoStub struct {
	available map[int64]float64
}

func (s *apiKeyCreateQuotaPackageRepoStub) CreateFromOrder(ctx context.Context, userID, groupID, orderID int64, quotaUSD float64, expiresAt time.Time) error {
	panic("unexpected CreateFromOrder call")
}

func (s *apiKeyCreateQuotaPackageRepoStub) GetAvailableTotal(ctx context.Context, userID, groupID int64) (float64, error) {
	return s.available[groupID], nil
}

func (s *apiKeyCreateQuotaPackageRepoStub) Deduct(ctx context.Context, userID, groupID int64, amount float64) error {
	panic("unexpected Deduct call")
}

type apiKeyCreateLegalAgreementRepoStub struct{}

func (s *apiKeyCreateLegalAgreementRepoStub) Upsert(ctx context.Context, agreement *UserLegalAgreement) error {
	return nil
}

func (s *apiKeyCreateLegalAgreementRepoStub) GetByUserID(ctx context.Context, userID int64) (*UserLegalAgreement, error) {
	return NewCurrentLegalAgreement(userID), nil
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
			svc := NewAPIKeyService(repo, apiKeyCreateUserRepoStub{}, nil, nil, nil, nil, &config.Config{})
			svc.SetLegalAgreementRepository(&apiKeyCreateLegalAgreementRepoStub{})

			key, err := svc.Create(context.Background(), 7, CreateAPIKeyRequest{Name: tt.inputName, LegalAccepted: true})
			require.NoError(t, err)
			require.NotNil(t, repo.created)
			require.Len(t, repo.created.Name, apiKeyDefaultNameLen)
			require.Equal(t, repo.created.Name, key.Name)
		})
	}
}

func TestAPIKeyService_Create_TrimsProvidedName(t *testing.T) {
	repo := &apiKeyCreateRepoStub{}
	svc := NewAPIKeyService(repo, apiKeyCreateUserRepoStub{}, nil, nil, nil, nil, &config.Config{})
	svc.SetLegalAgreementRepository(&apiKeyCreateLegalAgreementRepoStub{})

	key, err := svc.Create(context.Background(), 7, CreateAPIKeyRequest{Name: "  production  ", LegalAccepted: true})
	require.NoError(t, err)
	require.Equal(t, "production", key.Name)
}

func TestAPIKeyService_Create_AllowsQuotaPackageGroupWithoutActiveSubscription(t *testing.T) {
	groupID := int64(42)
	repo := &apiKeyCreateRepoStub{}
	groupRepo := &apiKeyCreateGroupRepoStub{
		group: &Group{
			ID:                  groupID,
			Status:              StatusActive,
			SubscriptionType:    SubscriptionTypeSubscription,
			QuotaPackageEnabled: true,
		},
	}
	subRepo := &apiKeyCreateSubscriptionRepoStub{}
	quotaRepo := &apiKeyCreateQuotaPackageRepoStub{available: map[int64]float64{groupID: 10}}
	svc := NewAPIKeyService(repo, apiKeyCreateUserRepoStub{}, groupRepo, subRepo, quotaRepo, nil, &config.Config{})
	svc.SetLegalAgreementRepository(&apiKeyCreateLegalAgreementRepoStub{})

	key, err := svc.Create(context.Background(), 7, CreateAPIKeyRequest{Name: "quota", GroupID: &groupID, LegalAccepted: true})

	require.NoError(t, err)
	require.NotNil(t, key.GroupID)
	require.Equal(t, groupID, *key.GroupID)
	require.NotNil(t, repo.created)
	require.NotNil(t, repo.created.GroupID)
	require.Equal(t, groupID, *repo.created.GroupID)
}

func TestAPIKeyService_GetAvailableGroupsIncludesQuotaPackageGroupWithoutActiveSubscription(t *testing.T) {
	groupID := int64(42)
	groupRepo := &apiKeyCreateGroupRepoStub{
		group: &Group{
			ID:                  groupID,
			Status:              StatusActive,
			SubscriptionType:    SubscriptionTypeSubscription,
			QuotaPackageEnabled: true,
		},
	}
	subRepo := &apiKeyCreateSubscriptionRepoStub{}
	quotaRepo := &apiKeyCreateQuotaPackageRepoStub{available: map[int64]float64{groupID: 10}}
	svc := NewAPIKeyService(&apiKeyCreateRepoStub{}, apiKeyCreateUserRepoStub{}, groupRepo, subRepo, quotaRepo, nil, &config.Config{})

	groups, err := svc.GetAvailableGroups(context.Background(), 7)

	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, groupID, groups[0].ID)
}

var _ APIKeyRepository = (*apiKeyCreateRepoStub)(nil)
var _ UserRepository = (*apiKeyCreateUserRepoStub)(nil)
var _ GroupRepository = (*apiKeyCreateGroupRepoStub)(nil)
var _ UserSubscriptionRepository = (*apiKeyCreateSubscriptionRepoStub)(nil)
var _ QuotaPackageRepository = (*apiKeyCreateQuotaPackageRepoStub)(nil)
var _ UserLegalAgreementRepository = (*apiKeyCreateLegalAgreementRepoStub)(nil)
