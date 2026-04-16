package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/stretchr/testify/require"
)

func TestSubSiteService_ApplyPublicSettings_KeepsMainSiteConfig(t *testing.T) {
	svc := &SubSiteService{}
	base := &PublicSettings{
		SiteName:         "Main",
		SiteLogo:         "main-logo",
		ThemeTemplate:    "main-theme",
		ThemeConfig:      `{"accent":"main"}`,
		CustomConfig:     `{"hero":"main"}`,
		RegistrationMode: SubSiteRegistrationInvite,
		EnableTopup:      true,
		AllowSubSite:     false,
		SubSitePriceFen:  0,
	}
	site := &SubSite{
		ID:                    1,
		Slug:                  "branch",
		CustomDomain:          "branch.example.com",
		Name:                  "Branch",
		SiteLogo:              "branch-logo",
		SiteSubtitle:          "branch-subtitle",
		HomeContent:           "<p>branch</p>",
		ThemeTemplate:         "branch-theme",
		ThemeConfig:           `{"accent":"branch"}`,
		CustomConfig:          `{"hero":"branch"}`,
		RegistrationMode:      SubSiteRegistrationClosed,
		EnableTopup:           false,
		AllowSubSite:          true,
		SubSitePriceFen:       39900,
		ConsumeRateMultiplier: 2,
		Status:                SubSiteStatusActive,
	}
	ctx := context.WithValue(context.Background(), ctxkey.SubSite, site)

	got := svc.ApplyPublicSettings(ctx, base)
	require.NotNil(t, got)
	require.Equal(t, "Branch", got.SiteName)
	require.Equal(t, "branch-logo", got.SiteLogo)
	require.Equal(t, "branch-subtitle", got.SiteSubtitle)
	require.Equal(t, "<p>branch</p>", got.HomeContent)
	require.True(t, got.AllowSubSite)
	require.Equal(t, 39900, got.SubSitePriceFen)

	require.Equal(t, "main-theme", got.ThemeTemplate)
	require.Equal(t, `{"accent":"main"}`, got.ThemeConfig)
	require.Equal(t, `{"hero":"main"}`, got.CustomConfig)
	require.Equal(t, SubSiteRegistrationInvite, got.RegistrationMode)
	require.True(t, got.EnableTopup)
}

func TestCurrentSubSiteConsumeRateMultiplier(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxkey.SubSite, &SubSite{
		Status:                SubSiteStatusActive,
		ConsumeRateMultiplier: 2.5,
	})
	require.Equal(t, 2.5, currentSubSiteConsumeRateMultiplier(ctx))

	disabledCtx := context.WithValue(context.Background(), ctxkey.SubSite, &SubSite{
		Status:                SubSiteStatusDisabled,
		ConsumeRateMultiplier: 3,
	})
	require.Equal(t, 1.0, currentSubSiteConsumeRateMultiplier(disabledCtx))
	require.Equal(t, 1.0, currentSubSiteConsumeRateMultiplier(context.Background()))
}
