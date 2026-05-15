package service

import (
	"context"
	"fmt"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrUserNotFound      = infraerrors.NotFound("USER_NOT_FOUND", "user not found")
	ErrPasswordIncorrect = infraerrors.BadRequest("PASSWORD_INCORRECT", "current password is incorrect")
	ErrInsufficientPerms = infraerrors.Forbidden("INSUFFICIENT_PERMISSIONS", "insufficient permissions")
)

// UserListFilters contains all filter options for listing users
type UserListFilters struct {
	Status        string           // User status filter
	Role          string           // User role filter
	Search        string           // Search in email, username
	Attributes    map[int64]string // Custom attribute filters: attributeID -> value
	CreatedAfter  string           // Filter users created on or after this date (YYYY-MM-DD)
	CreatedBefore string           // Filter users created on or before this date (YYYY-MM-DD)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetFirstAdmin(ctx context.Context) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error

	List(ctx context.Context, params pagination.PaginationParams) ([]User, *pagination.PaginationResult, error)
	ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters UserListFilters) ([]User, *pagination.PaginationResult, error)

	UpdateBalance(ctx context.Context, id int64, amount float64) error
	DeductBalance(ctx context.Context, id int64, amount float64) error
	GetBalance(ctx context.Context, id int64) (float64, error)
	UpdateConcurrency(ctx context.Context, id int64, amount int) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	RemoveGroupFromAllowedGroups(ctx context.Context, groupID int64) (int64, error)

	// TOTP 相关方法
	UpdateTotpSecret(ctx context.Context, userID int64, encryptedSecret *string) error
	EnableTotp(ctx context.Context, userID int64) error
	DisableTotp(ctx context.Context, userID int64) error

	// 邀请码
	GetByInviteCode(ctx context.Context, code string) (*User, error)

	// 来源统计
	GetDiscoverySourceStats(ctx context.Context, startTime, endTime time.Time) ([]DiscoverySourceStat, error)

	// 初始余额过期清理
	ClearExpiredInitialBalances(ctx context.Context) ([]int64, error)
}

type DiscoverySourceStat struct {
	Source string
	Count  int
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Email       *string `json:"email"`
	Username    *string `json:"username"`
	Concurrency *int    `json:"concurrency"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// UserService 用户服务
type UserService struct {
	userRepo             UserRepository
	legalRepo            UserLegalAgreementRepository
	authCacheInvalidator APIKeyAuthCacheInvalidator
	userCache            UserCache
}

type UserCache interface {
	Get(ctx context.Context, userID int64) (*User, error)
	Set(ctx context.Context, userID int64, user *User) error
	Delete(ctx context.Context, userID int64) error
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo UserRepository, authCacheInvalidator APIKeyAuthCacheInvalidator, userCache UserCache) *UserService {
	return &UserService{
		userRepo:             userRepo,
		authCacheInvalidator: authCacheInvalidator,
		userCache:            userCache,
	}
}

func (s *UserService) SetLegalAgreementRepository(repo UserLegalAgreementRepository) {
	s.legalRepo = repo
}

func (s *UserService) hydrateLegalAgreement(ctx context.Context, user *User) error {
	if s.legalRepo == nil || user == nil {
		return nil
	}
	agreement, err := s.legalRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		if IsLegalAgreementMissing(err) {
			user.LegalAgreementAccepted = false
			user.LegalAgreement = nil
			return nil
		}
		return fmt.Errorf("get legal agreement: %w", err)
	}
	user.LegalAgreement = agreement
	user.LegalAgreementAccepted = agreement.IsCurrent()
	return nil
}

func (s *UserService) AcceptLegalAgreement(ctx context.Context, userID int64) error {
	if s.legalRepo == nil {
		return ErrServiceUnavailable
	}
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if err := s.legalRepo.Upsert(ctx, NewCurrentLegalAgreement(userID)); err != nil {
		return err
	}
	if s.authCacheInvalidator != nil {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, userID)
	}
	if s.userCache != nil {
		_ = s.userCache.Delete(ctx, userID)
	}
	return nil
}

// GetFirstAdmin 获取首个管理员用户（用于 Admin API Key 认证）
func (s *UserService) GetFirstAdmin(ctx context.Context) (*User, error) {
	admin, err := s.userRepo.GetFirstAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("get first admin: %w", err)
	}
	return admin, nil
}

// GetProfile 获取用户资料
func (s *UserService) GetProfile(ctx context.Context, userID int64) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if err := s.hydrateLegalAgreement(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(ctx context.Context, userID int64, req UpdateProfileRequest) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	oldConcurrency := user.Concurrency

	// 更新字段
	if req.Email != nil {
		// 检查新邮箱是否已被使用
		exists, err := s.userRepo.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("check email exists: %w", err)
		}
		if exists && *req.Email != user.Email {
			return nil, ErrEmailExists
		}
		user.Email = *req.Email
	}

	if req.Username != nil {
		user.Username = *req.Username
	}

	if req.Concurrency != nil {
		user.Concurrency = *req.Concurrency
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	if s.authCacheInvalidator != nil && user.Concurrency != oldConcurrency {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, userID)
	}
	if s.userCache != nil {
		_ = s.userCache.Delete(ctx, userID)
	}

	return user, nil
}

// UpdateDiscoverySource 更新用户来源渠道
func (s *UserService) UpdateDiscoverySource(ctx context.Context, userID int64, source string) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	user.DiscoverySource = &source
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

// ChangePassword 修改密码
// Security: Increments TokenVersion to invalidate all existing JWT tokens
func (s *UserService) ChangePassword(ctx context.Context, userID int64, req ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	// 验证当前密码
	if !user.CheckPassword(req.CurrentPassword) {
		return ErrPasswordIncorrect
	}

	if err := user.SetPassword(req.NewPassword); err != nil {
		return fmt.Errorf("set password: %w", err)
	}

	// Increment TokenVersion to invalidate all existing tokens
	// This ensures that any tokens issued before the password change become invalid
	user.TokenVersion++

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

// GetByID 根据ID获取用户（管理员功能）
func (s *UserService) GetByID(ctx context.Context, id int64) (*User, error) {
	// Try cache first
	if s.userCache != nil {
		if user, err := s.userCache.Get(ctx, id); err == nil {
			if err := s.hydrateLegalAgreement(ctx, user); err != nil {
				return nil, err
			}
			return user, nil
		}
	}

	// Cache miss, query database
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if err := s.hydrateLegalAgreement(ctx, user); err != nil {
		return nil, err
	}

	// Update cache
	if s.userCache != nil {
		_ = s.userCache.Set(ctx, id, user)
	}

	return user, nil
}

// List 获取用户列表（管理员功能）
func (s *UserService) List(ctx context.Context, params pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) {
	users, pagination, err := s.userRepo.List(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("list users: %w", err)
	}
	return users, pagination, nil
}

// UpdateBalance 更新用户余额（管理员功能）
func (s *UserService) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	if err := s.userRepo.UpdateBalance(ctx, userID, amount); err != nil {
		return fmt.Errorf("update balance: %w", err)
	}
	if s.authCacheInvalidator != nil {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, userID)
	}
	if s.userCache != nil {
		_ = s.userCache.Delete(ctx, userID)
	}
	return nil
}

// UpdateConcurrency 更新用户并发数（管理员功能）
func (s *UserService) UpdateConcurrency(ctx context.Context, userID int64, concurrency int) error {
	if err := s.userRepo.UpdateConcurrency(ctx, userID, concurrency); err != nil {
		return fmt.Errorf("update concurrency: %w", err)
	}
	if s.authCacheInvalidator != nil {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, userID)
	}
	return nil
}

// UpdateStatus 更新用户状态（管理员功能）
func (s *UserService) UpdateStatus(ctx context.Context, userID int64, status string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	user.Status = status

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	if s.authCacheInvalidator != nil {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, userID)
	}
	if s.userCache != nil {
		_ = s.userCache.Delete(ctx, userID)
	}

	return nil
}

// Delete 删除用户（管理员功能）
func (s *UserService) Delete(ctx context.Context, userID int64) error {
	if s.authCacheInvalidator != nil {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, userID)
	}
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// GetDiscoverySourceStats 获取来源统计
func (s *UserService) GetDiscoverySourceStats(ctx context.Context, startTime, endTime time.Time) ([]DiscoverySourceStat, error) {
	return s.userRepo.GetDiscoverySourceStats(ctx, startTime, endTime)
}
