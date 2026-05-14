package service

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/util/responseheaders"
	"github.com/Wei-Shaw/sub2api/internal/util/urlvalidator"
	"github.com/gin-gonic/gin"
)

const (
	// ChatGPT internal API for OAuth accounts
	chatgptCodexURL = "https://chatgpt.com/backend-api/codex/responses"
	// OpenAI Platform base URL for API Key accounts.
	openaiPlatformBaseURL  = "https://api.openai.com"
	openaiStickySessionTTL = time.Hour // 粘性会话TTL
)

const (
	OpenAIEndpointResponses         = "/responses"
	OpenAIEndpointChatCompletions   = "/chat/completions"
	OpenAIEndpointImagesGenerations = "/images/generations"
	OpenAIEndpointImagesEdits       = "/images/edits"
)

// openaiSSEDataRe matches SSE data lines with optional whitespace after colon.
// Some upstream APIs return non-standard "data:" without space (should be "data: ").
var openaiSSEDataRe = regexp.MustCompile(`^data:\s*`)

// OpenAI allowed headers whitelist (for non-OAuth accounts)
var openaiAllowedHeaders = map[string]bool{
	"accept":          true,
	"accept-language": true,
	"content-type":    true,
	"conversation_id": true,
	"user-agent":      true,
	"originator":      true,
	"session_id":      true,
}

func OpenAIEndpointFromPath(path string) string {
	switch {
	case strings.HasSuffix(path, OpenAIEndpointChatCompletions):
		return OpenAIEndpointChatCompletions
	case strings.HasSuffix(path, OpenAIEndpointImagesGenerations):
		return OpenAIEndpointImagesGenerations
	case strings.HasSuffix(path, OpenAIEndpointImagesEdits):
		return OpenAIEndpointImagesEdits
	default:
		return OpenAIEndpointResponses
	}
}

func OpenAIUpstreamEndpoint(endpoint string) string {
	if endpoint == OpenAIEndpointChatCompletions {
		return OpenAIEndpointResponses
	}
	return endpoint
}

func IsOpenAIImageEndpoint(endpoint string) bool {
	return endpoint == OpenAIEndpointImagesGenerations || endpoint == OpenAIEndpointImagesEdits
}

// OpenAICodexUsageSnapshot represents Codex API usage limits from response headers
type OpenAICodexUsageSnapshot struct {
	PrimaryUsedPercent          *float64 `json:"primary_used_percent,omitempty"`
	PrimaryResetAfterSeconds    *int     `json:"primary_reset_after_seconds,omitempty"`
	PrimaryWindowMinutes        *int     `json:"primary_window_minutes,omitempty"`
	SecondaryUsedPercent        *float64 `json:"secondary_used_percent,omitempty"`
	SecondaryResetAfterSeconds  *int     `json:"secondary_reset_after_seconds,omitempty"`
	SecondaryWindowMinutes      *int     `json:"secondary_window_minutes,omitempty"`
	PrimaryOverSecondaryPercent *float64 `json:"primary_over_secondary_percent,omitempty"`
	UpdatedAt                   string   `json:"updated_at,omitempty"`
}

// NormalizedCodexLimits contains normalized 5h/7d rate limit data
type NormalizedCodexLimits struct {
	Used5hPercent   *float64
	Reset5hSeconds  *int
	Window5hMinutes *int
	Used7dPercent   *float64
	Reset7dSeconds  *int
	Window7dMinutes *int
}

// Normalize converts primary/secondary fields to canonical 5h/7d fields.
// Strategy: Compare window_minutes to determine which is 5h vs 7d.
// Returns nil if snapshot is nil or has no useful data.
func (s *OpenAICodexUsageSnapshot) Normalize() *NormalizedCodexLimits {
	if s == nil {
		return nil
	}

	result := &NormalizedCodexLimits{}

	primaryMins := 0
	secondaryMins := 0
	hasPrimaryWindow := false
	hasSecondaryWindow := false

	if s.PrimaryWindowMinutes != nil {
		primaryMins = *s.PrimaryWindowMinutes
		hasPrimaryWindow = true
	}
	if s.SecondaryWindowMinutes != nil {
		secondaryMins = *s.SecondaryWindowMinutes
		hasSecondaryWindow = true
	}

	// Determine mapping based on window_minutes
	use5hFromPrimary := false
	use7dFromPrimary := false

	if hasPrimaryWindow && hasSecondaryWindow {
		// Both known: smaller window is 5h, larger is 7d
		if primaryMins < secondaryMins {
			use5hFromPrimary = true
		} else {
			use7dFromPrimary = true
		}
	} else if hasPrimaryWindow {
		// Only primary known: classify by threshold (<=360 min = 6h -> 5h window)
		if primaryMins <= 360 {
			use5hFromPrimary = true
		} else {
			use7dFromPrimary = true
		}
	} else if hasSecondaryWindow {
		// Only secondary known: classify by threshold
		if secondaryMins <= 360 {
			// 5h from secondary, so primary (if any data) is 7d
			use7dFromPrimary = true
		} else {
			// 7d from secondary, so primary (if any data) is 5h
			use5hFromPrimary = true
		}
	} else {
		// No window_minutes: fall back to legacy assumption (primary=7d, secondary=5h)
		use7dFromPrimary = true
	}

	// Assign values
	if use5hFromPrimary {
		result.Used5hPercent = s.PrimaryUsedPercent
		result.Reset5hSeconds = s.PrimaryResetAfterSeconds
		result.Window5hMinutes = s.PrimaryWindowMinutes
		result.Used7dPercent = s.SecondaryUsedPercent
		result.Reset7dSeconds = s.SecondaryResetAfterSeconds
		result.Window7dMinutes = s.SecondaryWindowMinutes
	} else if use7dFromPrimary {
		result.Used7dPercent = s.PrimaryUsedPercent
		result.Reset7dSeconds = s.PrimaryResetAfterSeconds
		result.Window7dMinutes = s.PrimaryWindowMinutes
		result.Used5hPercent = s.SecondaryUsedPercent
		result.Reset5hSeconds = s.SecondaryResetAfterSeconds
		result.Window5hMinutes = s.SecondaryWindowMinutes
	}

	return result
}

// OpenAIUsage represents OpenAI API response usage
type OpenAIUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
}

// OpenAIForwardResult represents the result of forwarding
type OpenAIForwardResult struct {
	RequestID    string
	Usage        OpenAIUsage
	Model        string
	Stream       bool
	Duration     time.Duration
	FirstTokenMs *int
	ImageCount   int
	ImageSize    string
}

type openAIRequestPayload struct {
	Model          string
	Stream         bool
	PromptCacheKey string
	ImageSize      string
	ContentType    string
	JSONBody       map[string]any
	MultipartBody  *openAIMultipartBody
}

type openAIMultipartBody struct {
	Parts []openAIMultipartPart
}

type openAIMultipartPart struct {
	Header   textproto.MIMEHeader
	FormName string
	FileName string
	Data     []byte
}

type openAIImageResponseMeta struct {
	Usage      OpenAIUsage
	ImageCount int
	ImageSize  string
}

// OpenAIGatewayService handles OpenAI API gateway operations
type OpenAIGatewayService struct {
	accountRepo         AccountRepository
	usageLogRepo        UsageLogRepository
	userRepo            UserRepository
	userSubRepo         UserSubscriptionRepository
	quotaPackageRepo    QuotaPackageRepository
	orgRepo             OrganizationRepository
	orgMemberRepo       OrgMemberRepository
	orgProjectRepo      OrgProjectRepository
	auditService        *OrgAuditService
	cache               GatewayCache
	cfg                 *config.Config
	schedulerSnapshot   *SchedulerSnapshotService
	concurrencyService  *ConcurrencyService
	billingService      *BillingService
	rateLimitService    *RateLimitService
	billingCacheService *BillingCacheService
	httpUpstream        HTTPUpstream
	deferredService     *DeferredService
	openAITokenProvider *OpenAITokenProvider
	toolCorrector       *CodexToolCorrector
	subSiteService      *SubSiteService // 分站池扣费
	wechatNotifyService *WechatOfficialNotificationService
}

// NewOpenAIGatewayService creates a new OpenAIGatewayService
func NewOpenAIGatewayService(
	accountRepo AccountRepository,
	usageLogRepo UsageLogRepository,
	userRepo UserRepository,
	userSubRepo UserSubscriptionRepository,
	quotaPackageRepo QuotaPackageRepository,
	orgRepo OrganizationRepository,
	orgMemberRepo OrgMemberRepository,
	orgProjectRepo OrgProjectRepository,
	auditService *OrgAuditService,
	cache GatewayCache,
	cfg *config.Config,
	schedulerSnapshot *SchedulerSnapshotService,
	concurrencyService *ConcurrencyService,
	billingService *BillingService,
	rateLimitService *RateLimitService,
	billingCacheService *BillingCacheService,
	httpUpstream HTTPUpstream,
	deferredService *DeferredService,
	openAITokenProvider *OpenAITokenProvider,
	subSiteService *SubSiteService,
	wechatNotifyService *WechatOfficialNotificationService,
) *OpenAIGatewayService {
	return &OpenAIGatewayService{
		accountRepo:         accountRepo,
		usageLogRepo:        usageLogRepo,
		userRepo:            userRepo,
		userSubRepo:         userSubRepo,
		quotaPackageRepo:    quotaPackageRepo,
		orgRepo:             orgRepo,
		orgMemberRepo:       orgMemberRepo,
		orgProjectRepo:      orgProjectRepo,
		auditService:        auditService,
		cache:               cache,
		cfg:                 cfg,
		schedulerSnapshot:   schedulerSnapshot,
		concurrencyService:  concurrencyService,
		billingService:      billingService,
		rateLimitService:    rateLimitService,
		billingCacheService: billingCacheService,
		httpUpstream:        httpUpstream,
		deferredService:     deferredService,
		openAITokenProvider: openAITokenProvider,
		toolCorrector:       NewCodexToolCorrector(),
		subSiteService:      subSiteService,
		wechatNotifyService: wechatNotifyService,
	}
}

// GenerateSessionHash generates a sticky-session hash for OpenAI requests.
//
// Priority:
//  1. Header: session_id
//  2. Header: conversation_id
//  3. Body:   prompt_cache_key (opencode)
func (s *OpenAIGatewayService) GenerateSessionHash(c *gin.Context, reqBody map[string]any) string {
	if c == nil {
		return ""
	}

	sessionID := strings.TrimSpace(c.GetHeader("session_id"))
	if sessionID == "" {
		sessionID = strings.TrimSpace(c.GetHeader("conversation_id"))
	}
	if sessionID == "" && reqBody != nil {
		if v, ok := reqBody["prompt_cache_key"].(string); ok {
			sessionID = strings.TrimSpace(v)
		}
	}
	if sessionID == "" {
		return ""
	}

	hash := sha256.Sum256([]byte(sessionID))
	return hex.EncodeToString(hash[:])
}

// BindStickySession sets session -> account binding with standard TTL.
func (s *OpenAIGatewayService) BindStickySession(ctx context.Context, groupID *int64, sessionHash string, accountID int64) error {
	if sessionHash == "" || accountID <= 0 {
		return nil
	}
	return s.cache.SetSessionAccountID(ctx, derefGroupID(groupID), "openai:"+sessionHash, accountID, openaiStickySessionTTL)
}

// SelectAccount selects an OpenAI account with sticky session support
func (s *OpenAIGatewayService) SelectAccount(ctx context.Context, groupID *int64, sessionHash string) (*Account, error) {
	return s.SelectAccountForModel(ctx, groupID, sessionHash, "")
}

// SelectAccountForModel selects an account supporting the requested model
func (s *OpenAIGatewayService) SelectAccountForModel(ctx context.Context, groupID *int64, sessionHash string, requestedModel string) (*Account, error) {
	return s.SelectAccountForModelWithExclusions(ctx, groupID, sessionHash, requestedModel, nil)
}

// SelectAccountForModelWithExclusions selects an account supporting the requested model while excluding specified accounts.
// SelectAccountForModelWithExclusions 选择支持指定模型的账号，同时排除指定的账号。
func (s *OpenAIGatewayService) SelectAccountForModelWithExclusions(ctx context.Context, groupID *int64, sessionHash string, requestedModel string, excludedIDs map[int64]struct{}) (*Account, error) {
	cacheKey := "openai:" + sessionHash

	// 1. 尝试粘性会话命中
	// Try sticky session hit
	if account := s.tryStickySessionHit(ctx, groupID, sessionHash, cacheKey, requestedModel, excludedIDs); account != nil {
		return account, nil
	}

	// 2. 获取可调度的 OpenAI 账号
	// Get schedulable OpenAI accounts
	accounts, err := s.listSchedulableAccounts(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("query accounts failed: %w", err)
	}

	// 3. 按优先级 + LRU 选择最佳账号
	// Select by priority + LRU
	selected := s.selectBestAccount(accounts, requestedModel, excludedIDs)

	if selected == nil {
		if requestedModel != "" {
			return nil, fmt.Errorf("no available OpenAI accounts supporting model: %s", requestedModel)
		}
		return nil, errors.New("no available OpenAI accounts")
	}

	// 4. 设置粘性会话绑定
	// Set sticky session binding
	if sessionHash != "" {
		_ = s.cache.SetSessionAccountID(ctx, derefGroupID(groupID), cacheKey, selected.ID, openaiStickySessionTTL)
	}

	return selected, nil
}

// tryStickySessionHit 尝试从粘性会话获取账号。
// 如果命中且账号可用则返回账号；如果账号不可用则清理会话并返回 nil。
//
// tryStickySessionHit attempts to get account from sticky session.
// Returns account if hit and usable; clears session and returns nil if account is unavailable.
func (s *OpenAIGatewayService) tryStickySessionHit(ctx context.Context, groupID *int64, sessionHash, cacheKey, requestedModel string, excludedIDs map[int64]struct{}) *Account {
	if sessionHash == "" {
		return nil
	}

	accountID, err := s.cache.GetSessionAccountID(ctx, derefGroupID(groupID), cacheKey)
	if err != nil || accountID <= 0 {
		return nil
	}

	if _, excluded := excludedIDs[accountID]; excluded {
		return nil
	}

	account, err := s.getSchedulableAccount(ctx, accountID)
	if err != nil {
		return nil
	}

	// 检查账号是否需要清理粘性会话
	// Check if sticky session should be cleared
	if shouldClearStickySession(account) {
		_ = s.cache.DeleteSessionAccountID(ctx, derefGroupID(groupID), cacheKey)
		return nil
	}

	// 验证账号是否可用于当前请求
	// Verify account is usable for current request
	if !account.IsSchedulable() || !account.IsOpenAI() {
		return nil
	}
	if requestedModel != "" && !account.IsModelSupported(requestedModel) {
		return nil
	}

	// 刷新会话 TTL 并返回账号
	// Refresh session TTL and return account
	_ = s.cache.RefreshSessionTTL(ctx, derefGroupID(groupID), cacheKey, openaiStickySessionTTL)
	return account
}

// selectBestAccount 从候选账号中选择最佳账号（优先级 + LRU）。
// 返回 nil 表示无可用账号。
//
// selectBestAccount selects the best account from candidates (priority + LRU).
// Returns nil if no available account.
func (s *OpenAIGatewayService) selectBestAccount(accounts []Account, requestedModel string, excludedIDs map[int64]struct{}) *Account {
	var selected *Account

	for i := range accounts {
		acc := &accounts[i]

		// 跳过被排除的账号
		// Skip excluded accounts
		if _, excluded := excludedIDs[acc.ID]; excluded {
			continue
		}

		// 调度器快照可能暂时过时，这里重新检查可调度性和平台
		// Scheduler snapshots can be temporarily stale; re-check schedulability and platform
		if !acc.IsSchedulable() || !acc.IsOpenAI() {
			continue
		}

		// 检查模型支持
		// Check model support
		if requestedModel != "" && !acc.IsModelSupported(requestedModel) {
			continue
		}

		// 选择优先级最高且最久未使用的账号
		// Select highest priority and least recently used
		if selected == nil {
			selected = acc
			continue
		}

		if s.isBetterAccount(acc, selected) {
			selected = acc
		}
	}

	return selected
}

// isBetterAccount 判断 candidate 是否比 current 更优。
// 规则：优先级更高（数值更小）优先；同优先级时，未使用过的优先，其次是最久未使用的。
//
// isBetterAccount checks if candidate is better than current.
// Rules: higher priority (lower value) wins; same priority: never used > least recently used.
func (s *OpenAIGatewayService) isBetterAccount(candidate, current *Account) bool {
	// 优先级更高（数值更小）
	// Higher priority (lower value)
	if candidate.Priority < current.Priority {
		return true
	}
	if candidate.Priority > current.Priority {
		return false
	}

	// 同优先级，比较最后使用时间
	// Same priority, compare last used time
	switch {
	case candidate.LastUsedAt == nil && current.LastUsedAt != nil:
		// candidate 从未使用，优先
		return true
	case candidate.LastUsedAt != nil && current.LastUsedAt == nil:
		// current 从未使用，保持
		return false
	case candidate.LastUsedAt == nil && current.LastUsedAt == nil:
		// 都未使用，保持
		return false
	default:
		// 都使用过，选择最久未使用的
		return candidate.LastUsedAt.Before(*current.LastUsedAt)
	}
}

// SelectAccountWithLoadAwareness selects an account with load-awareness and wait plan.
func (s *OpenAIGatewayService) SelectAccountWithLoadAwareness(ctx context.Context, groupID *int64, sessionHash string, requestedModel string, excludedIDs map[int64]struct{}) (*AccountSelectionResult, error) {
	cfg := s.schedulingConfig()
	var stickyAccountID int64
	if sessionHash != "" && s.cache != nil {
		if accountID, err := s.cache.GetSessionAccountID(ctx, derefGroupID(groupID), "openai:"+sessionHash); err == nil {
			stickyAccountID = accountID
		}
	}
	if s.concurrencyService == nil || !cfg.LoadBatchEnabled {
		account, err := s.SelectAccountForModelWithExclusions(ctx, groupID, sessionHash, requestedModel, excludedIDs)
		if err != nil {
			return nil, err
		}
		result, err := s.tryAcquireAccountSlot(ctx, account.ID, account.Concurrency)
		if err == nil && result.Acquired {
			return &AccountSelectionResult{
				Account:     account,
				Acquired:    true,
				ReleaseFunc: result.ReleaseFunc,
			}, nil
		}
		if stickyAccountID > 0 && stickyAccountID == account.ID && s.concurrencyService != nil {
			waitingCount, _ := s.concurrencyService.GetAccountWaitingCount(ctx, account.ID)
			if waitingCount < cfg.StickySessionMaxWaiting {
				return &AccountSelectionResult{
					Account: account,
					WaitPlan: &AccountWaitPlan{
						AccountID:      account.ID,
						MaxConcurrency: account.Concurrency,
						Timeout:        cfg.StickySessionWaitTimeout,
						MaxWaiting:     cfg.StickySessionMaxWaiting,
					},
				}, nil
			}
		}
		return &AccountSelectionResult{
			Account: account,
			WaitPlan: &AccountWaitPlan{
				AccountID:      account.ID,
				MaxConcurrency: account.Concurrency,
				Timeout:        cfg.FallbackWaitTimeout,
				MaxWaiting:     cfg.FallbackMaxWaiting,
			},
		}, nil
	}

	accounts, err := s.listSchedulableAccounts(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, errors.New("no available accounts")
	}

	isExcluded := func(accountID int64) bool {
		if excludedIDs == nil {
			return false
		}
		_, excluded := excludedIDs[accountID]
		return excluded
	}

	// ============ Layer 1: Sticky session ============
	if sessionHash != "" {
		accountID, err := s.cache.GetSessionAccountID(ctx, derefGroupID(groupID), "openai:"+sessionHash)
		if err == nil && accountID > 0 && !isExcluded(accountID) {
			account, err := s.getSchedulableAccount(ctx, accountID)
			if err == nil {
				clearSticky := shouldClearStickySession(account)
				if clearSticky {
					_ = s.cache.DeleteSessionAccountID(ctx, derefGroupID(groupID), "openai:"+sessionHash)
				}
				if !clearSticky && account.IsSchedulable() && account.IsOpenAI() &&
					(requestedModel == "" || account.IsModelSupported(requestedModel)) {
					result, err := s.tryAcquireAccountSlot(ctx, accountID, account.Concurrency)
					if err == nil && result.Acquired {
						_ = s.cache.RefreshSessionTTL(ctx, derefGroupID(groupID), "openai:"+sessionHash, openaiStickySessionTTL)
						return &AccountSelectionResult{
							Account:     account,
							Acquired:    true,
							ReleaseFunc: result.ReleaseFunc,
						}, nil
					}

					waitingCount, _ := s.concurrencyService.GetAccountWaitingCount(ctx, accountID)
					if waitingCount < cfg.StickySessionMaxWaiting {
						return &AccountSelectionResult{
							Account: account,
							WaitPlan: &AccountWaitPlan{
								AccountID:      accountID,
								MaxConcurrency: account.Concurrency,
								Timeout:        cfg.StickySessionWaitTimeout,
								MaxWaiting:     cfg.StickySessionMaxWaiting,
							},
						}, nil
					}
				}
			}
		}
	}

	// ============ Layer 2: Load-aware selection ============
	candidates := make([]*Account, 0, len(accounts))
	for i := range accounts {
		acc := &accounts[i]
		if isExcluded(acc.ID) {
			continue
		}
		// Scheduler snapshots can be temporarily stale (bucket rebuild is throttled);
		// re-check schedulability here so recently rate-limited/overloaded accounts
		// are not selected again before the bucket is rebuilt.
		if !acc.IsSchedulable() {
			continue
		}
		if requestedModel != "" && !acc.IsModelSupported(requestedModel) {
			continue
		}
		candidates = append(candidates, acc)
	}

	if len(candidates) == 0 {
		return nil, errors.New("no available accounts")
	}

	accountLoads := make([]AccountWithConcurrency, 0, len(candidates))
	for _, acc := range candidates {
		accountLoads = append(accountLoads, AccountWithConcurrency{
			ID:             acc.ID,
			MaxConcurrency: acc.Concurrency,
		})
	}

	loadMap, err := s.concurrencyService.GetAccountsLoadBatch(ctx, accountLoads)
	if err != nil {
		ordered := append([]*Account(nil), candidates...)
		sortAccountsByPriorityAndLastUsed(ordered, false)
		for _, acc := range ordered {
			result, err := s.tryAcquireAccountSlot(ctx, acc.ID, acc.Concurrency)
			if err == nil && result.Acquired {
				if sessionHash != "" {
					_ = s.cache.SetSessionAccountID(ctx, derefGroupID(groupID), "openai:"+sessionHash, acc.ID, openaiStickySessionTTL)
				}
				return &AccountSelectionResult{
					Account:     acc,
					Acquired:    true,
					ReleaseFunc: result.ReleaseFunc,
				}, nil
			}
		}
	} else {
		type accountWithLoad struct {
			account  *Account
			loadInfo *AccountLoadInfo
		}
		var available []accountWithLoad
		for _, acc := range candidates {
			loadInfo := loadMap[acc.ID]
			if loadInfo == nil {
				loadInfo = &AccountLoadInfo{AccountID: acc.ID}
			}
			available = append(available, accountWithLoad{
				account:  acc,
				loadInfo: loadInfo,
			})
		}

		if len(available) > 0 {
			sort.SliceStable(available, func(i, j int) bool {
				a, b := available[i], available[j]
				if a.account.Priority != b.account.Priority {
					return a.account.Priority < b.account.Priority
				}
				switch {
				case a.account.LastUsedAt == nil && b.account.LastUsedAt != nil:
					return true
				case a.account.LastUsedAt != nil && b.account.LastUsedAt == nil:
					return false
				case a.account.LastUsedAt == nil && b.account.LastUsedAt == nil:
					return false
				default:
					return a.account.LastUsedAt.Before(*b.account.LastUsedAt)
				}
			})

			for _, item := range available {
				result, err := s.tryAcquireAccountSlot(ctx, item.account.ID, item.account.Concurrency)
				if err == nil && result.Acquired {
					if sessionHash != "" {
						_ = s.cache.SetSessionAccountID(ctx, derefGroupID(groupID), "openai:"+sessionHash, item.account.ID, openaiStickySessionTTL)
					}
					return &AccountSelectionResult{
						Account:     item.account,
						Acquired:    true,
						ReleaseFunc: result.ReleaseFunc,
					}, nil
				}
			}
		}
	}

	// ============ Layer 3: Fallback wait ============
	sortAccountsByPriorityAndLastUsed(candidates, false)
	for _, acc := range candidates {
		return &AccountSelectionResult{
			Account: acc,
			WaitPlan: &AccountWaitPlan{
				AccountID:      acc.ID,
				MaxConcurrency: acc.Concurrency,
				Timeout:        cfg.FallbackWaitTimeout,
				MaxWaiting:     cfg.FallbackMaxWaiting,
			},
		}, nil
	}

	return nil, errors.New("no available accounts")
}

func (s *OpenAIGatewayService) listSchedulableAccounts(ctx context.Context, groupID *int64) ([]Account, error) {
	if s.schedulerSnapshot != nil {
		accounts, _, err := s.schedulerSnapshot.ListSchedulableAccounts(ctx, groupID, PlatformOpenAI, false)
		return accounts, err
	}
	var accounts []Account
	var err error
	if s.cfg != nil && s.cfg.RunMode == config.RunModeSimple {
		accounts, err = s.accountRepo.ListSchedulableByPlatform(ctx, PlatformOpenAI)
	} else if groupID != nil {
		accounts, err = s.accountRepo.ListSchedulableByGroupIDAndPlatform(ctx, *groupID, PlatformOpenAI)
	} else {
		accounts, err = s.accountRepo.ListSchedulableByPlatform(ctx, PlatformOpenAI)
	}
	if err != nil {
		return nil, fmt.Errorf("query accounts failed: %w", err)
	}
	return accounts, nil
}

func (s *OpenAIGatewayService) tryAcquireAccountSlot(ctx context.Context, accountID int64, maxConcurrency int) (*AcquireResult, error) {
	if s.concurrencyService == nil {
		return &AcquireResult{Acquired: true, ReleaseFunc: func() {}}, nil
	}
	return s.concurrencyService.AcquireAccountSlot(ctx, accountID, maxConcurrency)
}

func (s *OpenAIGatewayService) getSchedulableAccount(ctx context.Context, accountID int64) (*Account, error) {
	if s.schedulerSnapshot != nil {
		return s.schedulerSnapshot.GetAccount(ctx, accountID)
	}
	return s.accountRepo.GetByID(ctx, accountID)
}

func (s *OpenAIGatewayService) schedulingConfig() config.GatewaySchedulingConfig {
	if s.cfg != nil {
		return s.cfg.Gateway.Scheduling
	}
	return config.GatewaySchedulingConfig{
		StickySessionMaxWaiting:  3,
		StickySessionWaitTimeout: 45 * time.Second,
		FallbackWaitTimeout:      30 * time.Second,
		FallbackMaxWaiting:       100,
		LoadBatchEnabled:         true,
		SlotCleanupInterval:      30 * time.Second,
	}
}

// GetAccessToken gets the access token for an OpenAI account
func (s *OpenAIGatewayService) GetAccessToken(ctx context.Context, account *Account) (string, string, error) {
	switch account.Type {
	case AccountTypeOAuth:
		// 使用 TokenProvider 获取缓存的 token
		if s.openAITokenProvider != nil {
			accessToken, err := s.openAITokenProvider.GetAccessToken(ctx, account)
			if err != nil {
				return "", "", err
			}
			return accessToken, "oauth", nil
		}
		// 降级：TokenProvider 未配置时直接从账号读取
		accessToken := account.GetOpenAIAccessToken()
		if accessToken == "" {
			return "", "", errors.New("access_token not found in credentials")
		}
		return accessToken, "oauth", nil
	case AccountTypeAPIKey:
		apiKey := account.GetOpenAIApiKey()
		if apiKey == "" {
			return "", "", errors.New("api_key not found in credentials")
		}
		return apiKey, "apikey", nil
	default:
		return "", "", fmt.Errorf("unsupported account type: %s", account.Type)
	}
}

func (s *OpenAIGatewayService) shouldFailoverUpstreamError(statusCode int) bool {
	return statusCode >= 400
}

func (s *OpenAIGatewayService) handleFailoverSideEffects(ctx context.Context, resp *http.Response, account *Account) {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	s.rateLimitService.HandleUpstreamError(ctx, account, resp.StatusCode, resp.Header, body)
}

func parseOpenAIRequestPayload(body []byte, contentType string) (*openAIRequestPayload, error) {
	payload := &openAIRequestPayload{
		ContentType: strings.TrimSpace(contentType),
	}

	mediaType := strings.ToLower(strings.TrimSpace(contentType))
	if mediaType != "" {
		parsedMediaType, _, err := mime.ParseMediaType(contentType)
		if err == nil {
			mediaType = strings.ToLower(parsedMediaType)
		}
	}

	if strings.HasPrefix(mediaType, "multipart/form-data") {
		multipartBody, err := parseOpenAIMultipartBody(body, contentType)
		if err != nil {
			return nil, err
		}
		payload.MultipartBody = multipartBody
		payload.Model = strings.TrimSpace(multipartBody.firstValue("model"))
		payload.Stream = parseMultipartBool(multipartBody.firstValue("stream"))
		payload.ImageSize = normalizeOpenAIImageSize(multipartBody.firstValue("size"))
		return payload, nil
	}

	var reqBody map[string]any
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}
	payload.JSONBody = reqBody
	payload.Model, _ = reqBody["model"].(string)
	payload.Stream, _ = reqBody["stream"].(bool)
	if v, ok := reqBody["prompt_cache_key"].(string); ok {
		payload.PromptCacheKey = strings.TrimSpace(v)
	}
	if v, ok := reqBody["size"].(string); ok {
		payload.ImageSize = normalizeOpenAIImageSize(v)
	}
	return payload, nil
}

func parseOpenAIMultipartBody(body []byte, contentType string) (*openAIMultipartBody, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("parse multipart content-type: %w", err)
	}
	if !strings.EqualFold(mediaType, "multipart/form-data") {
		return nil, fmt.Errorf("unsupported content-type: %s", contentType)
	}
	boundary := strings.TrimSpace(params["boundary"])
	if boundary == "" {
		return nil, errors.New("missing multipart boundary")
	}

	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	result := &openAIMultipartBody{}
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			return result, nil
		}
		if err != nil {
			return nil, fmt.Errorf("read multipart body: %w", err)
		}
		data, readErr := io.ReadAll(part)
		_ = part.Close()
		if readErr != nil {
			return nil, fmt.Errorf("read multipart part: %w", readErr)
		}

		result.Parts = append(result.Parts, openAIMultipartPart{
			Header:   cloneMIMEHeader(part.Header),
			FormName: part.FormName(),
			FileName: part.FileName(),
			Data:     data,
		})
	}
}

func (b *openAIMultipartBody) firstValue(name string) string {
	for _, part := range b.Parts {
		if part.FileName == "" && part.FormName == name {
			return strings.TrimSpace(string(part.Data))
		}
	}
	return ""
}

func (b *openAIMultipartBody) rebuildBody(fieldOverrides map[string]string) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for _, part := range b.Parts {
		header := cloneMIMEHeader(part.Header)
		header.Del("Content-Length")

		data := part.Data
		if part.FileName == "" {
			if updated, ok := fieldOverrides[part.FormName]; ok {
				data = []byte(updated)
			}
		}

		partWriter, err := writer.CreatePart(header)
		if err != nil {
			return nil, "", fmt.Errorf("create multipart part: %w", err)
		}
		if _, err := partWriter.Write(data); err != nil {
			return nil, "", fmt.Errorf("write multipart part: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("finalize multipart body: %w", err)
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

func cloneMIMEHeader(header textproto.MIMEHeader) textproto.MIMEHeader {
	cloned := make(textproto.MIMEHeader, len(header))
	for key, values := range header {
		cloned[key] = append([]string(nil), values...)
	}
	return cloned
}

func parseMultipartBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// Forward forwards request to OpenAI API
func (s *OpenAIGatewayService) Forward(ctx context.Context, c *gin.Context, account *Account, body []byte) (*OpenAIForwardResult, error) {
	startTime := time.Now()
	endpoint := OpenAIEndpointFromPath(c.Request.URL.Path)
	upstreamEndpoint := OpenAIUpstreamEndpoint(endpoint)
	reqPayload, err := parseOpenAIRequestPayload(body, c.GetHeader("Content-Type"))
	if err != nil {
		return nil, err
	}

	// Track if body needs re-serialization
	bodyModified := false
	if IsOpenAIImageEndpoint(endpoint) && strings.TrimSpace(reqPayload.Model) == "" {
		reqPayload.Model = DefaultOpenAIImageModel
		if reqPayload.JSONBody != nil {
			reqPayload.JSONBody["model"] = reqPayload.Model
		}
		bodyModified = true
	}
	originalModel := reqPayload.Model

	isCodexCLI := openai.IsCodexCLIRequest(c.GetHeader("User-Agent"))
	requestBodyPassthrough := account.IsRequestBodyPassthroughEnabled()

	if IsOpenAIImageEndpoint(endpoint) && account.Type == AccountTypeOAuth {
		parsedImagesReq, parseErr := ParseOpenAIImagesRequest(body, reqPayload.ContentType, endpoint)
		if parseErr != nil {
			return nil, parseErr
		}
		return s.forwardOpenAIImagesOAuth(ctx, c, account, parsedImagesReq, account.GetMappedModel(parsedImagesReq.Model))
	}

	// 对所有请求执行模型映射（包含 Codex CLI）。
	mappedModel := account.GetMappedModel(reqPayload.Model)
	if mappedModel != reqPayload.Model && !requestBodyPassthrough {
		log.Printf("[OpenAI] Model mapping applied: %s -> %s (account: %s, isCodexCLI: %v)", reqPayload.Model, mappedModel, account.Name, isCodexCLI)
		if reqPayload.JSONBody != nil {
			reqPayload.JSONBody["model"] = mappedModel
		}
		bodyModified = true
	}

	// OAuth 账号走 ChatGPT internal API，需要把客户端别名规范化成内部 Codex 模型名。
	// API Key 账号走 OpenAI Platform API，应保留用户请求/显式映射后的模型名，避免把
	// 平台模型（如 gpt-5-codex）改写成另一个内部模型后触发上游 404。
	if currentModel := firstNonEmptyModel(reqPayload, mappedModel); account.Type == AccountTypeOAuth && currentModel != "" && !requestBodyPassthrough {
		normalizedModel := normalizeCodexModel(currentModel)
		if normalizedModel != "" && normalizedModel != currentModel {
			log.Printf("[OpenAI] Codex model normalization: %s -> %s (account: %s, type: %s, isCodexCLI: %v)",
				currentModel, normalizedModel, account.Name, account.Type, isCodexCLI)
			if reqPayload.JSONBody != nil {
				reqPayload.JSONBody["model"] = normalizedModel
			}
			mappedModel = normalizedModel
			bodyModified = true
		}
	}

	// 规范化 reasoning.effort 参数（minimal -> none），与上游允许值对齐。
	if reasoning, ok := reqPayload.JSONBody["reasoning"].(map[string]any); ok && !requestBodyPassthrough {
		if effort, ok := reasoning["effort"].(string); ok && effort == "minimal" {
			reasoning["effort"] = "none"
			bodyModified = true
			log.Printf("[OpenAI] Normalized reasoning.effort: minimal -> none (account: %s)", account.Name)
		}
	}

	if upstreamEndpoint == OpenAIEndpointResponses && account.Type == AccountTypeOAuth && !isCodexCLI && !requestBodyPassthrough {
		codexResult := applyCodexOAuthTransform(reqPayload.JSONBody)
		if codexResult.Modified {
			bodyModified = true
		}
		if codexResult.NormalizedModel != "" {
			mappedModel = codexResult.NormalizedModel
		}
		if codexResult.PromptCacheKey != "" {
			reqPayload.PromptCacheKey = codexResult.PromptCacheKey
		}
	}

	// Handle max_output_tokens based on platform and account type
	if reqPayload.JSONBody != nil && !isCodexCLI && !requestBodyPassthrough {
		if maxOutputTokens, hasMaxOutputTokens := reqPayload.JSONBody["max_output_tokens"]; hasMaxOutputTokens {
			switch account.Platform {
			case PlatformOpenAI:
				// For OpenAI API Key, remove max_output_tokens (not supported)
				// For OpenAI OAuth (Responses API), keep it (supported)
				if account.Type == AccountTypeAPIKey && upstreamEndpoint != OpenAIEndpointResponses && !IsOpenAIImageEndpoint(upstreamEndpoint) {
					delete(reqPayload.JSONBody, "max_output_tokens")
					bodyModified = true
				}
			case PlatformAnthropic:
				// For Anthropic (Claude), convert to max_tokens
				delete(reqPayload.JSONBody, "max_output_tokens")
				if _, hasMaxTokens := reqPayload.JSONBody["max_tokens"]; !hasMaxTokens {
					reqPayload.JSONBody["max_tokens"] = maxOutputTokens
				}
				bodyModified = true
			case PlatformGemini:
				// For Gemini, remove (will be handled by Gemini-specific transform)
				delete(reqPayload.JSONBody, "max_output_tokens")
				bodyModified = true
			default:
				// For unknown platforms, remove to be safe
				delete(reqPayload.JSONBody, "max_output_tokens")
				bodyModified = true
			}
		}

		// Also handle max_completion_tokens (similar logic)
		if _, hasMaxCompletionTokens := reqPayload.JSONBody["max_completion_tokens"]; hasMaxCompletionTokens {
			if account.Type == AccountTypeAPIKey || account.Platform != PlatformOpenAI {
				delete(reqPayload.JSONBody, "max_completion_tokens")
				bodyModified = true
			}
		}
	}

	// Re-serialize body only if modified
	if bodyModified {
		if reqPayload.MultipartBody != nil {
			updatedBody, updatedContentType, rebuildErr := reqPayload.MultipartBody.rebuildBody(map[string]string{
				"model": mappedModel,
			})
			if rebuildErr != nil {
				return nil, rebuildErr
			}
			body = updatedBody
			reqPayload.ContentType = updatedContentType
		} else {
			body, err = json.Marshal(reqPayload.JSONBody)
			if err != nil {
				return nil, fmt.Errorf("serialize request body: %w", err)
			}
		}
	}

	// Get access token
	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		return nil, err
	}

	// Build upstream request
	upstreamReq, err := s.buildUpstreamRequest(ctx, c, account, body, token, reqPayload.Stream, reqPayload.PromptCacheKey, isCodexCLI, upstreamEndpoint, reqPayload.ContentType)
	if err != nil {
		return nil, err
	}

	// Get proxy URL
	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}

	// Capture upstream request body for ops retry of this attempt.
	if c != nil {
		c.Set(OpsUpstreamRequestBodyKey, string(body))
	}

	// Send request
	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		safeErr := sanitizeUpstreamErrorMessage(err.Error())
		setOpsUpstreamError(c, 0, safeErr, "")
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: 0,
			Kind:               "request_error",
			Message:            safeErr,
		})
		return nil, &UpstreamFailoverError{StatusCode: 0}
	}
	defer func() { _ = resp.Body.Close() }()

	// Handle error response
	if resp.StatusCode >= 400 {
		if s.shouldFailoverUpstreamError(resp.StatusCode) {
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
			_ = resp.Body.Close()
			resp.Body = io.NopCloser(bytes.NewReader(respBody))

			upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
			upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
			upstreamDetail := ""
			if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
				maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
				if maxBytes <= 0 {
					maxBytes = 2048
				}
				upstreamDetail = truncateString(string(respBody), maxBytes)
			}
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: resp.StatusCode,
				UpstreamRequestID:  resp.Header.Get("x-request-id"),
				Kind:               "failover",
				Message:            upstreamMsg,
				Detail:             upstreamDetail,
			})

			s.handleFailoverSideEffects(ctx, resp, account)
			return nil, &UpstreamFailoverError{StatusCode: resp.StatusCode}
		}
		return s.handleErrorResponse(ctx, resp, c, account)
	}

	// Handle normal response
	var usage *OpenAIUsage
	var firstTokenMs *int
	imageCount := 0
	imageSize := reqPayload.ImageSize
	if reqPayload.Stream {
		streamResult, err := s.handleStreamingResponse(ctx, resp, c, account, startTime, originalModel, mappedModel, endpoint)
		if err != nil {
			return nil, err
		}
		usage = streamResult.usage
		firstTokenMs = streamResult.firstTokenMs
		imageCount = streamResult.imageCount
		if streamResult.imageSize != "" {
			imageSize = streamResult.imageSize
		}
	} else {
		var responseMeta *openAIImageResponseMeta
		responseMeta, err = s.handleNonStreamingResponse(ctx, resp, c, account, originalModel, mappedModel, endpoint)
		if err != nil {
			return nil, err
		}
		usage = &responseMeta.Usage
		imageCount = responseMeta.ImageCount
		if responseMeta.ImageSize != "" {
			imageSize = responseMeta.ImageSize
		}
	}

	// Extract and save Codex usage snapshot from response headers (for OAuth accounts)
	if account.Type == AccountTypeOAuth {
		if snapshot := ParseCodexRateLimitHeaders(resp.Header); snapshot != nil {
			s.updateCodexUsageSnapshot(ctx, account.ID, snapshot)
		}
	}

	return &OpenAIForwardResult{
		RequestID:    resp.Header.Get("x-request-id"),
		Usage:        *usage,
		Model:        originalModel,
		Stream:       reqPayload.Stream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
		ImageCount:   imageCount,
		ImageSize:    imageSize,
	}, nil
}

func firstNonEmptyModel(reqPayload *openAIRequestPayload, fallback string) string {
	if strings.TrimSpace(fallback) != "" {
		return strings.TrimSpace(fallback)
	}
	if reqPayload != nil && strings.TrimSpace(reqPayload.Model) != "" {
		return strings.TrimSpace(reqPayload.Model)
	}
	return ""
}

func (s *OpenAIGatewayService) buildUpstreamRequest(ctx context.Context, c *gin.Context, account *Account, body []byte, token string, isStream bool, promptCacheKey string, isCodexCLI bool, endpoint string, contentType string) (*http.Request, error) {
	// Determine target URL based on account type
	var targetURL string
	switch account.Type {
	case AccountTypeOAuth:
		// OAuth accounts use ChatGPT internal API
		if endpoint != OpenAIEndpointResponses {
			return nil, fmt.Errorf("endpoint %s is not supported for OAuth accounts", endpoint)
		}
		targetURL = chatgptCodexURL
	case AccountTypeAPIKey:
		// API Key accounts use Platform API or custom base URL
		baseURL := account.GetOpenAIBaseURL()
		if baseURL == "" {
			targetURL = openaiPlatformBaseURL + "/v1" + endpoint
		} else {
			validatedURL, err := s.validateUpstreamBaseURL(baseURL)
			if err != nil {
				return nil, err
			}
			targetURL = strings.TrimRight(validatedURL, "/") + "/v1" + endpoint
		}
	default:
		targetURL = openaiPlatformBaseURL + "/v1" + endpoint
	}

	req, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Set authentication header
	req.Header.Set("authorization", "Bearer "+token)

	// Set headers specific to OAuth accounts (ChatGPT internal API)
	if account.Type == AccountTypeOAuth {
		// Required: set Host for ChatGPT API (must use req.Host, not Header.Set)
		req.Host = "chatgpt.com"
		// Required: set chatgpt-account-id header
		chatgptAccountID := account.GetChatGPTAccountID()
		if chatgptAccountID != "" {
			req.Header.Set("chatgpt-account-id", chatgptAccountID)
		}
	}

	// Whitelist passthrough headers
	for key, values := range c.Request.Header {
		lowerKey := strings.ToLower(key)
		if openaiAllowedHeaders[lowerKey] {
			for _, v := range values {
				req.Header.Add(key, v)
			}
		}
	}
	if account.Type == AccountTypeOAuth {
		req.Header.Set("OpenAI-Beta", "responses=experimental")
		if isCodexCLI {
			req.Header.Set("originator", "codex_cli_rs")
		} else {
			req.Header.Set("originator", "opencode")
		}
		req.Header.Set("accept", "text/event-stream")
		if promptCacheKey != "" {
			req.Header.Set("conversation_id", promptCacheKey)
			req.Header.Set("session_id", promptCacheKey)
		}
	}

	// Apply custom User-Agent if configured
	customUA := account.GetOpenAIUserAgent()
	if customUA != "" {
		req.Header.Set("user-agent", customUA)
	}

	// Ensure required headers exist
	if strings.TrimSpace(contentType) != "" {
		req.Header.Set("content-type", contentType)
	}
	if req.Header.Get("content-type") == "" {
		req.Header.Set("content-type", "application/json")
	}

	return req, nil
}

func (s *OpenAIGatewayService) handleErrorResponse(ctx context.Context, resp *http.Response, c *gin.Context, account *Account) (*OpenAIForwardResult, error) {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))

	upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(body))
	upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
	upstreamDetail := ""
	if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
		maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
		if maxBytes <= 0 {
			maxBytes = 2048
		}
		upstreamDetail = truncateString(string(body), maxBytes)
	}
	setOpsUpstreamError(c, resp.StatusCode, upstreamMsg, upstreamDetail)

	if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
		log.Printf(
			"OpenAI upstream error %d (account=%d platform=%s type=%s): %s",
			resp.StatusCode,
			account.ID,
			account.Platform,
			account.Type,
			truncateForLog(body, s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes),
		)
	}

	// Check custom error codes
	if !account.ShouldHandleErrorCode(resp.StatusCode) {
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: resp.StatusCode,
			UpstreamRequestID:  resp.Header.Get("x-request-id"),
			Kind:               "http_error",
			Message:            upstreamMsg,
			Detail:             upstreamDetail,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"type":    "upstream_error",
				"message": "Upstream gateway error",
			},
		})
		if upstreamMsg == "" {
			return nil, fmt.Errorf("upstream error: %d (not in custom error codes)", resp.StatusCode)
		}
		return nil, fmt.Errorf("upstream error: %d (not in custom error codes) message=%s", resp.StatusCode, upstreamMsg)
	}

	// Handle upstream error (mark account status)
	shouldDisable := false
	if s.rateLimitService != nil {
		shouldDisable = s.rateLimitService.HandleUpstreamError(ctx, account, resp.StatusCode, resp.Header, body)
	}
	kind := "http_error"
	if shouldDisable {
		kind = "failover"
	}
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:           account.Platform,
		AccountID:          account.ID,
		AccountName:        account.Name,
		UpstreamStatusCode: resp.StatusCode,
		UpstreamRequestID:  resp.Header.Get("x-request-id"),
		Kind:               kind,
		Message:            upstreamMsg,
		Detail:             upstreamDetail,
	})
	if shouldDisable {
		return nil, &UpstreamFailoverError{StatusCode: resp.StatusCode}
	}

	// Return appropriate error response
	var errType, errMsg string
	var statusCode int

	switch resp.StatusCode {
	case 401:
		statusCode = http.StatusBadGateway
		errType = "upstream_error"
		errMsg = "Upstream authentication failed, please contact administrator"
	case 402:
		statusCode = http.StatusBadGateway
		errType = "upstream_error"
		errMsg = "Upstream payment required: insufficient balance or billing issue"
	case 403:
		statusCode = http.StatusBadGateway
		errType = "upstream_error"
		errMsg = "Upstream access forbidden, please contact administrator"
	case 429:
		statusCode = http.StatusTooManyRequests
		errType = "rate_limit_error"
		errMsg = "Upstream rate limit exceeded, please retry later"
	default:
		statusCode = http.StatusBadGateway
		errType = "upstream_error"
		errMsg = "Upstream request failed"
	}

	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"type":    errType,
			"message": errMsg,
		},
	})

	if upstreamMsg == "" {
		return nil, fmt.Errorf("upstream error: %d", resp.StatusCode)
	}
	return nil, fmt.Errorf("upstream error: %d message=%s", resp.StatusCode, upstreamMsg)
}

// openaiStreamingResult streaming response result
type openaiStreamingResult struct {
	usage        *OpenAIUsage
	firstTokenMs *int
	imageCount   int
	imageSize    string
}

func (s *OpenAIGatewayService) handleStreamingResponse(ctx context.Context, resp *http.Response, c *gin.Context, account *Account, startTime time.Time, originalModel, mappedModel, endpoint string) (*openaiStreamingResult, error) {
	if isOpenAIChatCompletionsMode(c) {
		return s.handleChatCompletionsStreamingResponse(ctx, resp, c, account, startTime, originalModel, mappedModel)
	}

	if s.cfg != nil {
		responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.cfg.Security.ResponseHeaders)
	}

	// Set SSE response headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// Pass through other headers
	if v := resp.Header.Get("x-request-id"); v != "" {
		c.Header("x-request-id", v)
	}

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming not supported")
	}

	usage := &OpenAIUsage{}
	imageMeta := &openAIImageResponseMeta{}
	var firstTokenMs *int
	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.cfg != nil && s.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.cfg.Gateway.MaxLineSize
	}
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)

	type scanEvent struct {
		line string
		err  error
	}
	// 独立 goroutine 读取上游，避免读取阻塞影响 keepalive/超时处理
	events := make(chan scanEvent, 16)
	done := make(chan struct{})
	sendEvent := func(ev scanEvent) bool {
		select {
		case events <- ev:
			return true
		case <-done:
			return false
		}
	}
	var lastReadAt int64
	atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
	go func() {
		defer close(events)
		for scanner.Scan() {
			atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
			if !sendEvent(scanEvent{line: scanner.Text()}) {
				return
			}
		}
		if err := scanner.Err(); err != nil {
			_ = sendEvent(scanEvent{err: err})
		}
	}()
	defer close(done)

	streamInterval := time.Duration(0)
	if s.cfg != nil && s.cfg.Gateway.StreamDataIntervalTimeout > 0 {
		streamInterval = time.Duration(s.cfg.Gateway.StreamDataIntervalTimeout) * time.Second
	}
	// 仅监控上游数据间隔超时，不被下游写入阻塞影响
	var intervalTicker *time.Ticker
	if streamInterval > 0 {
		intervalTicker = time.NewTicker(streamInterval)
		defer intervalTicker.Stop()
	}
	var intervalCh <-chan time.Time
	if intervalTicker != nil {
		intervalCh = intervalTicker.C
	}

	keepaliveInterval := time.Duration(0)
	if s.cfg != nil && s.cfg.Gateway.StreamKeepaliveInterval > 0 {
		keepaliveInterval = time.Duration(s.cfg.Gateway.StreamKeepaliveInterval) * time.Second
	}
	// 下游 keepalive 仅用于防止代理空闲断开
	var keepaliveTicker *time.Ticker
	if keepaliveInterval > 0 {
		keepaliveTicker = time.NewTicker(keepaliveInterval)
		defer keepaliveTicker.Stop()
	}
	var keepaliveCh <-chan time.Time
	if keepaliveTicker != nil {
		keepaliveCh = keepaliveTicker.C
	}
	// 记录上次收到上游数据的时间，用于控制 keepalive 发送频率
	lastDataAt := time.Now()

	// 仅发送一次错误事件，避免多次写入导致协议混乱（写失败时尽力通知客户端）
	errorEventSent := false
	sendErrorEvent := func(reason string) {
		if errorEventSent {
			return
		}
		errorEventSent = true
		_, _ = fmt.Fprintf(w, "event: error\ndata: {\"error\":\"%s\"}\n\n", reason)
		flusher.Flush()
	}

	needModelReplace := originalModel != mappedModel

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs, imageCount: imageMeta.ImageCount, imageSize: imageMeta.ImageSize}, nil
			}
			if ev.err != nil {
				if errors.Is(ev.err, bufio.ErrTooLong) {
					log.Printf("SSE line too long: account=%d max_size=%d error=%v", account.ID, maxLineSize, ev.err)
					sendErrorEvent("response_too_large")
					return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs, imageCount: imageMeta.ImageCount, imageSize: imageMeta.ImageSize}, ev.err
				}
				sendErrorEvent("stream_read_error")
				return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs, imageCount: imageMeta.ImageCount, imageSize: imageMeta.ImageSize}, fmt.Errorf("stream read error: %w", ev.err)
			}

			line := ev.line
			lastDataAt = time.Now()

			// Extract data from SSE line (supports both "data: " and "data:" formats)
			if openaiSSEDataRe.MatchString(line) {
				data := openaiSSEDataRe.ReplaceAllString(line, "")

				// Replace model in response if needed
				if needModelReplace {
					line = s.replaceModelInSSELine(line, mappedModel, originalModel)
				}

				// Correct Codex tool calls if needed (apply_patch -> edit, etc.)
				if correctedData, corrected := s.toolCorrector.CorrectToolCallsInSSEData(data); corrected {
					line = "data: " + correctedData
				}

				// Forward line
				if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
					sendErrorEvent("write_failed")
					return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs, imageCount: imageMeta.ImageCount, imageSize: imageMeta.ImageSize}, err
				}
				flusher.Flush()

				// Record first token time
				if firstTokenMs == nil && data != "" && data != "[DONE]" {
					ms := int(time.Since(startTime).Milliseconds())
					firstTokenMs = &ms
				}
				s.parseSSEUsage(data, usage)
				if IsOpenAIImageEndpoint(endpoint) {
					s.parseOpenAIImageSSEEvent(data, imageMeta)
				}
			} else {
				// Forward non-data lines as-is
				if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
					sendErrorEvent("write_failed")
					return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs, imageCount: imageMeta.ImageCount, imageSize: imageMeta.ImageSize}, err
				}
				flusher.Flush()
			}

		case <-intervalCh:
			lastRead := time.Unix(0, atomic.LoadInt64(&lastReadAt))
			if time.Since(lastRead) < streamInterval {
				continue
			}
			log.Printf("Stream data interval timeout: account=%d model=%s interval=%s", account.ID, originalModel, streamInterval)
			// 处理流超时，可能标记账户为临时不可调度或错误状态
			if s.rateLimitService != nil {
				s.rateLimitService.HandleStreamTimeout(ctx, account, originalModel)
			}
			sendErrorEvent("stream_timeout")
			return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs, imageCount: imageMeta.ImageCount, imageSize: imageMeta.ImageSize}, fmt.Errorf("stream data interval timeout")

		case <-keepaliveCh:
			if time.Since(lastDataAt) < keepaliveInterval {
				continue
			}
			if _, err := fmt.Fprint(w, ":\n\n"); err != nil {
				return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMs, imageCount: imageMeta.ImageCount, imageSize: imageMeta.ImageSize}, err
			}
			flusher.Flush()
		}
	}

}

func (s *OpenAIGatewayService) replaceModelInSSELine(line, fromModel, toModel string) string {
	if !openaiSSEDataRe.MatchString(line) {
		return line
	}
	data := openaiSSEDataRe.ReplaceAllString(line, "")
	if data == "" || data == "[DONE]" {
		return line
	}

	var event map[string]any
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return line
	}

	// Replace model in response
	if m, ok := event["model"].(string); ok && m == fromModel {
		event["model"] = toModel
		newData, err := json.Marshal(event)
		if err != nil {
			return line
		}
		return "data: " + string(newData)
	}

	// Check nested response
	if response, ok := event["response"].(map[string]any); ok {
		if m, ok := response["model"].(string); ok && m == fromModel {
			response["model"] = toModel
			newData, err := json.Marshal(event)
			if err != nil {
				return line
			}
			return "data: " + string(newData)
		}
	}

	return line
}

// correctToolCallsInResponseBody 修正响应体中的工具调用
func (s *OpenAIGatewayService) correctToolCallsInResponseBody(body []byte) []byte {
	if len(body) == 0 {
		return body
	}

	bodyStr := string(body)
	corrected, changed := s.toolCorrector.CorrectToolCallsInSSEData(bodyStr)
	if changed {
		return []byte(corrected)
	}
	return body
}

func (s *OpenAIGatewayService) parseSSEUsage(data string, usage *OpenAIUsage) {
	// Parse response.completed event for usage (OpenAI Responses format)
	var event struct {
		Type     string `json:"type"`
		Response struct {
			Usage struct {
				InputTokens       int `json:"input_tokens"`
				OutputTokens      int `json:"output_tokens"`
				InputTokenDetails struct {
					CachedTokens int `json:"cached_tokens"`
				} `json:"input_tokens_details"`
			} `json:"usage"`
		} `json:"response"`
	}

	if json.Unmarshal([]byte(data), &event) == nil && event.Type == "response.completed" {
		usage.InputTokens = event.Response.Usage.InputTokens
		usage.OutputTokens = event.Response.Usage.OutputTokens
		usage.CacheReadInputTokens = event.Response.Usage.InputTokenDetails.CachedTokens
	}
}

func (s *OpenAIGatewayService) handleNonStreamingResponse(ctx context.Context, resp *http.Response, c *gin.Context, account *Account, originalModel, mappedModel, endpoint string) (*openAIImageResponseMeta, error) {
	if isOpenAIChatCompletionsMode(c) {
		usage, err := s.handleChatCompletionsNonStreamingResponse(ctx, resp, c, account, originalModel, mappedModel)
		if err != nil {
			return nil, err
		}
		return &openAIImageResponseMeta{Usage: *usage}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if account.Type == AccountTypeOAuth {
		bodyLooksLikeSSE := bytes.Contains(body, []byte("data:")) || bytes.Contains(body, []byte("event:"))
		if isEventStreamResponse(resp.Header) || bodyLooksLikeSSE {
			usage, oauthErr := s.handleOAuthSSEToJSON(resp, c, body, originalModel, mappedModel)
			if oauthErr != nil {
				return nil, oauthErr
			}
			return &openAIImageResponseMeta{Usage: *usage}, nil
		}
	}

	// Replace model in response if needed
	if originalModel != mappedModel {
		body = s.replaceModelInResponseBody(body, mappedModel, originalModel)
	}

	meta, err := s.parseOpenAINonStreamingResponse(body, endpoint)
	if err != nil {
		return nil, err
	}

	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.cfg.Security.ResponseHeaders)

	contentType := "application/json"
	if s.cfg != nil && !s.cfg.Security.ResponseHeaders.Enabled {
		if upstreamType := resp.Header.Get("Content-Type"); upstreamType != "" {
			contentType = upstreamType
		}
	}

	c.Data(resp.StatusCode, contentType, body)

	return meta, nil
}

func isEventStreamResponse(header http.Header) bool {
	contentType := strings.ToLower(header.Get("Content-Type"))
	return strings.Contains(contentType, "text/event-stream")
}

func (s *OpenAIGatewayService) parseOpenAINonStreamingResponse(body []byte, endpoint string) (*openAIImageResponseMeta, error) {
	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	meta := &openAIImageResponseMeta{
		Usage: OpenAIUsage{
			InputTokens:          intValue(nestedValue(response, "usage", "input_tokens")),
			OutputTokens:         intValue(nestedValue(response, "usage", "output_tokens")),
			CacheReadInputTokens: intValue(nestedValue(response, "usage", "input_tokens_details", "cached_tokens")),
		},
	}

	if !IsOpenAIImageEndpoint(endpoint) {
		return meta, nil
	}

	if data, ok := response["data"].([]any); ok {
		meta.ImageCount = len(data)
		if size := stringValue(response["size"]); size != "" {
			meta.ImageSize = normalizeOpenAIImageSize(size)
		}
		if meta.ImageSize == "" {
			for _, rawItem := range data {
				item, ok := rawItem.(map[string]any)
				if !ok {
					continue
				}
				if size := stringValue(item["size"]); size != "" {
					meta.ImageSize = normalizeOpenAIImageSize(size)
					break
				}
			}
		}
	}

	return meta, nil
}

func (s *OpenAIGatewayService) parseOpenAIImageSSEEvent(data string, meta *openAIImageResponseMeta) {
	if meta == nil || data == "" || data == "[DONE]" {
		return
	}

	var event map[string]any
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return
	}

	eventType := strings.TrimSpace(stringValue(event["type"]))
	switch eventType {
	case "image_generation.completed", "image_edit.completed":
		meta.ImageCount++
		if size := normalizeOpenAIImageSize(stringValue(event["size"])); size != "" {
			meta.ImageSize = size
		}
		if usage, ok := event["usage"].(map[string]any); ok {
			if meta.Usage.InputTokens == 0 {
				meta.Usage.InputTokens = intValue(usage["input_tokens"])
			}
			meta.Usage.OutputTokens += intValue(usage["output_tokens"])
			if details, ok := usage["input_tokens_details"].(map[string]any); ok {
				meta.Usage.CacheReadInputTokens += intValue(details["cached_tokens"])
			}
		}
	}
}

func normalizeOpenAIImageSize(size string) string {
	normalized := strings.ToLower(strings.TrimSpace(size))
	if normalized == "" || normalized == "auto" {
		return ""
	}

	switch normalized {
	case "1024x1024", "1024x-1024", "1024-x-1024":
		return "1K"
	case "1536x1024", "1024x1536", "1536x-1024", "1024x-1536", "1536-x-1024", "1024-x-1536":
		return "2K"
	case "1792x1024", "1024x1792", "1792x-1024", "1024x-1792", "1792-x-1024", "1024-x-1792":
		return "4K"
	default:
		return ""
	}
}

func nestedValue(value any, path ...string) any {
	current := value
	for _, key := range path {
		obj, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = obj[key]
	}
	return current
}

func (s *OpenAIGatewayService) handleOAuthSSEToJSON(resp *http.Response, c *gin.Context, body []byte, originalModel, mappedModel string) (*OpenAIUsage, error) {
	bodyText := string(body)
	finalResponse, ok := extractCodexFinalResponse(bodyText)

	usage := &OpenAIUsage{}
	if ok {
		var response struct {
			Usage struct {
				InputTokens       int `json:"input_tokens"`
				OutputTokens      int `json:"output_tokens"`
				InputTokenDetails struct {
					CachedTokens int `json:"cached_tokens"`
				} `json:"input_tokens_details"`
			} `json:"usage"`
		}
		if err := json.Unmarshal(finalResponse, &response); err == nil {
			usage.InputTokens = response.Usage.InputTokens
			usage.OutputTokens = response.Usage.OutputTokens
			usage.CacheReadInputTokens = response.Usage.InputTokenDetails.CachedTokens
		}
		body = finalResponse
		if originalModel != mappedModel {
			body = s.replaceModelInResponseBody(body, mappedModel, originalModel)
		}
		// Correct tool calls in final response
		body = s.correctToolCallsInResponseBody(body)
	} else {
		usage = s.parseSSEUsageFromBody(bodyText)
		if originalModel != mappedModel {
			bodyText = s.replaceModelInSSEBody(bodyText, mappedModel, originalModel)
		}
		body = []byte(bodyText)
	}

	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.cfg.Security.ResponseHeaders)

	contentType := "application/json; charset=utf-8"
	if !ok {
		contentType = resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "text/event-stream"
		}
	}
	c.Data(resp.StatusCode, contentType, body)

	return usage, nil
}

func extractCodexFinalResponse(body string) ([]byte, bool) {
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if !openaiSSEDataRe.MatchString(line) {
			continue
		}
		data := openaiSSEDataRe.ReplaceAllString(line, "")
		if data == "" || data == "[DONE]" {
			continue
		}
		var event struct {
			Type     string          `json:"type"`
			Response json.RawMessage `json:"response"`
		}
		if json.Unmarshal([]byte(data), &event) != nil {
			continue
		}
		if event.Type == "response.done" || event.Type == "response.completed" {
			if len(event.Response) > 0 {
				return event.Response, true
			}
		}
	}
	return nil, false
}

func (s *OpenAIGatewayService) parseSSEUsageFromBody(body string) *OpenAIUsage {
	usage := &OpenAIUsage{}
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if !openaiSSEDataRe.MatchString(line) {
			continue
		}
		data := openaiSSEDataRe.ReplaceAllString(line, "")
		if data == "" || data == "[DONE]" {
			continue
		}
		s.parseSSEUsage(data, usage)
	}
	return usage
}

func (s *OpenAIGatewayService) replaceModelInSSEBody(body, fromModel, toModel string) string {
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if !openaiSSEDataRe.MatchString(line) {
			continue
		}
		lines[i] = s.replaceModelInSSELine(line, fromModel, toModel)
	}
	return strings.Join(lines, "\n")
}

func (s *OpenAIGatewayService) validateUpstreamBaseURL(raw string) (string, error) {
	if s.cfg == nil {
		normalized, err := urlvalidator.ValidateURLFormat(raw, false)
		if err != nil {
			return "", fmt.Errorf("invalid base_url: %w", err)
		}
		return normalized, nil
	}
	if s.cfg != nil && !s.cfg.Security.URLAllowlist.Enabled {
		normalized, err := urlvalidator.ValidateURLFormat(raw, s.cfg.Security.URLAllowlist.AllowInsecureHTTP)
		if err != nil {
			return "", fmt.Errorf("invalid base_url: %w", err)
		}
		return normalized, nil
	}
	normalized, err := urlvalidator.ValidateHTTPSURL(raw, urlvalidator.ValidationOptions{
		AllowedHosts:     s.cfg.Security.URLAllowlist.UpstreamHosts,
		RequireAllowlist: true,
		AllowPrivate:     s.cfg.Security.URLAllowlist.AllowPrivateHosts,
	})
	if err != nil {
		return "", fmt.Errorf("invalid base_url: %w", err)
	}
	return normalized, nil
}

func (s *OpenAIGatewayService) replaceModelInResponseBody(body []byte, fromModel, toModel string) []byte {
	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return body
	}

	model, ok := resp["model"].(string)
	if !ok || model != fromModel {
		return body
	}

	resp["model"] = toModel
	newBody, err := json.Marshal(resp)
	if err != nil {
		return body
	}

	return newBody
}

// OpenAIRecordUsageInput input for recording usage
type OpenAIRecordUsageInput struct {
	Result       *OpenAIForwardResult
	APIKey       *APIKey
	User         *User
	Account      *Account
	Subscription *UserSubscription
	UserAgent    string // 请求的 User-Agent
	IPAddress    string // 请求的客户端 IP 地址
	// 企业计费字段（可选）
	Organization *Organization
	OrgMember    *OrgMember
	OrgProject   *OrgProject
	RequestBody  string // 请求体（用于审计日志）
}

// RecordUsage records usage and deducts balance
func (s *OpenAIGatewayService) RecordUsage(ctx context.Context, input *OpenAIRecordUsageInput) error {
	result := input.Result
	apiKey := input.APIKey
	user := input.User
	account := input.Account
	subscription := input.Subscription

	// 计算实际的新输入token（减去缓存读取的token）
	// 因为 input_tokens 包含了 cache_read_tokens，而缓存读取的token不应按输入价格计费
	actualInputTokens := result.Usage.InputTokens - result.Usage.CacheReadInputTokens
	if actualInputTokens < 0 {
		actualInputTokens = 0
	}

	// Get rate multiplier
	multiplier := s.cfg.Default.RateMultiplier
	if apiKey.GroupID != nil && apiKey.Group != nil {
		multiplier = apiKey.Group.RateMultiplier
	}
	subSiteRate := currentSubSiteConsumeRateMultiplier(ctx)
	var boundSubSite *SubSite
	var subSiteChain []*SubSite
	if s.subSiteService != nil {
		if site, err := s.subSiteService.GetBoundSubSiteForUser(ctx, user.ID); err == nil && site != nil && site.Status == SubSiteStatusActive {
			boundSubSite = site
			if chain, err := s.subSiteService.GetSiteChain(ctx, site.ID); err == nil && len(chain) > 0 {
				if s.subSiteService.ChainActive(chain) {
					subSiteChain = chain
					subSiteRate = s.subSiteService.CompoundConsumeRateForChain(chain)
				} else {
					boundSubSite = nil
					subSiteRate = DefaultSubSiteConsumeRate
				}
			} else if subSiteRate <= 0 || subSiteRate == DefaultSubSiteConsumeRate {
				subSiteRate = normalizeConsumeRateMultiplier(site.ConsumeRateMultiplier)
			}
		}
	}
	if subSiteRate > 0 && subSiteRate != DefaultSubSiteConsumeRate {
		multiplier *= subSiteRate
	}

	var cost *CostBreakdown
	if result.ImageCount > 0 {
		var groupConfig *ImagePriceConfig
		if apiKey.Group != nil {
			groupConfig = &ImagePriceConfig{
				Price1K: apiKey.Group.ImagePrice1K,
				Price2K: apiKey.Group.ImagePrice2K,
				Price4K: apiKey.Group.ImagePrice4K,
			}
		}
		cost = s.billingService.CalculateImageCost(result.Model, result.ImageSize, result.ImageCount, groupConfig, multiplier)
	} else {
		tokens := UsageTokens{
			InputTokens:         actualInputTokens,
			OutputTokens:        result.Usage.OutputTokens,
			CacheCreationTokens: result.Usage.CacheCreationInputTokens,
			CacheReadTokens:     result.Usage.CacheReadInputTokens,
		}
		var costErr error
		cost, costErr = s.billingService.CalculateCost(result.Model, tokens, multiplier)
		if costErr != nil {
			cost = &CostBreakdown{ActualCost: 0}
		}
	}

	// Determine billing type
	isSubscriptionBilling := subscription != nil && apiKey.Group != nil && apiKey.Group.IsSubscriptionType()
	isQuotaPackageBilling := subscription == nil && apiKey.GroupID != nil && apiKey.Group != nil && apiKey.Group.IsQuotaPackage()
	billingType := BillingTypeBalance
	switch {
	case isSubscriptionBilling:
		billingType = BillingTypeSubscription
	case isQuotaPackageBilling:
		billingType = BillingTypeQuotaPackage
	}

	// 应用账号计费倍率到真实扣费；展示层会按角色决定是否暴露成本明细。
	cost.ActualCost *= account.BillingRateMultiplier()

	// Create usage log
	durationMs := int(result.Duration.Milliseconds())
	accountRateMultiplier := account.BillingRateMultiplier()
	var imageSize *string
	if result.ImageSize != "" {
		imageSize = &result.ImageSize
	}
	usageLog := &UsageLog{
		UserID:                user.ID,
		APIKeyID:              apiKey.ID,
		AccountID:             account.ID,
		RequestID:             result.RequestID,
		Model:                 result.Model,
		InputTokens:           actualInputTokens,
		OutputTokens:          result.Usage.OutputTokens,
		CacheCreationTokens:   result.Usage.CacheCreationInputTokens,
		CacheReadTokens:       result.Usage.CacheReadInputTokens,
		InputCost:             cost.InputCost,
		OutputCost:            cost.OutputCost,
		CacheCreationCost:     cost.CacheCreationCost,
		CacheReadCost:         cost.CacheReadCost,
		TotalCost:             cost.TotalCost,
		ActualCost:            cost.ActualCost,
		RateMultiplier:        multiplier,
		AccountRateMultiplier: &accountRateMultiplier,
		BillingType:           billingType,
		Stream:                result.Stream,
		DurationMs:            &durationMs,
		FirstTokenMs:          result.FirstTokenMs,
		ImageCount:            result.ImageCount,
		ImageSize:             imageSize,
		CreatedAt:             time.Now(),
	}

	// 添加 UserAgent
	if input.UserAgent != "" {
		usageLog.UserAgent = &input.UserAgent
	}

	// 添加 IPAddress
	if input.IPAddress != "" {
		usageLog.IPAddress = &input.IPAddress
	}

	if apiKey.GroupID != nil {
		usageLog.GroupID = apiKey.GroupID
	}
	if subscription != nil {
		usageLog.SubscriptionID = &subscription.ID
	}
	// 添加企业关联
	if input.Organization != nil {
		usageLog.OrgID = &input.Organization.ID
	}
	if input.OrgMember != nil {
		usageLog.OrgMemberID = &input.OrgMember.ID
	}

	inserted, err := s.usageLogRepo.Create(ctx, usageLog)
	if s.cfg != nil && s.cfg.RunMode == config.RunModeSimple {
		log.Printf("[SIMPLE MODE] Usage recorded (not billed): user=%d, tokens=%d", usageLog.UserID, usageLog.TotalTokens())
		s.deferredService.ScheduleLastUsedUpdate(account.ID)
		return nil
	}

	shouldBill := inserted || err != nil

	// === 企业计费路径 ===
	isOrgBilling := input.Organization != nil && apiKey.OrgID != nil
	if isOrgBilling && shouldBill && cost.ActualCost > 0 {
		if input.Organization.BillingMode == OrgBillingModeBalance {
			if s.orgRepo != nil {
				_ = s.orgRepo.DeductBalance(ctx, input.Organization.ID, cost.ActualCost)
			}
		}
		if input.OrgMember != nil && s.orgMemberRepo != nil {
			_ = s.orgMemberRepo.IncrementUsage(ctx, input.OrgMember.ID, cost.ActualCost)
		}
		// 累加项目用量
		if input.OrgProject != nil && s.orgProjectRepo != nil {
			_ = s.orgProjectRepo.IncrementUsage(ctx, input.OrgProject.ID, cost.ActualCost)
		}
		// 写入审计日志
		if s.auditService != nil {
			var projectID *int64
			if input.OrgProject != nil {
				projectID = &input.OrgProject.ID
			}
			var memberID *int64
			if input.OrgMember != nil {
				memberID = &input.OrgMember.ID
			}
			model := result.Model
			inputTokens := result.Usage.InputTokens
			outputTokens := result.Usage.OutputTokens
			go func() {
				_ = s.auditService.WriteAuditLog(context.Background(), &WriteAuditInput{
					OrgID:        input.Organization.ID,
					UserID:       user.ID,
					MemberID:     memberID,
					ProjectID:    projectID,
					UsageLogID:   &usageLog.ID,
					Action:       AuditActionAPIRequest,
					Model:        &model,
					AuditMode:    input.Organization.AuditMode,
					InputTokens:  &inputTokens,
					OutputTokens: &outputTokens,
					CostUSD:      &cost.ActualCost,
					IPAddress:    &input.IPAddress,
					UserAgent:    &input.UserAgent,
					RequestBody:  input.RequestBody,
				})
			}()
		}
		s.deferredService.ScheduleLastUsedUpdate(account.ID)
		return nil
	}

	// Deduct based on billing type
	if isQuotaPackageBilling {
		if shouldBill && cost.ActualCost > 0 {
			if s.quotaPackageRepo == nil {
				log.Printf("Quota package repository is unavailable: user=%d group=%d", user.ID, *apiKey.GroupID)
			} else if err := s.quotaPackageRepo.Deduct(ctx, user.ID, *apiKey.GroupID, cost.ActualCost); err != nil {
				log.Printf("Deduct quota package failed: %v", err)
			} else if s.wechatNotifyService != nil {
				if remaining, err := s.quotaPackageRepo.GetAvailableTotal(ctx, user.ID, *apiKey.GroupID); err == nil {
					go s.wechatNotifyService.NotifyQuotaAfterDeduct(context.Background(), user.ID, apiKey.Group, remaining)
				}
			}
		}
	} else if isSubscriptionBilling {
		if shouldBill && cost.ActualCost > 0 {
			if err := s.userSubRepo.IncrementUsage(ctx, subscription.ID, cost.ActualCost); err == nil && s.wechatNotifyService != nil {
				go s.wechatNotifyService.NotifySubscriptionAfterUsage(context.Background(), user.ID, subscription, apiKey.Group, cost.ActualCost)
			}
			s.billingCacheService.QueueUpdateSubscriptionUsage(user.ID, *apiKey.GroupID, cost.ActualCost)
		}
	} else {
		if shouldBill && cost.ActualCost > 0 {
			if err := s.userRepo.DeductBalance(ctx, user.ID, cost.ActualCost); err == nil && s.wechatNotifyService != nil {
				go s.wechatNotifyService.NotifyBalanceAfterDeduct(context.Background(), user, cost.ActualCost)
			}
			s.billingCacheService.QueueDeductBalance(user.ID, cost.ActualCost)
			if boundSubSite != nil && subSiteRate > 0 && s.subSiteService != nil && len(subSiteChain) > 0 {
				s.subSiteService.DebitPoolForConsumption(ctx, boundSubSite.ID, user.ID, usageLog.ID, cost.ActualCost, subSiteRate)
			}
		}
	}

	// Schedule batch update for account last_used_at
	s.deferredService.ScheduleLastUsedUpdate(account.ID)

	return nil
}

// ParseCodexRateLimitHeaders extracts Codex usage limits from response headers.
// Exported for use in ratelimit_service when handling OpenAI 429 responses.
func ParseCodexRateLimitHeaders(headers http.Header) *OpenAICodexUsageSnapshot {
	snapshot := &OpenAICodexUsageSnapshot{}
	hasData := false

	// Helper to parse float64 from header
	parseFloat := func(key string) *float64 {
		if v := headers.Get(key); v != "" {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return &f
			}
		}
		return nil
	}

	// Helper to parse int from header
	parseInt := func(key string) *int {
		if v := headers.Get(key); v != "" {
			if i, err := strconv.Atoi(v); err == nil {
				return &i
			}
		}
		return nil
	}

	// Primary (weekly) limits
	if v := parseFloat("x-codex-primary-used-percent"); v != nil {
		snapshot.PrimaryUsedPercent = v
		hasData = true
	}
	if v := parseInt("x-codex-primary-reset-after-seconds"); v != nil {
		snapshot.PrimaryResetAfterSeconds = v
		hasData = true
	}
	if v := parseInt("x-codex-primary-window-minutes"); v != nil {
		snapshot.PrimaryWindowMinutes = v
		hasData = true
	}

	// Secondary (5h) limits
	if v := parseFloat("x-codex-secondary-used-percent"); v != nil {
		snapshot.SecondaryUsedPercent = v
		hasData = true
	}
	if v := parseInt("x-codex-secondary-reset-after-seconds"); v != nil {
		snapshot.SecondaryResetAfterSeconds = v
		hasData = true
	}
	if v := parseInt("x-codex-secondary-window-minutes"); v != nil {
		snapshot.SecondaryWindowMinutes = v
		hasData = true
	}

	// Overflow ratio
	if v := parseFloat("x-codex-primary-over-secondary-limit-percent"); v != nil {
		snapshot.PrimaryOverSecondaryPercent = v
		hasData = true
	}

	if !hasData {
		return nil
	}

	snapshot.UpdatedAt = time.Now().Format(time.RFC3339)
	return snapshot
}

// updateCodexUsageSnapshot saves the Codex usage snapshot to account's Extra field
func (s *OpenAIGatewayService) updateCodexUsageSnapshot(ctx context.Context, accountID int64, snapshot *OpenAICodexUsageSnapshot) {
	if snapshot == nil {
		return
	}

	// Convert snapshot to map for merging into Extra
	updates := make(map[string]any)

	// Save raw primary/secondary fields for debugging/tracing
	if snapshot.PrimaryUsedPercent != nil {
		updates["codex_primary_used_percent"] = *snapshot.PrimaryUsedPercent
	}
	if snapshot.PrimaryResetAfterSeconds != nil {
		updates["codex_primary_reset_after_seconds"] = *snapshot.PrimaryResetAfterSeconds
	}
	if snapshot.PrimaryWindowMinutes != nil {
		updates["codex_primary_window_minutes"] = *snapshot.PrimaryWindowMinutes
	}
	if snapshot.SecondaryUsedPercent != nil {
		updates["codex_secondary_used_percent"] = *snapshot.SecondaryUsedPercent
	}
	if snapshot.SecondaryResetAfterSeconds != nil {
		updates["codex_secondary_reset_after_seconds"] = *snapshot.SecondaryResetAfterSeconds
	}
	if snapshot.SecondaryWindowMinutes != nil {
		updates["codex_secondary_window_minutes"] = *snapshot.SecondaryWindowMinutes
	}
	if snapshot.PrimaryOverSecondaryPercent != nil {
		updates["codex_primary_over_secondary_percent"] = *snapshot.PrimaryOverSecondaryPercent
	}
	updates["codex_usage_updated_at"] = snapshot.UpdatedAt

	// Normalize to canonical 5h/7d fields
	if normalized := snapshot.Normalize(); normalized != nil {
		if normalized.Used5hPercent != nil {
			updates["codex_5h_used_percent"] = *normalized.Used5hPercent
		}
		if normalized.Reset5hSeconds != nil {
			updates["codex_5h_reset_after_seconds"] = *normalized.Reset5hSeconds
		}
		if normalized.Window5hMinutes != nil {
			updates["codex_5h_window_minutes"] = *normalized.Window5hMinutes
		}
		if normalized.Used7dPercent != nil {
			updates["codex_7d_used_percent"] = *normalized.Used7dPercent
		}
		if normalized.Reset7dSeconds != nil {
			updates["codex_7d_reset_after_seconds"] = *normalized.Reset7dSeconds
		}
		if normalized.Window7dMinutes != nil {
			updates["codex_7d_window_minutes"] = *normalized.Window7dMinutes
		}
	}

	// Update account's Extra field asynchronously
	go func() {
		updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.accountRepo.UpdateExtra(updateCtx, accountID, updates)
	}()
}
