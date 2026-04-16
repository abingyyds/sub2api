package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

func (s *SubSiteService) GetPlatformConfig(ctx context.Context) (*PlatformSubSiteConfig, error) {
	cfg := &PlatformSubSiteConfig{
		EntryEnabled:         s.readSettingBool(ctx, SettingKeySubSiteEntryEnabled, false),
		Enabled:              s.readSettingBool(ctx, SettingKeySubSiteSelfServiceEnabled, true),
		ActivationPriceFen:   s.readSettingInt(ctx, SettingKeySubSiteActivationPriceFen, DefaultSubSiteActivationPriceFen),
		ValidityDays:         s.readSettingInt(ctx, SettingKeySubSiteActivationValidityDays, DefaultSubSiteValidityDays),
		MaxLevel:             s.readSettingInt(ctx, SettingKeySubSiteMaxLevel, DefaultSubSiteMaxLevel),
		DefaultThemeTemplate: normalizeThemeTemplate(s.readSettingString(ctx, SettingKeySubSiteDefaultThemeTemplate, SubSiteThemeTemplateStarter)),
		ThemeTemplates:       DefaultSubSiteThemeTemplates,
	}
	if cfg.ValidityDays <= 0 {
		cfg.ValidityDays = DefaultSubSiteValidityDays
	}
	cfg.MaxLevel = clampMaxLevel(cfg.MaxLevel)
	if !isValidThemeTemplate(cfg.DefaultThemeTemplate) {
		cfg.DefaultThemeTemplate = SubSiteThemeTemplateStarter
	}
	return cfg, nil
}

func clampMaxLevel(v int) int {
	if v <= 0 {
		return DefaultSubSiteMaxLevel
	}
	if v > MaxSubSiteLevelHardLimit {
		return MaxSubSiteLevelHardLimit
	}
	return v
}

// MaxLevel 返回 admin 配置的最大层级（带缓存读取 setting）。
func (s *SubSiteService) MaxLevel(ctx context.Context) int {
	return clampMaxLevel(s.readSettingInt(ctx, SettingKeySubSiteMaxLevel, DefaultSubSiteMaxLevel))
}

func (s *SubSiteService) UpdatePlatformConfig(ctx context.Context, input UpdatePlatformSubSiteConfigInput) (*PlatformSubSiteConfig, error) {
	template := normalizeThemeTemplate(input.DefaultThemeTemplate)
	if !isValidThemeTemplate(template) {
		return nil, ErrSubSiteInvalidThemeTemplate
	}
	if input.ActivationPriceFen < 0 {
		return nil, infraerrors.BadRequest("SUBSITE_PLATFORM_PRICE_INVALID", "activation price must be greater than or equal to 0")
	}
	if input.ValidityDays <= 0 {
		return nil, infraerrors.BadRequest("SUBSITE_PLATFORM_VALIDITY_INVALID", "validity days must be positive")
	}
	if s.settingRepo == nil {
		return nil, infraerrors.ServiceUnavailable("SETTING_REPOSITORY_UNAVAILABLE", "setting repository is unavailable")
	}
	updates := map[string]string{
		SettingKeySubSiteEntryEnabled:           fmt.Sprintf("%t", input.EntryEnabled),
		SettingKeySubSiteSelfServiceEnabled:     fmt.Sprintf("%t", input.Enabled),
		SettingKeySubSiteActivationPriceFen:     fmt.Sprintf("%d", input.ActivationPriceFen),
		SettingKeySubSiteActivationValidityDays: fmt.Sprintf("%d", input.ValidityDays),
		SettingKeySubSiteMaxLevel:               fmt.Sprintf("%d", clampMaxLevel(input.MaxLevel)),
		SettingKeySubSiteDefaultThemeTemplate:   template,
	}
	if err := s.settingRepo.SetMultiple(ctx, updates); err != nil {
		return nil, err
	}
	return s.GetPlatformConfig(ctx)
}

func (s *SubSiteService) GetOpenInfo(ctx context.Context) (*SubSiteOpenInfo, error) {
	platformCfg, err := s.GetPlatformConfig(ctx)
	if err != nil {
		return nil, err
	}

	if current, ok := s.GetCurrent(ctx); ok && current != nil {
		maxLevel := platformCfg.MaxLevel
		enabled := current.Status == SubSiteStatusActive &&
			current.AllowSubSite &&
			current.Level < maxLevel &&
			current.SubSitePriceFen > 0
		defaultTemplate := current.ThemeTemplate
		if defaultTemplate == "" {
			defaultTemplate = platformCfg.DefaultThemeTemplate
		}
		return &SubSiteOpenInfo{
			Enabled:              enabled,
			Scope:                "subsite",
			ParentSubSiteID:      &current.ID,
			ParentSubSiteName:    current.Name,
			Level:                current.Level + 1,
			MaxLevel:             maxLevel,
			PriceFen:             current.SubSitePriceFen,
			ValidityDays:         platformCfg.ValidityDays,
			Currency:             "CNY",
			AllowCustomDomain:    true,
			DefaultThemeTemplate: defaultTemplate,
			ThemeTemplates:       DefaultSubSiteThemeTemplates,
		}, nil
	}

	return &SubSiteOpenInfo{
		Enabled:              platformCfg.Enabled && platformCfg.ActivationPriceFen > 0,
		Scope:                "platform",
		Level:                1,
		MaxLevel:             platformCfg.MaxLevel,
		PriceFen:             platformCfg.ActivationPriceFen,
		ValidityDays:         platformCfg.ValidityDays,
		Currency:             "CNY",
		AllowCustomDomain:    true,
		DefaultThemeTemplate: platformCfg.DefaultThemeTemplate,
		ThemeTemplates:       DefaultSubSiteThemeTemplates,
	}, nil
}

func (s *SubSiteService) ValidateRegistration(ctx context.Context, inviteCode string) error {
	site, ok := s.GetCurrent(ctx)
	if !ok || site == nil {
		return nil
	}
	switch normalizeRegistrationMode(site.RegistrationMode) {
	case SubSiteRegistrationClosed:
		return ErrSubSiteRegistrationClosed
	case SubSiteRegistrationInvite:
		if strings.TrimSpace(inviteCode) == "" {
			return ErrSubSiteInviteRequired
		}
	}
	return nil
}

// EnsureUserScopeMatches 确认 user 的站点归属与当前请求域名匹配。
// 规则：
//   - 当前在主站（ctx 中无分站）：user 未绑定任何分站才允许。
//   - 当前在分站 A：user 未绑定 → 允许（登录成功后由调用方绑定）；user 绑到 A → 允许；
//     user 绑到其他分站 → 返回 ErrSubSiteUserScopeMismatch。
func (s *SubSiteService) EnsureUserScopeMatches(ctx context.Context, userID int64) error {
	if s == nil || s.repo == nil || userID <= 0 {
		return nil
	}
	current, _ := s.GetCurrent(ctx)
	bound, err := s.repo.GetBoundSubSiteByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if current == nil {
		if bound != nil {
			return ErrSubSiteUserScopeMismatch
		}
		return nil
	}
	if bound == nil {
		return nil
	}
	if bound.ID != current.ID {
		return ErrSubSiteUserScopeMismatch
	}
	return nil
}

// BindCurrentUserStrict 在当前分站上下文下强制绑定用户；若已绑到其他分站或当前不在分站上下文则返回错误。
func (s *SubSiteService) BindCurrentUserStrict(ctx context.Context, userID int64) error {
	site, ok := s.GetCurrent(ctx)
	if !ok || site == nil || userID <= 0 {
		return nil
	}
	bound, err := s.repo.GetBoundSubSiteByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if bound != nil && bound.ID != site.ID {
		return ErrSubSiteUserScopeMismatch
	}
	if bound != nil && bound.ID == site.ID {
		return nil
	}
	if err := s.repo.BindUser(ctx, site.ID, userID, "register"); err != nil {
		return err
	}
	s.invalidateBoundCache(userID)
	return nil
}

func (s *SubSiteService) ListOwnedSites(ctx context.Context, ownerUserID int64) ([]SubSite, error) {
	if ownerUserID <= 0 {
		return []SubSite{}, nil
	}
	items, err := s.repo.ListByOwner(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		item := items[i]
		if _, err := s.populateComputedFields(ctx, &item, true); err != nil {
			return nil, err
		}
		items[i] = item
	}
	return items, nil
}

func (s *SubSiteService) GetOwnedSite(ctx context.Context, ownerUserID int64, siteID int64) (*SubSite, error) {
	site, err := s.repo.GetByID(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if site.OwnerUserID != ownerUserID {
		return nil, ErrSubSiteForbidden
	}
	return s.populateComputedFields(ctx, site, true)
}

// AuthorizeOwner 验证 user 是否是 siteID 的 owner，返回 site（不填充 computed fields）。
// 用于中间件级鉴权 —— 要求 ownerUserID>0 且 siteID>0。
func (s *SubSiteService) AuthorizeOwner(ctx context.Context, ownerUserID int64, siteID int64) (*SubSite, error) {
	if ownerUserID <= 0 || siteID <= 0 {
		return nil, ErrSubSiteForbidden
	}
	site, err := s.repo.GetByID(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if site == nil || site.OwnerUserID != ownerUserID {
		return nil, ErrSubSiteForbidden
	}
	return site, nil
}

func (s *SubSiteService) UpdateOwnedSite(ctx context.Context, ownerUserID int64, input UpdateOwnedSubSiteInput) (*SubSite, error) {
	current, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if current.OwnerUserID != ownerUserID {
		return nil, ErrSubSiteForbidden
	}
	input.OwnerUserID = ownerUserID
	input.ParentSubSiteID = current.ParentSubSiteID
	input.Status = current.Status
	return s.Update(ctx, UpdateSubSiteInput(input))
}

func (s *SubSiteService) CreateActivationRequest(ctx context.Context, userID int64, input CreateSubSiteActivationInput) (*SubSiteActivationRequest, *SubSiteOpenInfo, error) {
	openInfo, err := s.GetOpenInfo(ctx)
	if err != nil {
		return nil, nil, err
	}
	if !openInfo.Enabled {
		return nil, nil, ErrSubSiteOpenDisabled
	}
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return nil, nil, ErrSubSiteOwnerNotFound
	}
	_, slug, customDomain, themeTemplate, registrationMode, err := s.normalizeAndValidateInput(
		userID,
		input.Name,
		input.Slug,
		input.CustomDomain,
		SubSiteStatusActive,
		input.ThemeTemplate,
		input.RegistrationMode,
	)
	if err != nil {
		return nil, nil, err
	}
	if exists, err := s.repo.ExistsBySlug(ctx, slug, 0); err != nil {
		return nil, nil, err
	} else if exists {
		return nil, nil, ErrSubSiteSlugExists
	}
	if customDomain != "" {
		if exists, err := s.repo.ExistsByDomain(ctx, customDomain, 0); err != nil {
			return nil, nil, err
		} else if exists {
			return nil, nil, ErrSubSiteDomainExists
		}
	}

	if input.AllowSubSite != nil && *input.AllowSubSite && input.SubSitePriceFen <= 0 {
		input.SubSitePriceFen = openInfo.PriceFen
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = slug
	input.CustomDomain = customDomain
	input.ThemeTemplate = themeTemplate
	input.RegistrationMode = registrationMode
	input.HomeContent = strings.TrimSpace(input.HomeContent)
	input.SiteLogo = strings.TrimSpace(input.SiteLogo)
	input.SiteFavicon = strings.TrimSpace(input.SiteFavicon)
	input.SiteSubtitle = strings.TrimSpace(input.SiteSubtitle)
	input.Announcement = strings.TrimSpace(input.Announcement)
	input.ContactInfo = strings.TrimSpace(input.ContactInfo)
	input.DocURL = strings.TrimSpace(input.DocURL)

	request := &SubSiteActivationRequest{
		UserID:          userID,
		ParentSubSiteID: openInfo.ParentSubSiteID,
		Level:           openInfo.Level,
		ValidityDays:    openInfo.ValidityDays,
		Site:            input,
	}
	return request, openInfo, nil
}

func (s *SubSiteService) SaveActivationRequest(ctx context.Context, request *SubSiteActivationRequest) error {
	if request == nil {
		return ErrSubSiteActivationNotFound
	}
	return s.repo.CreateActivationRequest(ctx, request)
}

func (s *SubSiteService) ActivatePaidOrder(ctx context.Context, order *PaymentOrder) (*SubSite, error) {
	if order == nil {
		return nil, ErrPaymentOrderNotFound
	}
	request, err := s.repo.GetActivationRequestByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	if request.ActivatedSubSiteID != nil && *request.ActivatedSubSiteID > 0 {
		site, err := s.repo.GetByID(ctx, *request.ActivatedSubSiteID)
		if err != nil {
			return nil, err
		}
		return s.populateComputedFields(ctx, site, true)
	}

	expiresAt := time.Now().Add(time.Duration(request.ValidityDays) * 24 * time.Hour)
	createInput := CreateSubSiteInput{
		OwnerUserID:           request.UserID,
		ParentSubSiteID:       request.ParentSubSiteID,
		Name:                  request.Site.Name,
		Slug:                  request.Site.Slug,
		CustomDomain:          request.Site.CustomDomain,
		Status:                SubSiteStatusActive,
		SiteLogo:              request.Site.SiteLogo,
		SiteFavicon:           request.Site.SiteFavicon,
		SiteSubtitle:          request.Site.SiteSubtitle,
		Announcement:          request.Site.Announcement,
		ContactInfo:           request.Site.ContactInfo,
		DocURL:                request.Site.DocURL,
		HomeContent:           request.Site.HomeContent,
		ThemeTemplate:         request.Site.ThemeTemplate,
		RegistrationMode:      request.Site.RegistrationMode,
		EnableTopup:           request.Site.EnableTopup,
		AllowSubSite:          request.Site.AllowSubSite,
		SubSitePriceFen:       request.Site.SubSitePriceFen,
		ConsumeRateMultiplier: request.Site.ConsumeRateMultiplier,
		SubscriptionExpiredAt: &expiresAt,
	}
	site, err := s.Create(ctx, createInput)
	if err != nil {
		return nil, err
	}
	if err := s.repo.MarkActivationRequestCompleted(ctx, order.ID, site.ID); err != nil {
		return nil, err
	}
	return site, nil
}
