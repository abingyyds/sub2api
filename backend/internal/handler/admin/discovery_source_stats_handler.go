package admin

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type DiscoverySourceStatsHandler struct {
	userService *service.UserService
}

func NewDiscoverySourceStatsHandler(userService *service.UserService) *DiscoverySourceStatsHandler {
	return &DiscoverySourceStatsHandler{userService: userService}
}

type DiscoverySourceStats struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

type DiscoverySourceStatsResponse struct {
	Day7  []DiscoverySourceStats `json:"day7"`
	Day1  []DiscoverySourceStats `json:"day1"`
	Total int                    `json:"total"`
}

func (h *DiscoverySourceStatsHandler) GetStats(c *gin.Context) {
	now := time.Now()
	day1Start := now.Add(-24 * time.Hour)
	day7Start := now.Add(-7 * 24 * time.Hour)

	stats7, err := h.userService.GetDiscoverySourceStats(c.Request.Context(), day7Start, now)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	stats1, err := h.userService.GetDiscoverySourceStats(c.Request.Context(), day1Start, now)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	total7 := 0
	for _, s := range stats7 {
		total7 += s.Count
	}

	response.Success(c, DiscoverySourceStatsResponse{
		Day7:  convertStats(stats7),
		Day1:  convertStats(stats1),
		Total: total7,
	})
}

func convertStats(stats []service.DiscoverySourceStat) []DiscoverySourceStats {
	result := make([]DiscoverySourceStats, len(stats))
	for i, s := range stats {
		result[i] = DiscoverySourceStats{
			Source: s.Source,
			Count:  s.Count,
		}
	}
	return result
}
