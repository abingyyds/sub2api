package service

import (
	"context"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrAnnouncementNotFound = infraerrors.NotFound("ANNOUNCEMENT_NOT_FOUND", "announcement not found")
)

// AnnouncementRepository defines the data access interface for announcements.
type AnnouncementRepository interface {
	List(ctx context.Context, params pagination.PaginationParams) ([]Announcement, *pagination.PaginationResult, error)
	GetByID(ctx context.Context, id int64) (*Announcement, error)
	Create(ctx context.Context, a *Announcement) error
	Update(ctx context.Context, a *Announcement) error
	Delete(ctx context.Context, id int64) error
	ListActive(ctx context.Context) ([]Announcement, error)
}

// AnnouncementService handles announcement business logic.
type AnnouncementService struct {
	repo AnnouncementRepository
}

// NewAnnouncementService creates a new AnnouncementService.
func NewAnnouncementService(repo AnnouncementRepository) *AnnouncementService {
	return &AnnouncementService{repo: repo}
}

func (s *AnnouncementService) List(ctx context.Context, params pagination.PaginationParams) ([]Announcement, *pagination.PaginationResult, error) {
	return s.repo.List(ctx, params)
}

func (s *AnnouncementService) GetByID(ctx context.Context, id int64) (*Announcement, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AnnouncementService) Create(ctx context.Context, a *Announcement) error {
	return s.repo.Create(ctx, a)
}

func (s *AnnouncementService) Update(ctx context.Context, a *Announcement) error {
	return s.repo.Update(ctx, a)
}

func (s *AnnouncementService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *AnnouncementService) ListActive(ctx context.Context) ([]Announcement, error) {
	return s.repo.ListActive(ctx)
}
