package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	NotificationChannelWechatOfficial = "wechat_official"

	NotificationEventLowBalance      = "low_balance"
	NotificationEventBalanceDepleted = "balance_depleted"
	NotificationEventLowQuota        = "low_quota"
	NotificationEventQuotaDepleted   = "quota_depleted"
	NotificationEventLowSubscription = "low_subscription"
	NotificationEventSubscriptionEnd = "subscription_end"

	wechatOfficialAccessTokenURL  = "https://api.weixin.qq.com/cgi-bin/token"
	wechatOfficialTemplateSendURL = "https://api.weixin.qq.com/cgi-bin/message/template/send"
	wechatOfficialOAuthURL        = "https://open.weixin.qq.com/connect/oauth2/authorize"
	wechatOfficialOAuthTokenURL   = "https://api.weixin.qq.com/sns/oauth2/access_token"

	defaultWechatNotifyCooldownHours = 24
	defaultWechatStateTTL            = 10 * time.Minute
)

var (
	ErrWechatOfficialDisabled  = infraerrors.Forbidden("WECHAT_OFFICIAL_DISABLED", "wechat official account notifications are disabled")
	ErrWechatOfficialUnbound   = infraerrors.NotFound("WECHAT_OFFICIAL_UNBOUND", "wechat official account is not bound")
	ErrWechatOfficialNotReady  = infraerrors.BadRequest("WECHAT_OFFICIAL_NOT_READY", "wechat official account notification settings are incomplete")
	ErrWechatOfficialBindState = infraerrors.BadRequest("WECHAT_OFFICIAL_BIND_STATE_INVALID", "wechat official bind state is invalid")
)

type WechatBinding struct {
	ID        int64
	UserID    int64
	OpenID    string
	Enabled   bool
	BoundAt   time.Time
	UnboundAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WechatBindingStatus struct {
	Enabled      bool       `json:"enabled"`
	Configured   bool       `json:"configured"`
	Bound        bool       `json:"bound"`
	OpenIDMasked string     `json:"openid_masked,omitempty"`
	BoundAt      *time.Time `json:"bound_at,omitempty"`
}

type WechatBindURL struct {
	URL string `json:"url"`
}

type WechatOfficialConfig struct {
	Enabled                   bool
	AppID                     string
	AppSecret                 string
	TemplateLowBalance        string
	TemplateLowQuota          string
	TemplateSubscriptionLimit string
	BindRedirectURL           string
	NotifyURL                 string
	LowBalanceThreshold       float64
	LowQuotaThreshold         float64
	LowSubscriptionThreshold  float64
	CooldownHours             int
}

type WechatOfficialRepository interface {
	GetBinding(ctx context.Context, userID int64) (*WechatBinding, error)
	Bind(ctx context.Context, userID int64, openID string) error
	Unbind(ctx context.Context, userID int64) error
	ShouldSend(ctx context.Context, userID int64, channel, eventType, resourceKey string, cooldown time.Duration) (bool, error)
}

type WechatOfficialNotificationService struct {
	settingRepo SettingRepository
	repo        WechatOfficialRepository
	userRepo    UserRepository
	cfg         *config.Config
	httpClient  *http.Client

	tokenMu     sync.Mutex
	accessToken string
	tokenAppID  string
	tokenExpiry time.Time
}

func NewWechatOfficialNotificationService(
	settingRepo SettingRepository,
	repo WechatOfficialRepository,
	userRepo UserRepository,
	cfg *config.Config,
) *WechatOfficialNotificationService {
	return &WechatOfficialNotificationService{
		settingRepo: settingRepo,
		repo:        repo,
		userRepo:    userRepo,
		cfg:         cfg,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *WechatOfficialNotificationService) Status(ctx context.Context, userID int64) (*WechatBindingStatus, error) {
	cfg, err := s.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	status := &WechatBindingStatus{
		Enabled:    cfg.Enabled,
		Configured: cfg.IsConfigured(),
	}
	if !cfg.Enabled || !cfg.IsConfigured() || s.repo == nil {
		return status, nil
	}
	binding, err := s.repo.GetBinding(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrWechatOfficialUnbound) {
			return status, nil
		}
		return nil, err
	}
	status.Bound = true
	status.OpenIDMasked = maskOpenID(binding.OpenID)
	status.BoundAt = &binding.BoundAt
	return status, nil
}

func (s *WechatOfficialNotificationService) BuildBindURL(ctx context.Context, userID int64, returnTo string) (string, error) {
	cfg, err := s.GetConfig(ctx)
	if err != nil {
		return "", err
	}
	if !cfg.Enabled {
		return "", ErrWechatOfficialDisabled
	}
	if !cfg.IsConfigured() {
		return "", ErrWechatOfficialNotReady
	}
	if strings.TrimSpace(cfg.BindRedirectURL) == "" {
		return "", ErrWechatOfficialNotReady
	}

	state, err := s.SignBindState(userID, returnTo, time.Now())
	if err != nil {
		return "", err
	}
	values := url.Values{}
	values.Set("appid", cfg.AppID)
	values.Set("redirect_uri", cfg.BindRedirectURL)
	values.Set("response_type", "code")
	values.Set("scope", "snsapi_base")
	values.Set("state", state)
	return wechatOfficialOAuthURL + "?" + values.Encode() + "#wechat_redirect", nil
}

func (s *WechatOfficialNotificationService) CompleteBind(ctx context.Context, code, state string) (string, error) {
	userID, returnTo, err := s.VerifyBindState(state)
	if err != nil {
		return "", err
	}
	cfg, err := s.GetConfig(ctx)
	if err != nil {
		return "", err
	}
	if !cfg.Enabled {
		return "", ErrWechatOfficialDisabled
	}
	if !cfg.IsConfigured() {
		return "", ErrWechatOfficialNotReady
	}
	openID, err := s.exchangeOAuthCode(ctx, cfg, code)
	if err != nil {
		return "", err
	}
	if s.repo == nil {
		return "", ErrWechatOfficialNotReady
	}
	if err := s.repo.Bind(ctx, userID, openID); err != nil {
		return "", err
	}
	if strings.TrimSpace(returnTo) == "" {
		returnTo = "/profile"
	}
	return returnTo, nil
}

func (s *WechatOfficialNotificationService) Unbind(ctx context.Context, userID int64) error {
	if s.repo == nil {
		return ErrWechatOfficialUnbound
	}
	return s.repo.Unbind(ctx, userID)
}

func (s *WechatOfficialNotificationService) NotifyBalanceAfterDeduct(ctx context.Context, user *User, cost float64) {
	if user == nil || cost <= 0 {
		return
	}
	remaining := user.Balance - cost
	if s.userRepo != nil {
		if freshBalance, err := s.userRepo.GetBalance(ctx, user.ID); err == nil {
			remaining = freshBalance
		}
	}
	s.notifyBalance(ctx, user.ID, remaining)
}

func (s *WechatOfficialNotificationService) NotifyBalanceUnavailable(ctx context.Context, user *User) {
	if user == nil {
		return
	}
	s.notifyBalance(ctx, user.ID, user.Balance)
}

func (s *WechatOfficialNotificationService) notifyBalance(ctx context.Context, userID int64, remaining float64) {
	cfg, ok := s.readyConfig(ctx)
	if !ok || cfg.TemplateLowBalance == "" {
		return
	}
	threshold := cfg.LowBalanceThreshold
	if remaining > threshold {
		return
	}
	eventType := NotificationEventLowBalance
	title := "账户余额不足"
	if remaining <= 0 {
		eventType = NotificationEventBalanceDepleted
		title = "账户余额已耗尽"
	}
	data := map[string]templateValue{
		"first":    {Value: title},
		"keyword1": {Value: title},
		"keyword2": {Value: fmt.Sprintf("$%.4f", math.Max(remaining, 0))},
		"keyword3": {Value: fmt.Sprintf("$%.4f", threshold)},
		"keyword4": {Value: time.Now().Format("2006-01-02 15:04:05")},
		"remark":   {Value: "请及时充值或切换可用套餐，避免 API 请求受影响。"},
	}
	s.deliverWithCooldown(ctx, cfg, userID, eventType, "balance", cfg.TemplateLowBalance, title, data)
}

func (s *WechatOfficialNotificationService) NotifyQuotaAfterDeduct(ctx context.Context, userID int64, group *Group, remaining float64) {
	if group == nil {
		return
	}
	s.notifyQuota(ctx, userID, group, remaining)
}

func (s *WechatOfficialNotificationService) NotifyQuotaUnavailable(ctx context.Context, userID int64, group *Group) {
	s.notifyQuota(ctx, userID, group, 0)
}

func (s *WechatOfficialNotificationService) notifyQuota(ctx context.Context, userID int64, group *Group, remaining float64) {
	cfg, ok := s.readyConfig(ctx)
	if !ok || cfg.TemplateLowQuota == "" || group == nil {
		return
	}
	threshold := cfg.LowQuotaThreshold
	if remaining > threshold {
		return
	}
	eventType := NotificationEventLowQuota
	title := "额度包余额不足"
	if remaining <= 0 {
		eventType = NotificationEventQuotaDepleted
		title = "额度包已耗尽"
	}
	resourceKey := fmt.Sprintf("quota:%d", group.ID)
	data := map[string]templateValue{
		"first":    {Value: title},
		"keyword1": {Value: group.Name},
		"keyword2": {Value: fmt.Sprintf("$%.4f", math.Max(remaining, 0))},
		"keyword3": {Value: fmt.Sprintf("$%.4f", threshold)},
		"keyword4": {Value: time.Now().Format("2006-01-02 15:04:05")},
		"remark":   {Value: "请及时购买额度包或切换可用分组。"},
	}
	s.deliverWithCooldown(ctx, cfg, userID, eventType, resourceKey, cfg.TemplateLowQuota, title, data)
}

func (s *WechatOfficialNotificationService) NotifySubscriptionAfterUsage(ctx context.Context, userID int64, sub *UserSubscription, group *Group, cost float64) {
	if sub == nil || group == nil {
		return
	}
	remaining, label, ok := subscriptionRemaining(sub, group, cost)
	if !ok {
		return
	}
	s.notifySubscription(ctx, userID, sub.ID, group.Name, label, remaining)
}

func (s *WechatOfficialNotificationService) NotifySubscriptionUnavailable(ctx context.Context, userID int64, group *Group) {
	name := "订阅套餐"
	groupID := int64(0)
	if group != nil {
		name = group.Name
		groupID = group.ID
	}
	s.notifySubscription(ctx, userID, groupID, name, "可用额度", 0)
}

func (s *WechatOfficialNotificationService) notifySubscription(ctx context.Context, userID, resourceID int64, groupName, windowLabel string, remaining float64) {
	cfg, ok := s.readyConfig(ctx)
	if !ok || cfg.TemplateSubscriptionLimit == "" {
		return
	}
	threshold := cfg.LowSubscriptionThreshold
	if remaining > threshold {
		return
	}
	eventType := NotificationEventLowSubscription
	title := "套餐额度不足"
	if remaining <= 0 {
		eventType = NotificationEventSubscriptionEnd
		title = "套餐额度已用尽"
	}
	resourceKey := fmt.Sprintf("subscription:%d", resourceID)
	data := map[string]templateValue{
		"first":    {Value: title},
		"keyword1": {Value: groupName},
		"keyword2": {Value: windowLabel},
		"keyword3": {Value: fmt.Sprintf("$%.4f", math.Max(remaining, 0))},
		"keyword4": {Value: time.Now().Format("2006-01-02 15:04:05")},
		"remark":   {Value: "请及时购买新套餐或切换到余额计费。"},
	}
	s.deliverWithCooldown(ctx, cfg, userID, eventType, resourceKey, cfg.TemplateSubscriptionLimit, title, data)
}

func (s *WechatOfficialNotificationService) readyConfig(ctx context.Context) (WechatOfficialConfig, bool) {
	cfg, err := s.GetConfig(ctx)
	if err != nil || !cfg.Enabled || !cfg.IsConfigured() {
		return WechatOfficialConfig{}, false
	}
	return cfg, true
}

func (s *WechatOfficialNotificationService) deliverWithCooldown(ctx context.Context, cfg WechatOfficialConfig, userID int64, eventType, resourceKey, templateID, title string, data map[string]templateValue) {
	if s.repo == nil {
		return
	}
	binding, err := s.repo.GetBinding(ctx, userID)
	if err != nil {
		if !errors.Is(err, ErrWechatOfficialUnbound) {
			log.Printf("[WechatOfficial] get binding failed: user=%d event=%s err=%v", userID, eventType, err)
		}
		return
	}
	if binding == nil || !binding.Enabled || strings.TrimSpace(binding.OpenID) == "" {
		return
	}
	cooldown := time.Duration(cfg.CooldownHours) * time.Hour
	if cooldown <= 0 {
		cooldown = defaultWechatNotifyCooldownHours * time.Hour
	}
	ok, err := s.repo.ShouldSend(ctx, userID, NotificationChannelWechatOfficial, eventType, resourceKey, cooldown)
	if err != nil || !ok {
		if err != nil {
			log.Printf("[WechatOfficial] cooldown check failed: user=%d event=%s err=%v", userID, eventType, err)
		}
		return
	}
	if err := s.sendTemplateMessage(ctx, cfg, binding.OpenID, templateID, cfg.NotifyURL, data); err != nil {
		log.Printf("[WechatOfficial] send template failed: user=%d event=%s title=%s err=%v", userID, eventType, title, err)
	}
}

func (s *WechatOfficialNotificationService) sendTemplateMessage(ctx context.Context, cfg WechatOfficialConfig, openID, templateID, targetURL string, data map[string]templateValue) error {
	accessToken, err := s.getAccessToken(ctx, cfg)
	if err != nil {
		return err
	}
	payload := map[string]any{
		"touser":      openID,
		"template_id": templateID,
		"data":        data,
	}
	if strings.TrimSpace(targetURL) != "" {
		payload["url"] = targetURL
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	endpoint := wechatOfficialTemplateSendURL + "?access_token=" + url.QueryEscape(accessToken)
	resp, err := s.httpClient.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	var wxResp wechatAPIResponse
	_ = json.Unmarshal(respBody, &wxResp)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("wechat template send status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	if wxResp.ErrCode != 0 {
		if wxResp.ErrCode == 40001 || wxResp.ErrCode == 42001 {
			s.invalidateAccessToken()
		}
		return fmt.Errorf("wechat template send errcode=%d errmsg=%s", wxResp.ErrCode, wxResp.ErrMsg)
	}
	return nil
}

func (s *WechatOfficialNotificationService) getAccessToken(ctx context.Context, cfg WechatOfficialConfig) (string, error) {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()

	if s.accessToken != "" && s.tokenAppID == cfg.AppID && time.Now().Before(s.tokenExpiry.Add(-2*time.Minute)) {
		return s.accessToken, nil
	}

	values := url.Values{}
	values.Set("grant_type", "client_credential")
	values.Set("appid", cfg.AppID)
	values.Set("secret", cfg.AppSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wechatOfficialAccessTokenURL+"?"+values.Encode(), nil)
	if err != nil {
		return "", err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	var tokenResp wechatAccessTokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("wechat access_token status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("wechat access_token errcode=%d errmsg=%s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}
	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("wechat access_token empty response")
	}
	expiresIn := tokenResp.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 7200
	}
	s.accessToken = tokenResp.AccessToken
	s.tokenAppID = cfg.AppID
	s.tokenExpiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	return s.accessToken, nil
}

func (s *WechatOfficialNotificationService) invalidateAccessToken() {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()
	s.accessToken = ""
	s.tokenAppID = ""
	s.tokenExpiry = time.Time{}
}

func (s *WechatOfficialNotificationService) exchangeOAuthCode(ctx context.Context, cfg WechatOfficialConfig, code string) (string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", ErrWechatOfficialBindState
	}
	values := url.Values{}
	values.Set("appid", cfg.AppID)
	values.Set("secret", cfg.AppSecret)
	values.Set("code", code)
	values.Set("grant_type", "authorization_code")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wechatOfficialOAuthTokenURL+"?"+values.Encode(), nil)
	if err != nil {
		return "", err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	var oauthResp wechatOAuthTokenResponse
	if err := json.Unmarshal(respBody, &oauthResp); err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("wechat oauth token status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	if oauthResp.ErrCode != 0 {
		return "", fmt.Errorf("wechat oauth token errcode=%d errmsg=%s", oauthResp.ErrCode, oauthResp.ErrMsg)
	}
	if strings.TrimSpace(oauthResp.OpenID) == "" {
		return "", fmt.Errorf("wechat oauth token missing openid")
	}
	return strings.TrimSpace(oauthResp.OpenID), nil
}

func (s *WechatOfficialNotificationService) GetConfig(ctx context.Context) (WechatOfficialConfig, error) {
	if s.settingRepo == nil {
		return WechatOfficialConfig{}, ErrWechatOfficialNotReady
	}
	keys := []string{
		SettingKeyWechatOfficialEnabled,
		SettingKeyWechatOfficialAppID,
		SettingKeyWechatOfficialAppSecret,
		SettingKeyWechatOfficialTemplateLowBalance,
		SettingKeyWechatOfficialTemplateLowQuota,
		SettingKeyWechatOfficialTemplateSubscriptionLimit,
		SettingKeyWechatOfficialBindRedirectURL,
		SettingKeyWechatOfficialNotifyURL,
		SettingKeyWechatOfficialLowBalanceThreshold,
		SettingKeyWechatOfficialLowQuotaThreshold,
		SettingKeyWechatOfficialLowSubscriptionThreshold,
		SettingKeyWechatOfficialCooldownHours,
	}
	settings, err := s.settingRepo.GetMultiple(ctx, keys)
	if err != nil {
		return WechatOfficialConfig{}, err
	}
	cfg := WechatOfficialConfig{
		Enabled:                   settings[SettingKeyWechatOfficialEnabled] == "true",
		AppID:                     strings.TrimSpace(settings[SettingKeyWechatOfficialAppID]),
		AppSecret:                 strings.TrimSpace(settings[SettingKeyWechatOfficialAppSecret]),
		TemplateLowBalance:        strings.TrimSpace(settings[SettingKeyWechatOfficialTemplateLowBalance]),
		TemplateLowQuota:          strings.TrimSpace(settings[SettingKeyWechatOfficialTemplateLowQuota]),
		TemplateSubscriptionLimit: strings.TrimSpace(settings[SettingKeyWechatOfficialTemplateSubscriptionLimit]),
		BindRedirectURL:           strings.TrimSpace(settings[SettingKeyWechatOfficialBindRedirectURL]),
		NotifyURL:                 strings.TrimSpace(settings[SettingKeyWechatOfficialNotifyURL]),
		LowBalanceThreshold:       parseSettingFloatDefault(settings[SettingKeyWechatOfficialLowBalanceThreshold], 1),
		LowQuotaThreshold:         parseSettingFloatDefault(settings[SettingKeyWechatOfficialLowQuotaThreshold], 1),
		LowSubscriptionThreshold:  parseSettingFloatDefault(settings[SettingKeyWechatOfficialLowSubscriptionThreshold], 1),
		CooldownHours:             parseSettingIntDefault(settings[SettingKeyWechatOfficialCooldownHours], defaultWechatNotifyCooldownHours),
	}
	if cfg.CooldownHours <= 0 {
		cfg.CooldownHours = defaultWechatNotifyCooldownHours
	}
	return cfg, nil
}

func (cfg WechatOfficialConfig) IsConfigured() bool {
	return cfg.AppID != "" &&
		cfg.AppSecret != "" &&
		cfg.BindRedirectURL != "" &&
		(cfg.TemplateLowBalance != "" || cfg.TemplateLowQuota != "" || cfg.TemplateSubscriptionLimit != "")
}

func (s *WechatOfficialNotificationService) SignBindState(userID int64, returnTo string, now time.Time) (string, error) {
	if userID <= 0 {
		return "", ErrWechatOfficialBindState
	}
	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}
	payload := wechatBindStatePayload{
		UserID:   userID,
		ReturnTo: sanitizeWechatReturnTo(returnTo),
		TS:       now.Unix(),
		Nonce:    hex.EncodeToString(nonceBytes),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signature := s.signStatePayload(payloadEncoded)
	return payloadEncoded + "." + signature, nil
}

func (s *WechatOfficialNotificationService) VerifyBindState(state string) (int64, string, error) {
	state = strings.TrimSpace(state)
	parts := strings.Split(state, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return 0, "", ErrWechatOfficialBindState
	}
	expected := s.signStatePayload(parts[0])
	if !hmac.Equal([]byte(expected), []byte(parts[1])) {
		return 0, "", ErrWechatOfficialBindState
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, "", ErrWechatOfficialBindState
	}
	var payload wechatBindStatePayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return 0, "", ErrWechatOfficialBindState
	}
	if payload.UserID <= 0 || payload.TS <= 0 {
		return 0, "", ErrWechatOfficialBindState
	}
	if time.Since(time.Unix(payload.TS, 0)) > defaultWechatStateTTL {
		return 0, "", ErrWechatOfficialBindState
	}
	return payload.UserID, sanitizeWechatReturnTo(payload.ReturnTo), nil
}

func (s *WechatOfficialNotificationService) signStatePayload(payload string) string {
	secret := "wechat-official-bind-state"
	if s != nil && s.cfg != nil && strings.TrimSpace(s.cfg.JWT.Secret) != "" {
		secret = strings.TrimSpace(s.cfg.JWT.Secret)
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

type wechatBindStatePayload struct {
	UserID   int64  `json:"uid"`
	ReturnTo string `json:"return_to,omitempty"`
	TS       int64  `json:"ts"`
	Nonce    string `json:"nonce"`
}

type templateValue struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

type wechatAPIResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type wechatAccessTokenResponse struct {
	wechatAPIResponse
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type wechatOAuthTokenResponse struct {
	wechatAPIResponse
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

func subscriptionRemaining(sub *UserSubscription, group *Group, additionalCost float64) (float64, string, bool) {
	type candidate struct {
		label     string
		limit     float64
		used      float64
		remaining float64
	}
	candidates := make([]candidate, 0, 3)
	if group.HasDailyLimit() {
		remaining := *group.DailyLimitUSD - (sub.DailyUsageUSD + additionalCost)
		candidates = append(candidates, candidate{label: "日额度", limit: *group.DailyLimitUSD, used: sub.DailyUsageUSD, remaining: remaining})
	}
	if group.HasWeeklyLimit() {
		remaining := *group.WeeklyLimitUSD - (sub.WeeklyUsageUSD + additionalCost)
		candidates = append(candidates, candidate{label: "周额度", limit: *group.WeeklyLimitUSD, used: sub.WeeklyUsageUSD, remaining: remaining})
	}
	if group.HasMonthlyLimit() {
		remaining := *group.MonthlyLimitUSD - (sub.MonthlyUsageUSD + additionalCost)
		candidates = append(candidates, candidate{label: "月额度", limit: *group.MonthlyLimitUSD, used: sub.MonthlyUsageUSD, remaining: remaining})
	}
	if len(candidates) == 0 {
		return 0, "", false
	}
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.remaining < best.remaining {
			best = c
		}
	}
	return best.remaining, best.label, true
}

func sanitizeWechatReturnTo(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "/profile"
	}
	if strings.ContainsAny(raw, "\r\n") || strings.HasPrefix(raw, "//") {
		return "/profile"
	}
	if strings.HasPrefix(raw, "/") {
		return raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "/profile"
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "/profile"
	}
	return raw
}

func maskOpenID(openID string) string {
	openID = strings.TrimSpace(openID)
	if len(openID) <= 8 {
		return openID
	}
	return openID[:4] + "..." + openID[len(openID)-4:]
}

func parseSettingFloatDefault(raw string, fallback float64) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil || v < 0 {
		return fallback
	}
	return v
}

func parseSettingIntDefault(raw string, fallback int) int {
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return fallback
	}
	return v
}
