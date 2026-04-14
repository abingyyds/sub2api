package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type subSiteRepository struct {
	db *sql.DB
}

func NewSubSiteRepository(sqlDB *sql.DB) service.SubSiteRepository {
	return &subSiteRepository{db: sqlDB}
}

func (r *subSiteRepository) List(ctx context.Context, params pagination.PaginationParams, search, status string) ([]service.SubSite, *pagination.PaginationResult, error) {
	var (
		conds []string
		args  []any
	)
	conds = append(conds, "1=1")
	if search = strings.TrimSpace(search); search != "" {
		args = append(args, "%"+search+"%")
		idx := len(args)
		conds = append(conds, fmt.Sprintf("(s.name ILIKE $%d OR s.slug ILIKE $%d OR COALESCE(s.custom_domain, '') ILIKE $%d OR COALESCE(u.email, '') ILIKE $%d)", idx, idx, idx, idx))
	}
	if status = strings.TrimSpace(strings.ToLower(status)); status != "" {
		args = append(args, status)
		conds = append(conds, fmt.Sprintf("s.status = $%d", len(args)))
	}
	where := strings.Join(conds, " AND ")

	countQuery := `SELECT COUNT(*) FROM sub_sites s LEFT JOIN users u ON u.id = s.owner_user_id WHERE ` + where
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, nil, err
	}

	selectArgs := append([]any{}, args...)
	selectArgs = append(selectArgs, params.Limit(), params.Offset())
	query := fmt.Sprintf(`
		SELECT s.id, s.owner_user_id, COALESCE(u.email, ''), s.name, s.slug, COALESCE(s.custom_domain, ''), s.status,
			COALESCE(s.site_logo, ''), COALESCE(s.site_favicon, ''), COALESCE(s.site_subtitle, ''), COALESCE(s.announcement, ''),
			COALESCE(s.contact_info, ''), COALESCE(s.doc_url, ''), COALESCE(s.home_content, ''), COALESCE(s.theme_config, ''),
			s.created_at, s.updated_at,
			COALESCE((SELECT COUNT(*) FROM sub_site_users su WHERE su.sub_site_id = s.id), 0) AS user_count
		FROM sub_sites s
		LEFT JOIN users u ON u.id = s.owner_user_id
		WHERE %s
		ORDER BY s.id DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)
	rows, err := r.db.QueryContext(ctx, query, selectArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	var items []service.SubSite
	for rows.Next() {
		var site service.SubSite
		if err := rows.Scan(
			&site.ID, &site.OwnerUserID, &site.OwnerEmail, &site.Name, &site.Slug, &site.CustomDomain, &site.Status,
			&site.SiteLogo, &site.SiteFavicon, &site.SiteSubtitle, &site.Announcement,
			&site.ContactInfo, &site.DocURL, &site.HomeContent, &site.ThemeConfig,
			&site.CreatedAt, &site.UpdatedAt, &site.UserCount,
		); err != nil {
			return nil, nil, err
		}
		items = append(items, site)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return items, paginationResultFromTotal(total, params), nil
}

func (r *subSiteRepository) GetByID(ctx context.Context, id int64) (*service.SubSite, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT s.id, s.owner_user_id, COALESCE(u.email, ''), s.name, s.slug, COALESCE(s.custom_domain, ''), s.status,
			COALESCE(s.site_logo, ''), COALESCE(s.site_favicon, ''), COALESCE(s.site_subtitle, ''), COALESCE(s.announcement, ''),
			COALESCE(s.contact_info, ''), COALESCE(s.doc_url, ''), COALESCE(s.home_content, ''), COALESCE(s.theme_config, ''),
			s.created_at, s.updated_at,
			COALESCE((SELECT COUNT(*) FROM sub_site_users su WHERE su.sub_site_id = s.id), 0) AS user_count
		FROM sub_sites s
		LEFT JOIN users u ON u.id = s.owner_user_id
		WHERE s.id = $1
	`, id)
	return scanSubSite(row)
}

func (r *subSiteRepository) GetByDomain(ctx context.Context, domain string) (*service.SubSite, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT s.id, s.owner_user_id, COALESCE(u.email, ''), s.name, s.slug, COALESCE(s.custom_domain, ''), s.status,
			COALESCE(s.site_logo, ''), COALESCE(s.site_favicon, ''), COALESCE(s.site_subtitle, ''), COALESCE(s.announcement, ''),
			COALESCE(s.contact_info, ''), COALESCE(s.doc_url, ''), COALESCE(s.home_content, ''), COALESCE(s.theme_config, ''),
			s.created_at, s.updated_at,
			COALESCE((SELECT COUNT(*) FROM sub_site_users su WHERE su.sub_site_id = s.id), 0) AS user_count
		FROM sub_sites s
		LEFT JOIN users u ON u.id = s.owner_user_id
		WHERE LOWER(COALESCE(s.custom_domain, '')) = LOWER($1)
	`, domain)
	return scanSubSite(row)
}

func (r *subSiteRepository) GetBySlug(ctx context.Context, slug string) (*service.SubSite, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT s.id, s.owner_user_id, COALESCE(u.email, ''), s.name, s.slug, COALESCE(s.custom_domain, ''), s.status,
			COALESCE(s.site_logo, ''), COALESCE(s.site_favicon, ''), COALESCE(s.site_subtitle, ''), COALESCE(s.announcement, ''),
			COALESCE(s.contact_info, ''), COALESCE(s.doc_url, ''), COALESCE(s.home_content, ''), COALESCE(s.theme_config, ''),
			s.created_at, s.updated_at,
			COALESCE((SELECT COUNT(*) FROM sub_site_users su WHERE su.sub_site_id = s.id), 0) AS user_count
		FROM sub_sites s
		LEFT JOIN users u ON u.id = s.owner_user_id
		WHERE LOWER(s.slug) = LOWER($1)
	`, slug)
	return scanSubSite(row)
}

func (r *subSiteRepository) ExistsBySlug(ctx context.Context, slug string, excludeID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sub_sites WHERE LOWER(slug) = LOWER($1)`
	args := []any{slug}
	if excludeID > 0 {
		query += ` AND id <> $2`
		args = append(args, excludeID)
	}
	query += `)`
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *subSiteRepository) ExistsByDomain(ctx context.Context, domain string, excludeID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sub_sites WHERE LOWER(COALESCE(custom_domain, '')) = LOWER($1) AND COALESCE(custom_domain, '') <> ''`
	args := []any{domain}
	if excludeID > 0 {
		query += ` AND id <> $2`
		args = append(args, excludeID)
	}
	query += `)`
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *subSiteRepository) Create(ctx context.Context, site *service.SubSite) error {
	query := `
		INSERT INTO sub_sites (
			owner_user_id, name, slug, custom_domain, status,
			site_logo, site_favicon, site_subtitle, announcement,
			contact_info, doc_url, home_content, theme_config,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, NULLIF($4, ''), $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13,
			NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		site.OwnerUserID, site.Name, site.Slug, site.CustomDomain, site.Status,
		site.SiteLogo, site.SiteFavicon, site.SiteSubtitle, site.Announcement,
		site.ContactInfo, site.DocURL, site.HomeContent, site.ThemeConfig,
	).Scan(&site.ID, &site.CreatedAt, &site.UpdatedAt)
}

func (r *subSiteRepository) Update(ctx context.Context, site *service.SubSite) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE sub_sites
		SET owner_user_id = $2,
			name = $3,
			slug = $4,
			custom_domain = NULLIF($5, ''),
			status = $6,
			site_logo = $7,
			site_favicon = $8,
			site_subtitle = $9,
			announcement = $10,
			contact_info = $11,
			doc_url = $12,
			home_content = $13,
			theme_config = $14,
			updated_at = NOW()
		WHERE id = $1
	`,
		site.ID, site.OwnerUserID, site.Name, site.Slug, site.CustomDomain, site.Status,
		site.SiteLogo, site.SiteFavicon, site.SiteSubtitle, site.Announcement,
		site.ContactInfo, site.DocURL, site.HomeContent, site.ThemeConfig,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return service.ErrSubSiteNotFound
	}
	return nil
}

func (r *subSiteRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM sub_sites WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return service.ErrSubSiteNotFound
	}
	return nil
}

func (r *subSiteRepository) BindUser(ctx context.Context, siteID int64, userID int64, source string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO sub_site_users (sub_site_id, user_id, source, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET sub_site_id = EXCLUDED.sub_site_id, source = EXCLUDED.source, updated_at = NOW()
	`, siteID, userID, source)
	return err
}

func scanSubSite(row *sql.Row) (*service.SubSite, error) {
	var site service.SubSite
	if err := row.Scan(
		&site.ID, &site.OwnerUserID, &site.OwnerEmail, &site.Name, &site.Slug, &site.CustomDomain, &site.Status,
		&site.SiteLogo, &site.SiteFavicon, &site.SiteSubtitle, &site.Announcement,
		&site.ContactInfo, &site.DocURL, &site.HomeContent, &site.ThemeConfig,
		&site.CreatedAt, &site.UpdatedAt, &site.UserCount,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrSubSiteNotFound
		}
		return nil, err
	}
	return &site, nil
}
