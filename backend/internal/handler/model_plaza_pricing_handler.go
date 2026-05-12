package handler

import (
	"sort"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

const pricingTablePerMillion = 1_000_000

type modelPlazaPricingTableResponse struct {
	Groups []modelPlazaPricingGroup `json:"groups"`
	Items  []modelPlazaPricingItem  `json:"items"`
}

type modelPlazaPricingGroup struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Platform          string `json:"platform"`
	DisplayPrice      string `json:"display_price"`
	DisplayDiscount   string `json:"display_discount"`
	Description       string `json:"description"`
	SubscriptionType  string `json:"subscription_type"`
	HasExplicitModels bool   `json:"has_explicit_models"`
	ModelsCount       int    `json:"models_count"`
}

type modelPlazaPricingItem struct {
	Model                 string                             `json:"model"`
	Aliases               []string                           `json:"aliases,omitempty"`
	Platform              string                             `json:"platform"`
	Provider              string                             `json:"provider"`
	Mode                  string                             `json:"mode"`
	SupportsPromptCaching bool                               `json:"supports_prompt_caching"`
	MaxInputTokens        int                                `json:"max_input_tokens"`
	MaxOutputTokens       int                                `json:"max_output_tokens"`
	Official              modelPlazaPricingMetrics           `json:"official"`
	GroupPrices           map[int64]modelPlazaPricingMetrics `json:"group_prices"`
}

type modelPlazaPricingMetrics struct {
	InputPerMillion               float64  `json:"input_per_million"`
	OutputPerMillion              float64  `json:"output_per_million"`
	CacheWritePerMillion          float64  `json:"cache_write_per_million"`
	CacheReadPerMillion           float64  `json:"cache_read_per_million"`
	InputPerMillionAbove200K      *float64 `json:"input_per_million_above_200k,omitempty"`
	OutputPerMillionAbove200K     *float64 `json:"output_per_million_above_200k,omitempty"`
	CacheWritePerMillionAbove200K *float64 `json:"cache_write_per_million_above_200k,omitempty"`
	CacheReadPerMillionAbove200K  *float64 `json:"cache_read_per_million_above_200k,omitempty"`
}

type pricingRowAccumulator struct {
	Model                 string
	Platform              string
	Provider              string
	Mode                  string
	SupportsPromptCaching bool
	MaxInputTokens        int
	MaxOutputTokens       int
	Official              modelPlazaPricingMetrics
	Aliases               map[string]struct{}
	GroupPrices           map[int64]modelPlazaPricingMetrics
}

// PricingTable returns official pricing rows plus derived group pricing rows.
// GET /api/v1/model-plaza/pricing-table
func (h *ModelPlazaHandler) PricingTable(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	groups, err := h.apiKeyService.GetAvailableGroups(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	groups = filterModelPlazaVisibleGroups(groups)

	groupIDs := make([]int64, 0, len(groups))
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID)
	}

	modelsByGroup := h.gatewayService.GetAvailableModelsByGroups(c.Request.Context(), groupIDs, "")
	items, groupSummaries := h.buildPricingTable(groups, modelsByGroup)

	response.Success(c, modelPlazaPricingTableResponse{
		Groups: groupSummaries,
		Items:  items,
	})
}

func (h *ModelPlazaHandler) buildPricingTable(groups []service.Group, modelsByGroup map[int64][]string) ([]modelPlazaPricingItem, []modelPlazaPricingGroup) {
	sortedGroups := append([]service.Group(nil), groups...)
	sort.Slice(sortedGroups, func(i, j int) bool {
		left := sortedGroups[i]
		right := sortedGroups[j]
		if platformSortOrder(left.Platform) != platformSortOrder(right.Platform) {
			return platformSortOrder(left.Platform) < platformSortOrder(right.Platform)
		}
		return strings.ToLower(left.Name) < strings.ToLower(right.Name)
	})

	groupSummaries := make([]modelPlazaPricingGroup, 0, len(sortedGroups))
	rowMap := make(map[string]*pricingRowAccumulator)
	fallbackEntries := h.groupFallbackPricingEntries()

	for _, group := range sortedGroups {
		models := modelsByGroup[group.ID]
		hasExplicitModels := len(models) > 0

		groupSummaries = append(groupSummaries, modelPlazaPricingGroup{
			ID:                group.ID,
			Name:              group.Name,
			Platform:          group.Platform,
			DisplayPrice:      group.DisplayPrice,
			DisplayDiscount:   group.DisplayDiscount,
			Description:       group.Description,
			SubscriptionType:  group.SubscriptionType,
			HasExplicitModels: hasExplicitModels,
			ModelsCount:       len(models),
		})

		if hasExplicitModels {
			for _, requestedModel := range models {
				matchedModel, pricing := h.pricingService.GetMatchedModelPricing(requestedModel)
				if !isPricingTableModel(group.Platform, matchedModel, pricing) {
					continue
				}

				row := upsertPricingRow(rowMap, matchedModel, group.Platform, pricing)
				if alias := strings.TrimSpace(requestedModel); alias != "" && !strings.EqualFold(alias, row.Model) {
					row.Aliases[alias] = struct{}{}
				}
				row.GroupPrices[group.ID] = buildPricingMetrics(pricing, group.EffectiveDisplayRateMultiplier())
			}
			continue
		}

		for _, entry := range fallbackEntries[group.Platform] {
			if !isPricingTableModel(group.Platform, entry.Model, entry.Pricing) {
				continue
			}

			row := upsertPricingRow(rowMap, entry.Model, group.Platform, entry.Pricing)
			row.GroupPrices[group.ID] = buildPricingMetrics(entry.Pricing, group.EffectiveDisplayRateMultiplier())
		}
	}

	keys := make([]string, 0, len(rowMap))
	for key := range rowMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		left := rowMap[keys[i]]
		right := rowMap[keys[j]]
		if platformSortOrder(left.Platform) != platformSortOrder(right.Platform) {
			return platformSortOrder(left.Platform) < platformSortOrder(right.Platform)
		}
		return left.Model < right.Model
	})

	items := make([]modelPlazaPricingItem, 0, len(keys))
	for _, key := range keys {
		row := rowMap[key]
		aliases := make([]string, 0, len(row.Aliases))
		for alias := range row.Aliases {
			if !strings.EqualFold(alias, row.Model) {
				aliases = append(aliases, alias)
			}
		}
		sort.Strings(aliases)

		items = append(items, modelPlazaPricingItem{
			Model:                 row.Model,
			Aliases:               aliases,
			Platform:              row.Platform,
			Provider:              row.Provider,
			Mode:                  row.Mode,
			SupportsPromptCaching: row.SupportsPromptCaching,
			MaxInputTokens:        row.MaxInputTokens,
			MaxOutputTokens:       row.MaxOutputTokens,
			Official:              row.Official,
			GroupPrices:           row.GroupPrices,
		})
	}

	return items, groupSummaries
}

func (h *ModelPlazaHandler) groupFallbackPricingEntries() map[string][]service.NamedLiteLLMModelPricing {
	entries := h.pricingService.ListModelPricingEntries()
	result := map[string][]service.NamedLiteLLMModelPricing{
		"anthropic": {},
		"openai":    {},
		"gemini":    {},
	}

	for _, entry := range entries {
		switch fallbackPlatformForPricingEntry(entry.Model, entry.Pricing) {
		case "anthropic":
			result["anthropic"] = append(result["anthropic"], entry)
		case "openai":
			result["openai"] = append(result["openai"], entry)
		case "gemini":
			result["gemini"] = append(result["gemini"], entry)
		}
	}

	for platform := range result {
		sort.Slice(result[platform], func(i, j int) bool {
			left := normalizeModelPlazaPricingKey(result[platform][i].Model)
			right := normalizeModelPlazaPricingKey(result[platform][j].Model)
			return left < right
		})
	}

	return result
}

func upsertPricingRow(rowMap map[string]*pricingRowAccumulator, model string, platform string, pricing *service.LiteLLMModelPricing) *pricingRowAccumulator {
	key := normalizeModelPlazaPricingKey(model)
	if row, ok := rowMap[key]; ok {
		if row.Provider == "" {
			row.Provider = pricing.LiteLLMProvider
		}
		if row.Mode == "" {
			row.Mode = pricing.Mode
		}
		if pricing.SupportsPromptCaching {
			row.SupportsPromptCaching = true
		}
		if pricing.MaxInputTokens > row.MaxInputTokens {
			row.MaxInputTokens = pricing.MaxInputTokens
		}
		if pricing.MaxOutputTokens > row.MaxOutputTokens {
			row.MaxOutputTokens = pricing.MaxOutputTokens
		}
		return row
	}

	row := &pricingRowAccumulator{
		Model:                 normalizeModelPlazaPricingKey(model),
		Platform:              platform,
		Provider:              pricing.LiteLLMProvider,
		Mode:                  pricing.Mode,
		SupportsPromptCaching: pricing.SupportsPromptCaching,
		MaxInputTokens:        pricing.MaxInputTokens,
		MaxOutputTokens:       pricing.MaxOutputTokens,
		Official:              buildPricingMetrics(pricing, 1),
		Aliases:               make(map[string]struct{}),
		GroupPrices:           make(map[int64]modelPlazaPricingMetrics),
	}
	rowMap[key] = row
	return row
}

func buildPricingMetrics(pricing *service.LiteLLMModelPricing, multiplier float64) modelPlazaPricingMetrics {
	if multiplier < 0 {
		multiplier = 1
	}

	return modelPlazaPricingMetrics{
		InputPerMillion:               pricing.InputCostPerToken * pricingTablePerMillion * multiplier,
		OutputPerMillion:              pricing.OutputCostPerToken * pricingTablePerMillion * multiplier,
		CacheWritePerMillion:          pricing.CacheCreationInputTokenCost * pricingTablePerMillion * multiplier,
		CacheReadPerMillion:           pricing.CacheReadInputTokenCost * pricingTablePerMillion * multiplier,
		InputPerMillionAbove200K:      scaleFloat64Ptr(pricing.InputCostPerTokenAbove200KTokens, multiplier),
		OutputPerMillionAbove200K:     scaleFloat64Ptr(pricing.OutputCostPerTokenAbove200KTokens, multiplier),
		CacheWritePerMillionAbove200K: scaleFloat64Ptr(pricing.CacheCreationInputTokenCostAbove200KTokens, multiplier),
		CacheReadPerMillionAbove200K:  scaleFloat64Ptr(pricing.CacheReadInputTokenCostAbove200KTokens, multiplier),
	}
}

func scaleFloat64Ptr(value *float64, multiplier float64) *float64 {
	if value == nil {
		return nil
	}

	scaled := *value * pricingTablePerMillion * multiplier
	return &scaled
}

func isPricingTableModel(platform string, model string, pricing *service.LiteLLMModelPricing) bool {
	if pricing == nil {
		return false
	}

	mode := strings.ToLower(strings.TrimSpace(pricing.Mode))
	if mode != "chat" && mode != "completion" && mode != "responses" {
		return false
	}
	if pricing.InputCostPerToken == 0 && pricing.OutputCostPerToken == 0 {
		return false
	}

	normalized := normalizeModelPlazaPricingKey(model)

	switch platform {
	case "anthropic":
		if pricing.LiteLLMProvider != "anthropic" {
			return false
		}
		return strings.HasPrefix(normalized, "claude-")
	case "openai":
		if pricing.LiteLLMProvider != "openai" {
			return false
		}
		if strings.HasPrefix(normalized, "ft:") || hasAnyFragment(normalized, []string{
			"audio", "realtime", "search-preview", "transcribe", "tts", "moderation",
		}) {
			return false
		}
		return hasAnyPrefix(normalized, []string{
			"gpt-", "chatgpt-", "o1", "o3", "o4", "codex", "gpt5", "gpt-5",
		})
	case "gemini":
		if pricing.LiteLLMProvider != "vertex_ai-language-models" && pricing.LiteLLMProvider != "gemini" {
			return false
		}
		if hasAnyFragment(normalized, []string{
			"native-audio", "live-", "image-generation", "preview-image", "embedding",
		}) {
			return false
		}
		return strings.HasPrefix(strings.TrimPrefix(normalized, "gemini/"), "gemini-")
	default:
		return false
	}
}

func fallbackPlatformForPricingEntry(model string, pricing *service.LiteLLMModelPricing) string {
	if pricing == nil {
		return ""
	}

	switch pricing.LiteLLMProvider {
	case "anthropic":
		return "anthropic"
	case "openai":
		return "openai"
	case "vertex_ai-language-models":
		if strings.HasPrefix(strings.TrimPrefix(normalizeModelPlazaPricingKey(model), "gemini/"), "gemini-") {
			return "gemini"
		}
	}

	return ""
}

func normalizeModelPlazaPricingKey(model string) string {
	normalized := strings.ToLower(strings.TrimSpace(model))
	normalized = strings.TrimLeft(normalized, "/")
	normalized = strings.TrimPrefix(normalized, "models/")
	normalized = strings.TrimPrefix(normalized, "publishers/google/models/")

	if idx := strings.LastIndex(normalized, "/publishers/google/models/"); idx != -1 {
		normalized = normalized[idx+len("/publishers/google/models/"):]
	}
	if idx := strings.LastIndex(normalized, "/models/"); idx != -1 {
		normalized = normalized[idx+len("/models/"):]
	}

	normalized = strings.TrimPrefix(normalized, "gemini/")
	return strings.TrimSpace(normalized)
}

func hasAnyPrefix(value string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func hasAnyFragment(value string, fragments []string) bool {
	for _, fragment := range fragments {
		if strings.Contains(value, fragment) {
			return true
		}
	}
	return false
}

func platformSortOrder(platform string) int {
	switch platform {
	case "anthropic":
		return 0
	case "openai":
		return 1
	case "gemini":
		return 2
	default:
		return 99
	}
}
