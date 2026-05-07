package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrSubSiteNotFound             = infraerrors.NotFound("SUBSITE_NOT_FOUND", "sub-site not found")
	ErrSubSiteSlugExists           = infraerrors.Conflict("SUBSITE_SLUG_EXISTS", "sub-site slug already exists")
	ErrSubSiteDomainExists         = infraerrors.Conflict("SUBSITE_DOMAIN_EXISTS", "sub-site domain already exists")
	ErrSubSiteInvalidSlug          = infraerrors.BadRequest("SUBSITE_INVALID_SLUG", "sub-site slug must contain only lowercase letters, numbers and hyphens")
	ErrSubSiteOwnerNotFound        = infraerrors.BadRequest("SUBSITE_OWNER_NOT_FOUND", "sub-site owner user not found")
	ErrSubSiteInvalidStatus        = infraerrors.BadRequest("SUBSITE_INVALID_STATUS", "sub-site status must be pending, active or disabled")
	ErrSubSiteInvalidThemeTemplate = infraerrors.BadRequest("SUBSITE_INVALID_THEME_TEMPLATE", "sub-site theme template is invalid")
	ErrSubSiteInvalidRegistration  = infraerrors.BadRequest("SUBSITE_INVALID_REGISTRATION_MODE", "sub-site registration mode must be open, invite or closed")
	ErrSubSiteParentNotFound       = infraerrors.BadRequest("SUBSITE_PARENT_NOT_FOUND", "parent sub-site not found")
	ErrSubSiteLevelExceeded        = infraerrors.BadRequest("SUBSITE_LEVEL_EXCEEDED", "sub-site level exceeds the supported hierarchy")
	ErrSubSiteRegistrationClosed   = infraerrors.Forbidden("SUBSITE_REGISTRATION_CLOSED", "registration is closed for this sub-site")
	ErrSubSiteInviteRequired       = infraerrors.BadRequest("SUBSITE_INVITE_REQUIRED", "this sub-site requires an invite code for registration")
	ErrSubSiteOpenDisabled         = infraerrors.Forbidden("SUBSITE_OPEN_DISABLED", "sub-site self-service opening is disabled")
	ErrSubSiteActivationNotFound   = infraerrors.NotFound("SUBSITE_ACTIVATION_NOT_FOUND", "sub-site activation request not found")
	ErrSubSiteForbidden            = infraerrors.Forbidden("SUBSITE_FORBIDDEN", "you do not have permission to manage this sub-site")
	ErrSubSiteUserScopeMismatch    = infraerrors.Forbidden("SUBSITE_USER_SCOPE_MISMATCH", "your account is bound to a different site")
	ErrSubSitePoolInsufficient     = infraerrors.BadRequest("SUBSITE_POOL_INSUFFICIENT", "sub-site pool balance is insufficient")
)

var subSiteSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type SubSiteRepository interface {
	List(ctx context.Context, params pagination.PaginationParams, search, status string) ([]SubSite, *pagination.PaginationResult, error)
	ListByOwner(ctx context.Context, ownerUserID int64) ([]SubSite, error)
	GetByID(ctx context.Context, id int64) (*SubSite, error)
	GetByDomain(ctx context.Context, domain string) (*SubSite, error)
	GetBySlug(ctx context.Context, slug string) (*SubSite, error)
	GetBoundSubSiteByUserID(ctx context.Context, userID int64) (*SubSite, error)
	ExistsBySlug(ctx context.Context, slug string, excludeID int64) (bool, error)
	ExistsByDomain(ctx context.Context, domain string, excludeID int64) (bool, error)
	Create(ctx context.Context, site *SubSite) error
	Update(ctx context.Context, site *SubSite) error
	UpdateMode(ctx context.Context, siteID int64, newMode string) error
	SubmitHomeContentReview(ctx context.Context, siteID int64, pendingContent string) error
	ReviewHomeContent(ctx context.Context, siteID int64, approved bool, reviewerID int64, reviewNote string) error
	IncrementTotalWithdrawnFen(ctx context.Context, siteID int64, amountFen int64) error
	Delete(ctx context.Context, id int64) error
	BindUser(ctx context.Context, siteID int64, userID int64, source string) error
	UnbindUser(ctx context.Context, userID int64) error
	// 级联停用：把 sub_sites 中 parent 链包含 rootID 的所有后代（递归）status 置为 newStatus。
	// 返回受影响的分站 id 列表。
	CascadeUpdateStatus(ctx context.Context, rootID int64, newStatus string) ([]int64, error)
	// 分站余额池：原子增减，返回变动后余额（正数入账、负数出账；允许负余额：透支在业务层记 warn）
	AdjustBalance(ctx context.Context, siteID int64, deltaFen int64, entry SubSiteLedgerEntry) (int64, error)
	// ApplyUserBalanceAndPoolLedger 在同一数据库事务内更新用户余额并写入一组分站池流水。
	// 用于线下加余额 / 自动进货这类必须账实一致的路径；所有负数池变动都禁止透支。
	ApplyUserBalanceAndPoolLedger(ctx context.Context, userID int64, balanceDelta float64, entries []SubSiteLedgerEntry) error
	ListLedger(ctx context.Context, siteID int64, params pagination.PaginationParams, txType string) ([]SubSiteLedgerEntry, *pagination.PaginationResult, error)
	CreateActivationRequest(ctx context.Context, request *SubSiteActivationRequest) error
	GetActivationRequestByOrderID(ctx context.Context, orderID int64) (*SubSiteActivationRequest, error)
	MarkActivationRequestCompleted(ctx context.Context, orderID int64, subSiteID int64) error
}

type subSiteCacheEntry struct {
	site      *SubSite
	expiresAt time.Time
}

type SubSiteService struct {
	repo            SubSiteRepository
	userRepo        UserRepository
	settingRepo     SettingRepository
	mainDomains     map[string]struct{}
	subdomainSuffix string
	cacheTTL        time.Duration
	cacheMu         sync.RWMutex
	hostCache       map[string]subSiteCacheEntry
	boundCacheMu    sync.RWMutex
	boundCache      map[int64]boundSubSiteEntry
	onUpdate        func()
}

func NewSubSiteService(repo SubSiteRepository, userRepo UserRepository, settingRepo SettingRepository) *SubSiteService {
	return &SubSiteService{
		repo:            repo,
		userRepo:        userRepo,
		settingRepo:     settingRepo,
		mainDomains:     parseMainDomains(os.Getenv("SUBSITE_MAIN_DOMAINS")),
		subdomainSuffix: normalizeHost(os.Getenv("SUBSITE_SUBDOMAIN_SUFFIX")),
		cacheTTL:        time.Minute,
		hostCache:       make(map[string]subSiteCacheEntry),
		boundCache:      make(map[int64]boundSubSiteEntry),
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

func normalizeStatus(status string) string {
	status = strings.TrimSpace(strings.ToLower(status))
	if status == "" {
		return SubSiteStatusActive
	}
	return status
}

func normalizeThemeTemplate(template string) string {
	template = strings.TrimSpace(strings.ToLower(template))
	if template == "" {
		return SubSiteThemeTemplateStarter
	}
	return template
}

func normalizeRegistrationMode(mode string) string {
	mode = strings.TrimSpace(strings.ToLower(mode))
	if mode == "" {
		return SubSiteRegistrationOpen
	}
	return mode
}

func normalizeConsumeRateMultiplier(multiplier float64) float64 {
	if multiplier <= 0 {
		return DefaultSubSiteConsumeRate
	}
	return multiplier
}

func isValidThemeTemplate(template string) bool {
	for _, item := range DefaultSubSiteThemeTemplates {
		if item.Key == template {
			return true
		}
	}
	return false
}

func boolValue(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

// resolveParent 解析父分站并计算新层级。
// excludeID 用于 Update 场景：防止把自己挂到自己或自己的后代之下形成环。
// 通过沿 parent 链向上遍历（以 MaxSubSiteLevelHardLimit 作为安全上界），若遇到 excludeID 即判定为环。
func (s *SubSiteService) resolveParent(ctx context.Context, parentID *int64, excludeID int64) (*SubSite, int, error) {
	if parentID == nil || *parentID <= 0 {
		return nil, 1, nil
	}
	if excludeID > 0 && *parentID == excludeID {
		return nil, 0, ErrSubSiteLevelExceeded
	}
	parent, err := s.repo.GetByID(ctx, *parentID)
	if err != nil {
		return nil, 0, ErrSubSiteParentNotFound
	}
	// 沿 parent 链向上遍历验证无环
	cursor := parent
	for hop := 0; hop < MaxSubSiteLevelHardLimit+1 && cursor != nil; hop++ {
		if excludeID > 0 && cursor.ID == excludeID {
			return nil, 0, ErrSubSiteLevelExceeded
		}
		if cursor.ParentSubSiteID == nil || *cursor.ParentSubSiteID <= 0 {
			break
		}
		next, err := s.repo.GetByID(ctx, *cursor.ParentSubSiteID)
		if err != nil {
			return nil, 0, ErrSubSiteParentNotFound
		}
		cursor = next
	}
	maxLevel := s.MaxLevel(ctx)
	level := parent.Level + 1
	if level > maxLevel {
		return nil, 0, ErrSubSiteLevelExceeded
	}
	return parent, level, nil
}

func (s *SubSiteService) normalizeAndValidateInput(ownerUserID int64, name, slug, customDomain, status, themeTemplate, registrationMode string) (string, string, string, string, string, error) {
	name = strings.TrimSpace(name)
	slug = normalizeSlug(slug)
	customDomain = normalizeDomain(customDomain)
	status = normalizeStatus(status)
	themeTemplate = normalizeThemeTemplate(themeTemplate)
	registrationMode = normalizeRegistrationMode(registrationMode)

	if name == "" || ownerUserID <= 0 {
		return "", "", "", "", "", infraerrors.BadRequest("SUBSITE_INVALID_INPUT", "sub-site name and owner_user_id are required")
	}
	if !subSiteSlugPattern.MatchString(slug) {
		return "", "", "", "", "", ErrSubSiteInvalidSlug
	}
	if status != SubSiteStatusPending && status != SubSiteStatusActive && status != SubSiteStatusDisabled {
		return "", "", "", "", "", ErrSubSiteInvalidStatus
	}
	if !isValidThemeTemplate(themeTemplate) {
		return "", "", "", "", "", ErrSubSiteInvalidThemeTemplate
	}
	if registrationMode != SubSiteRegistrationOpen && registrationMode != SubSiteRegistrationInvite && registrationMode != SubSiteRegistrationClosed {
		return "", "", "", "", "", ErrSubSiteInvalidRegistration
	}
	return name, slug, customDomain, themeTemplate, registrationMode, nil
}

func (s *SubSiteService) applyInputToSite(site *SubSite, input CreateSubSiteInput, level int, status string, themeTemplate string, registrationMode string) {
	site.OwnerUserID = input.OwnerUserID
	site.ParentSubSiteID = input.ParentSubSiteID
	site.Level = level
	site.Name = strings.TrimSpace(input.Name)
	site.Slug = normalizeSlug(input.Slug)
	site.CustomDomain = normalizeDomain(input.CustomDomain)
	site.Status = status
	site.Mode = normalizeMode(input.Mode)
	site.SiteLogo = strings.TrimSpace(input.SiteLogo)
	site.SiteFavicon = strings.TrimSpace(input.SiteFavicon)
	site.SiteSubtitle = strings.TrimSpace(input.SiteSubtitle)
	site.Announcement = strings.TrimSpace(input.Announcement)
	site.ContactInfo = strings.TrimSpace(input.ContactInfo)
	site.DocURL = strings.TrimSpace(input.DocURL)
	site.HomeContent = strings.TrimSpace(input.HomeContent)
	site.ThemeTemplate = themeTemplate
	site.RegistrationMode = registrationMode
	site.EnableTopup = boolValue(input.EnableTopup, true)
	site.AllowSubSite = boolValue(input.AllowSubSite, false)
	site.SubSitePriceFen = input.SubSitePriceFen
	site.ConsumeRateMultiplier = normalizeConsumeRateMultiplier(input.ConsumeRateMultiplier)
	site.AllowOnlineTopup = boolValue(input.AllowOnlineTopup, true)
	site.AllowOfflineTopup = boolValue(input.AllowOfflineTopup, true)
	site.OwnerPaymentConfig = input.OwnerPaymentConfig
	site.SubscriptionExpiredAt = input.SubscriptionExpiredAt
}

// normalizeMode 把任意模式值对齐到枚举；未知或空回退 pool。
// rate 模式要求 owner 显式选择（admin 端），pool 是默认。
func normalizeMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case SubSiteModeRate:
		return SubSiteModeRate
	default:
		return SubSiteModePool
	}
}

func (s *SubSiteService) populateComputedFields(ctx context.Context, site *SubSite, includePricing bool) (*SubSite, error) {
	if site == nil {
		return nil, nil
	}
	site.EntryURL = s.buildEntryURL(site)
	return site, nil
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
	name, slug, customDomain, themeTemplate, registrationMode, err := s.normalizeAndValidateInput(
		input.OwnerUserID,
		input.Name,
		input.Slug,
		input.CustomDomain,
		input.Status,
		input.ThemeTemplate,
		input.RegistrationMode,
	)
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
	parent, level, err := s.resolveParent(ctx, input.ParentSubSiteID, 0)
	if err != nil {
		return nil, err
	}
	site := &SubSite{}
	input.Name = name
	input.Slug = slug
	input.CustomDomain = customDomain
	s.applyInputToSite(site, input, level, normalizeStatus(input.Status), themeTemplate, registrationMode)
	if parent != nil {
		site.ParentSubSiteName = parent.Name
	}
	if err := s.repo.Create(ctx, site); err != nil {
		return nil, err
	}
	s.invalidateCaches()
	created, err := s.repo.GetByID(ctx, site.ID)
	if err != nil {
		return site, nil
	}
	return s.populateComputedFields(ctx, created, true)
}

func (s *SubSiteService) Update(ctx context.Context, input UpdateSubSiteInput) (*SubSite, error) {
	if input.ID <= 0 {
		return nil, ErrSubSiteNotFound
	}
	name, slug, customDomain, themeTemplate, registrationMode, err := s.normalizeAndValidateInput(
		input.OwnerUserID,
		input.Name,
		input.Slug,
		input.CustomDomain,
		input.Status,
		input.ThemeTemplate,
		input.RegistrationMode,
	)
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
	parent, level, err := s.resolveParent(ctx, input.ParentSubSiteID, input.ID)
	if err != nil {
		return nil, err
	}
	input.Name = name
	input.Slug = slug
	input.CustomDomain = customDomain
	createInput := CreateSubSiteInput{
		OwnerUserID:           input.OwnerUserID,
		ParentSubSiteID:       input.ParentSubSiteID,
		Name:                  input.Name,
		Slug:                  input.Slug,
		CustomDomain:          input.CustomDomain,
		Status:                input.Status,
		Mode:                  input.Mode,
		SiteLogo:              input.SiteLogo,
		SiteFavicon:           input.SiteFavicon,
		SiteSubtitle:          input.SiteSubtitle,
		Announcement:          input.Announcement,
		ContactInfo:           input.ContactInfo,
		DocURL:                input.DocURL,
		HomeContent:           input.HomeContent,
		ThemeTemplate:         input.ThemeTemplate,
		RegistrationMode:      input.RegistrationMode,
		EnableTopup:           input.EnableTopup,
		AllowSubSite:          input.AllowSubSite,
		SubSitePriceFen:       input.SubSitePriceFen,
		ConsumeRateMultiplier: input.ConsumeRateMultiplier,
		AllowOnlineTopup:      input.AllowOnlineTopup,
		AllowOfflineTopup:     input.AllowOfflineTopup,
		OwnerPaymentConfig:    input.OwnerPaymentConfig,
		SubscriptionExpiredAt: input.SubscriptionExpiredAt,
	}
	s.applyInputToSite(current, createInput, level, normalizeStatus(input.Status), themeTemplate, registrationMode)
	if parent != nil {
		current.ParentSubSiteName = parent.Name
	}
	if err := s.repo.Update(ctx, current); err != nil {
		return nil, err
	}
	// 级联停用：当 root 从非 disabled 变为 disabled 时把后代也停掉；
	// 这里不做反向级联激活（避免误激活已故意停掉的子站）。
	if current.Status == SubSiteStatusDisabled {
		if _, err := s.repo.CascadeUpdateStatus(ctx, current.ID, SubSiteStatusDisabled); err != nil {
			return nil, err
		}
	}
	s.invalidateCaches()
	updated, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return current, nil
	}
	return s.populateComputedFields(ctx, updated, true)
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

func (s *SubSiteService) invalidateCaches() {
	s.cacheMu.Lock()
	s.hostCache = make(map[string]subSiteCacheEntry)
	s.cacheMu.Unlock()
	s.boundCacheMu.Lock()
	s.boundCache = make(map[int64]boundSubSiteEntry)
	s.boundCacheMu.Unlock()
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

func (s *SubSiteService) GetByID(ctx context.Context, id int64) (*SubSite, error) {
	return s.repo.GetByID(ctx, id)
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
	cloned.AllowSubSite = site.AllowSubSite
	cloned.SubSitePriceFen = site.SubSitePriceFen
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
	if site.ThemeTemplate != "" {
		cloned.ThemeTemplate = site.ThemeTemplate
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

func (s *SubSiteService) CurrentConsumeRateMultiplier(ctx context.Context) float64 {
	site, ok := s.GetCurrent(ctx)
	if !ok || site == nil || site.Status != SubSiteStatusActive {
		return DefaultSubSiteConsumeRate
	}
	return normalizeConsumeRateMultiplier(site.ConsumeRateMultiplier)
}

func currentSubSiteConsumeRateMultiplier(ctx context.Context) float64 {
	if ctx == nil {
		return DefaultSubSiteConsumeRate
	}
	site, ok := ctx.Value(ctxkey.SubSite).(*SubSite)
	if !ok || site == nil || site.Status != SubSiteStatusActive {
		return DefaultSubSiteConsumeRate
	}
	return normalizeConsumeRateMultiplier(site.ConsumeRateMultiplier)
}

func (s *SubSiteService) readSettingInt(ctx context.Context, key string, fallback int) int {
	if s.settingRepo == nil {
		return fallback
	}
	value, err := s.settingRepo.GetValue(ctx, key)
	if err != nil || strings.TrimSpace(value) == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return parsed
}

func (s *SubSiteService) readSettingBool(ctx context.Context, key string, fallback bool) bool {
	if s.settingRepo == nil {
		return fallback
	}
	value, err := s.settingRepo.GetValue(ctx, key)
	if err != nil || strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(strings.ToLower(value)) == "true"
}

func (s *SubSiteService) readSettingString(ctx context.Context, key string, fallback string) string {
	if s.settingRepo == nil {
		return fallback
	}
	value, err := s.settingRepo.GetValue(ctx, key)
	if err != nil {
		return fallback
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
