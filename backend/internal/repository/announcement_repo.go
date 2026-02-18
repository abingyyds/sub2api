package repository

import (
	"context"
	"database/sql"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type announcementRepository struct {
	db *sql.DB
}

func NewAnnouncementRepository(sqlDB *sql.DB) service.AnnouncementRepository {
	return &announcementRepository{db: sqlDB}
}

func (r *announcementRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.Announcement, *pagination.PaginationResult, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM announcements`).Scan(&total); err != nil {
		return nil, nil, err
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, content, status, priority, created_at, updated_at FROM announcements ORDER BY priority DESC, id DESC LIMIT $1 OFFSET $2`,
		params.Limit(), params.Offset())
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var results []service.Announcement
	for rows.Next() {
		var a service.Announcement
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Status, &a.Priority, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, nil, err
		}
		results = append(results, a)
	}
	return results, paginationResultFromTotal(total, params), rows.Err()
}

func (r *announcementRepository) GetByID(ctx context.Context, id int64) (*service.Announcement, error) {
	a := &service.Announcement{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, content, status, priority, created_at, updated_at FROM announcements WHERE id = $1`, id).
		Scan(&a.ID, &a.Title, &a.Content, &a.Status, &a.Priority, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, service.ErrAnnouncementNotFound
	}
	return a, err
}

func (r *announcementRepository) Create(ctx context.Context, a *service.Announcement) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO announcements (title, content, status, priority, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at`,
		a.Title, a.Content, a.Status, a.Priority).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *announcementRepository) Update(ctx context.Context, a *service.Announcement) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE announcements SET title=$1, content=$2, status=$3, priority=$4, updated_at=NOW() WHERE id=$5`,
		a.Title, a.Content, a.Status, a.Priority, a.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return service.ErrAnnouncementNotFound
	}
	return nil
}

func (r *announcementRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM announcements WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return service.ErrAnnouncementNotFound
	}
	return nil
}

func (r *announcementRepository) ListActive(ctx context.Context) ([]service.Announcement, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, content, status, priority, created_at, updated_at FROM announcements WHERE status = 'active' ORDER BY priority DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []service.Announcement
	for rows.Next() {
		var a service.Announcement
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Status, &a.Priority, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, a)
	}
	return results, rows.Err()
}
