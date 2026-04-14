package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrSubSiteNotFound      = infraerrors.NotFound("SUBSITE_NOT_FOUND", "sub-site not found")
	ErrSubSiteSlugExists    = infraerrors.Conflict("SUBSITE_SLUG_EXISTS", "sub-site slug already exists")
	ErrSubSiteDomainExists  = infraerrors.Conflict("SUBSITE_DOMAIN_EXISTS", "sub-site domain already exists")
	ErrSubSiteInvalidSlug   = infraerrors.BadRequest("SUBSITE_INVALID_SLUG", "sub-site slug must contain only lowercase letters, numbers and hyphens")
	ErrSubSiteOwnerNotFound = infraerrors.BadRequest("SUBSITE_OWNER_NOT_FOUND", "sub-site owner user not found")
)

var subSiteSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type SubSiteRepository interface {
	List(ctx context.Context, params pagination.PaginationParams, search, status string) ([]SubSite, *pagination.PaginationResult, error)
	GetByID(ctx context.Context, id int64) (*SubSite, error)
	GetByDomain(ctx context.Context, domain string) (*SubSite, error)
	GetBySlug(ctx context.Context, slug string) (*SubSite, error)
	ExistsBySlug(ctx context.Context, slug string, excludeID int64) (bool, error)
	ExistsByDomain(ctx context.Context, domain string, excludeID int64) (bool, error)
	Create(ctx context.Context, site *SubSite) error
	Update(ctx context.Context, site *SubSite) error
	Delete(ctx context.Context, id int64) error
	BindUser(ctx context.Context, siteID int64, userID int64, source string) error
}

type subSiteCacheEntry struct {
	site      *SubSite
	expiresAt time.Time
}

type SubSiteService struct {
	repo            SubSiteRepository
	userRepo        UserRepository
	mainDomains     map[string]struct{}
	subdomainSuffix string
	cacheTTL        time.Duration
	cacheMu         sync.RWMutex
	hostCache       map[string]subSiteCacheEntry
	onUpdate        func()
}

func NewSubSiteService(repo SubSiteRepository, userRepo UserRepository) *SubSiteService {
	return &SubSiteService{
		repo:            repo,
		userRepo:        userRepo,
		mainDomains:     parseMainDomains(os.Getenv("SUBSITE_MAIN_DOMAINS")),
		subdomainSuffix: normalizeHost(os.Getenv("SUBSITE_SUBDOMAIN_SUFFIX")),
		cacheTTL:        time.Minute,
		hostCache:       make(map[string]subSiteCacheEntry),
	}
}

func (s *SubSiteService) SetOnUpdateCallback(callback func()) {
	s.onUpdate = callback
}

func parseMainDomains(raw string) map[string]struct{} {
	result := map[string]struct{}{
		"localhost":  {},
		"127.0.0.1":  {},
		"0.0.0.0":    {},
		"::1":        {},
		"[::1]":      {},
		"localhost.": {},
		"127.0.0.1.": {},
		"0.0.0.0.":   {},
	}
	for _, item := range strings.Split(raw, ",") {
		host := normalizeHost(item)
		if host != "" {
			result[host] = struct{}{}
		}
	}
	return result
}

func normalizeHost(host string) string {
	host = strings.TrimSpace(strings.ToLower(host))
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimSuffix(host, "/")
	if host == "" {
		return ""
	}
	if strings.HasPrefix(host, "[") && strings.Contains(host, "]") {
		if parsed, _, err := net.SplitHostPort(host); err == nil {
			return strings.Trim(parsed, "[]")
		}
	}
	if parsed, _, err := net.SplitHostPort(host); err == nil {
		return strings.TrimSpace(strings.ToLower(parsed))
	}
	if idx := strings.Index(host, ":"); idx >= 0 && strings.Count(host, ":") == 1 {
		return host[:idx]
	}
	return host
}

func normalizeDomain(domain string) string {
	return normalizeHost(domain)
}

func normalizeSlug(slug string) string {
	slug = strings.ToLower(strings.TrimSpace(slug))
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ReplaceAll(slug, " ", "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	return strings.Trim(slug, "-")
}

func (s *SubSiteService) normalizeAndValidateInput(ownerUserID int64, name, slug, customDomain, status string) (string, string, string, error) {
	name = strings.TrimSpace(name)
	slug = normalizeSlug(slug)
	customDomain = normalizeDomain(customDomain)
	status = strings.TrimSpace(strings.ToLower(status))
	if name == "" || ownerUserID <= 0 {
		return "", "", "", infraerrors.BadRequest("SUBSITE_INVALID_INPUT", "sub-site name and owner_user_id are required")
	}
	if !subSiteSlugPattern.MatchString(slug) {
		return "", "", "", ErrSubSiteInvalidSlug
	}
	if status == "" {
		status = SubSiteStatusActive
	}
	if status != SubSiteStatusActive && status != SubSiteStatusDisabled {
		return "", "", "", infraerrors.BadRequest("SUBSITE_INVALID_STATUS", "sub-site status must be active or disabled")
	}
	return name, slug, customDomain, nil
}

func (s *SubSiteService) List(ctx context.Context, params pagination.PaginationParams, search, status string) ([]SubSite, *pagination.PaginationResult, error) {
	sites, pag, err := s.repo.List(ctx, params, strings.TrimSpace(search), strings.TrimSpace(strings.ToLower(status)))
	if err != nil {
		return nil, nil, err
	}
	for i := range sites {
		sites[i].EntryURL = s.buildEntryURL(&sites[i])
	}
	return sites, pag, nil
}

func (s *SubSiteService) Create(ctx context.Context, input CreateSubSiteInput) (*SubSite, error) {
	name, slug, customDomain, err := s.normalizeAndValidateInput(input.OwnerUserID, input.Name, input.Slug, input.CustomDomain, input.Status)
	if err != nil {
		return nil, err
	}
	if _, err := s.userRepo.GetByID(ctx, input.OwnerUserID); err != nil {
		return nil, ErrSubSiteOwnerNotFound
	}
	if exists, err := s.repo.ExistsBySlug(ctx, slug, 0); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrSubSiteSlugExists
	}
	if customDomain != "" {
		if exists, err := s.repo.ExistsByDomain(ctx, customDomain, 0); err != nil {
			return nil, err
		} else if exists {
			return nil, ErrSubSiteDomainExists
		}
	}
	site := &SubSite{
		OwnerUserID:  input.OwnerUserID,
		Name:         name,
		Slug:         slug,
		CustomDomain: customDomain,
		Status:       normalizeStatus(input.Status),
		SiteLogo:     strings.TrimSpace(input.SiteLogo),
		SiteFavicon:  strings.TrimSpace(input.SiteFavicon),
		SiteSubtitle: strings.TrimSpace(input.SiteSubtitle),
		Announcement: strings.TrimSpace(input.Announcement),
		ContactInfo:  strings.TrimSpace(input.ContactInfo),
		DocURL:       strings.TrimSpace(input.DocURL),
		HomeContent:  strings.TrimSpace(input.HomeContent),
		ThemeConfig:  strings.TrimSpace(input.ThemeConfig),
	}
	if err := s.repo.Create(ctx, site); err != nil {
		return nil, err
	}
	s.invalidateCaches()
	created, err := s.repo.GetByID(ctx, site.ID)
	if err != nil {
		return site, nil
	}
	created.EntryURL = s.buildEntryURL(created)
	return created, nil
}

func (s *SubSiteService) Update(ctx context.Context, input UpdateSubSiteInput) (*SubSite, error) {
	if input.ID <= 0 {
		return nil, ErrSubSiteNotFound
	}
	name, slug, customDomain, err := s.normalizeAndValidateInput(input.OwnerUserID, input.Name, input.Slug, input.CustomDomain, input.Status)
	if err != nil {
		return nil, err
	}
	current, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if _, err := s.userRepo.GetByID(ctx, input.OwnerUserID); err != nil {
		return nil, ErrSubSiteOwnerNotFound
	}
	if exists, err := s.repo.ExistsBySlug(ctx, slug, input.ID); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrSubSiteSlugExists
	}
	if customDomain != "" {
		if exists, err := s.repo.ExistsByDomain(ctx, customDomain, input.ID); err != nil {
			return nil, err
		} else if exists {
			return nil, ErrSubSiteDomainExists
		}
	}
	current.OwnerUserID = input.OwnerUserID
	current.Name = name
	current.Slug = slug
	current.CustomDomain = customDomain
	current.Status = normalizeStatus(input.Status)
	current.SiteLogo = strings.TrimSpace(input.SiteLogo)
	current.SiteFavicon = strings.TrimSpace(input.SiteFavicon)
	current.SiteSubtitle = strings.TrimSpace(input.SiteSubtitle)
	current.Announcement = strings.TrimSpace(input.Announcement)
	current.ContactInfo = strings.TrimSpace(input.ContactInfo)
	current.DocURL = strings.TrimSpace(input.DocURL)
	current.HomeContent = strings.TrimSpace(input.HomeContent)
	current.ThemeConfig = strings.TrimSpace(input.ThemeConfig)
	if err := s.repo.Update(ctx, current); err != nil {
		return nil, err
	}
	s.invalidateCaches()
	updated, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return current, nil
	}
	updated.EntryURL = s.buildEntryURL(updated)
	return updated, nil
}

func (s *SubSiteService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrSubSiteNotFound
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCaches()
	return nil
}

func normalizeStatus(status string) string {
	status = strings.TrimSpace(strings.ToLower(status))
	if status == "" {
		return SubSiteStatusActive
	}
	return status
}

func (s *SubSiteService) invalidateCaches() {
	s.cacheMu.Lock()
	s.hostCache = make(map[string]subSiteCacheEntry)
	s.cacheMu.Unlock()
	if s.onUpdate != nil {
		s.onUpdate()
	}
}

func (s *SubSiteService) ResolveByHost(ctx context.Context, host string) (*SubSite, error) {
	host = normalizeHost(host)
	if host == "" || s.isMainDomain(host) {
		return nil, nil
	}
	if cached := s.getCachedHost(host); cached != nil || s.isCachedMiss(host) {
		return cached, nil
	}
	if site, err := s.repo.GetByDomain(ctx, host); err == nil && site != nil && site.Status == SubSiteStatusActive {
		site.EntryURL = s.buildEntryURL(site)
		s.setCachedHost(host, site)
		return site, nil
	} else if err != nil && !infraerrors.IsNotFound(err) {
		return nil, err
	}
	if slug := s.extractSubdomain(host); slug != "" {
		site, err := s.repo.GetBySlug(ctx, slug)
		if err == nil && site != nil && site.Status == SubSiteStatusActive {
			site.EntryURL = s.buildEntryURL(site)
			s.setCachedHost(host, site)
			return site, nil
		}
		if err != nil && !infraerrors.IsNotFound(err) {
			return nil, err
		}
	}
	s.setCachedHost(host, nil)
	return nil, nil
}

func (s *SubSiteService) getCachedHost(host string) *SubSite {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()
	entry, ok := s.hostCache[host]
	if !ok || time.Now().After(entry.expiresAt) || entry.site == nil {
		return nil
	}
	clone := *entry.site
	return &clone
}

func (s *SubSiteService) isCachedMiss(host string) bool {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()
	entry, ok := s.hostCache[host]
	return ok && entry.site == nil && time.Now().Before(entry.expiresAt)
}

func (s *SubSiteService) setCachedHost(host string, site *SubSite) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	var clone *SubSite
	if site != nil {
		copied := *site
		clone = &copied
	}
	s.hostCache[host] = subSiteCacheEntry{
		site:      clone,
		expiresAt: time.Now().Add(s.cacheTTL),
	}
}

func (s *SubSiteService) isMainDomain(host string) bool {
	host = normalizeHost(host)
	if host == "" {
		return true
	}
	if _, ok := s.mainDomains[host]; ok {
		return true
	}
	if strings.HasSuffix(host, ".up.railway.app") {
		return true
	}
	return false
}

func (s *SubSiteService) extractSubdomain(host string) string {
	host = normalizeHost(host)
	suffix := s.subdomainSuffix
	if suffix == "" {
		return ""
	}
	suffix = "." + strings.TrimPrefix(suffix, ".")
	if !strings.HasSuffix(host, suffix) {
		return ""
	}
	sub := strings.TrimSuffix(host, suffix)
	if sub == "" || strings.Contains(sub, ".") {
		return ""
	}
	if sub == "www" || sub == "api" {
		return ""
	}
	return sub
}

func (s *SubSiteService) buildEntryURL(site *SubSite) string {
	if site == nil {
		return ""
	}
	if site.CustomDomain != "" {
		return fmt.Sprintf("https://%s", site.CustomDomain)
	}
	if s.subdomainSuffix != "" && site.Slug != "" {
		return fmt.Sprintf("https://%s.%s", site.Slug, strings.TrimPrefix(s.subdomainSuffix, "."))
	}
	return ""
}

func (s *SubSiteService) GetCurrent(ctx context.Context) (*SubSite, bool) {
	if ctx == nil {
		return nil, false
	}
	site, ok := ctx.Value(ctxkey.SubSite).(*SubSite)
	if !ok || site == nil {
		return nil, false
	}
	return site, true
}

func (s *SubSiteService) BindCurrentUser(ctx context.Context, userID int64) error {
	site, ok := s.GetCurrent(ctx)
	if !ok || site == nil || userID <= 0 {
		return nil
	}
	return s.repo.BindUser(ctx, site.ID, userID, "register")
}

func (s *SubSiteService) ApplyPublicSettings(ctx context.Context, base *PublicSettings) *PublicSettings {
	if base == nil {
		return nil
	}
	cloned := *base
	site, ok := s.GetCurrent(ctx)
	if !ok || site == nil {
		return &cloned
	}
	cloned.IsSubSite = true
	cloned.SubSiteSlug = site.Slug
	cloned.SubSiteDomain = site.CustomDomain
	if site.Name != "" {
		cloned.SiteName = site.Name
	}
	if site.SiteLogo != "" {
		cloned.SiteLogo = site.SiteLogo
	}
	if site.SiteFavicon != "" {
		cloned.SiteFavicon = site.SiteFavicon
	} else if site.SiteLogo != "" {
		cloned.SiteFavicon = site.SiteLogo
	}
	if site.SiteSubtitle != "" {
		cloned.SiteSubtitle = site.SiteSubtitle
	}
	if site.ContactInfo != "" {
		cloned.ContactInfo = site.ContactInfo
	}
	if site.DocURL != "" {
		cloned.DocURL = site.DocURL
	}
	if site.HomeContent != "" {
		cloned.HomeContent = site.HomeContent
	}
	return &cloned
}

func (s *SubSiteService) CurrentCacheKey(ctx context.Context) string {
	if cacheKey, ok := ctx.Value(ctxkey.SubSiteCacheKey).(string); ok && strings.TrimSpace(cacheKey) != "" {
		return cacheKey
	}
	return "main"
}

func (s *SubSiteService) SiteNameOrEmpty(ctx context.Context) string {
	site, ok := s.GetCurrent(ctx)
	if !ok || site == nil {
		return ""
	}
	return strings.TrimSpace(site.Name)
}
