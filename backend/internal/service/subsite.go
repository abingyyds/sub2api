package service

import "time"

const (
	SubSiteStatusActive   = "active"
	SubSiteStatusDisabled = "disabled"
)

type SubSite struct {
	ID           int64     `json:"id"`
	OwnerUserID  int64     `json:"owner_user_id"`
	OwnerEmail   string    `json:"owner_email,omitempty"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	CustomDomain string    `json:"custom_domain,omitempty"`
	Status       string    `json:"status"`
	SiteLogo     string    `json:"site_logo,omitempty"`
	SiteFavicon  string    `json:"site_favicon,omitempty"`
	SiteSubtitle string    `json:"site_subtitle,omitempty"`
	Announcement string    `json:"announcement,omitempty"`
	ContactInfo  string    `json:"contact_info,omitempty"`
	DocURL       string    `json:"doc_url,omitempty"`
	HomeContent  string    `json:"home_content,omitempty"`
	ThemeConfig  string    `json:"theme_config,omitempty"`
	UserCount    int64     `json:"user_count,omitempty"`
	EntryURL     string    `json:"entry_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateSubSiteInput struct {
	OwnerUserID  int64  `json:"owner_user_id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	CustomDomain string `json:"custom_domain"`
	Status       string `json:"status"`
	SiteLogo     string `json:"site_logo"`
	SiteFavicon  string `json:"site_favicon"`
	SiteSubtitle string `json:"site_subtitle"`
	Announcement string `json:"announcement"`
	ContactInfo  string `json:"contact_info"`
	DocURL       string `json:"doc_url"`
	HomeContent  string `json:"home_content"`
	ThemeConfig  string `json:"theme_config"`
}

type UpdateSubSiteInput struct {
	ID           int64  `json:"id"`
	OwnerUserID  int64  `json:"owner_user_id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	CustomDomain string `json:"custom_domain"`
	Status       string `json:"status"`
	SiteLogo     string `json:"site_logo"`
	SiteFavicon  string `json:"site_favicon"`
	SiteSubtitle string `json:"site_subtitle"`
	Announcement string `json:"announcement"`
	ContactInfo  string `json:"contact_info"`
	DocURL       string `json:"doc_url"`
	HomeContent  string `json:"home_content"`
	ThemeConfig  string `json:"theme_config"`
}
