package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrInviteCodeNotFound = infraerrors.NotFound("INVITE_CODE_NOT_FOUND", "invite code not found")
	ErrInviteCodeExists   = infraerrors.Conflict("INVITE_CODE_EXISTS", "invite code already exists")
	ErrInviteCodeDisabled = infraerrors.BadRequest("INVITE_CODE_DISABLED", "invite code is disabled")
	ErrInviteCodeMaxUses  = infraerrors.BadRequest("INVITE_CODE_MAX_USES", "invite code has reached max uses")
)

type AdminInviteCode struct {
	ID         int64
	Code       string
	SourceName string
	CreatedBy  int64
	UsedCount  int
	MaxUses    *int
	Enabled    bool
	Notes      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type AdminInviteCodeRepository interface {
	Create(ctx context.Context, code *AdminInviteCode) error
	GetByID(ctx context.Context, id int64) (*AdminInviteCode, error)
	GetByCode(ctx context.Context, code string) (*AdminInviteCode, error)
	Update(ctx context.Context, code *AdminInviteCode) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, params pagination.PaginationParams) ([]AdminInviteCode, *pagination.PaginationResult, error)
	IncrementUsedCount(ctx context.Context, id int64) error
}

type AdminInviteCodeService struct {
	repo AdminInviteCodeRepository
}

func NewAdminInviteCodeService(repo AdminInviteCodeRepository) *AdminInviteCodeService {
	return &AdminInviteCodeService{repo: repo}
}

func (s *AdminInviteCodeService) Create(ctx context.Context, sourceName string, createdBy int64, maxUses *int, notes string) (*AdminInviteCode, error) {
	code := generateInviteCode()
	inviteCode := &AdminInviteCode{
		Code:       code,
		SourceName: sourceName,
		CreatedBy:  createdBy,
		MaxUses:    maxUses,
		Enabled:    true,
		Notes:      notes,
	}
	if err := s.repo.Create(ctx, inviteCode); err != nil {
		return nil, fmt.Errorf("create invite code: %w", err)
	}
	return inviteCode, nil
}

func (s *AdminInviteCodeService) GetByCode(ctx context.Context, code string) (*AdminInviteCode, error) {
	return s.repo.GetByCode(ctx, code)
}

func (s *AdminInviteCodeService) ValidateAndUse(ctx context.Context, code string) (string, error) {
	inviteCode, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return "", err
	}
	if !inviteCode.Enabled {
		return "", ErrInviteCodeDisabled
	}
	if inviteCode.MaxUses != nil && inviteCode.UsedCount >= *inviteCode.MaxUses {
		return "", ErrInviteCodeMaxUses
	}
	if err := s.repo.IncrementUsedCount(ctx, inviteCode.ID); err != nil {
		return "", fmt.Errorf("increment used count: %w", err)
	}
	return inviteCode.SourceName, nil
}

func (s *AdminInviteCodeService) List(ctx context.Context, params pagination.PaginationParams) ([]AdminInviteCode, *pagination.PaginationResult, error) {
	return s.repo.List(ctx, params)
}

func (s *AdminInviteCodeService) Update(ctx context.Context, id int64, sourceName *string, maxUses *int, enabled *bool, notes *string) (*AdminInviteCode, error) {
	code, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sourceName != nil {
		code.SourceName = *sourceName
	}
	if maxUses != nil {
		code.MaxUses = maxUses
	}
	if enabled != nil {
		code.Enabled = *enabled
	}
	if notes != nil {
		code.Notes = *notes
	}
	if err := s.repo.Update(ctx, code); err != nil {
		return nil, fmt.Errorf("update invite code: %w", err)
	}
	return code, nil
}

func (s *AdminInviteCodeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func generateInviteCode() string {
	b := make([]byte, 12)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:16]
}
