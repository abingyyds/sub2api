package handler

import (
	"sort"
	"sync"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ModelPlazaHandler handles model plaza requests
type ModelPlazaHandler struct {
	apiKeyService  *service.APIKeyService
	gatewayService *service.GatewayService
}

// NewModelPlazaHandler creates a new ModelPlazaHandler
func NewModelPlazaHandler(apiKeyService *service.APIKeyService, gatewayService *service.GatewayService) *ModelPlazaHandler {
	return &ModelPlazaHandler{
		apiKeyService:  apiKeyService,
		gatewayService: gatewayService,
	}
}

// GroupModels represents a group with its available models
type GroupModels struct {
	Group  dto.Group `json:"group"`
	Models []string  `json:"models"`
}

// List returns all available groups with their models for the current user
// GET /api/v1/model-plaza
func (h *ModelPlazaHandler) List(c *gin.Context) {
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

	result := make([]GroupModels, len(groups))

	// Bound concurrency so we avoid a full serial waterfall without turning the
	// request into an unbounded burst of DB work when many groups are visible.
	sem := make(chan struct{}, 6)
	var wg sync.WaitGroup

	for i := range groups {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			group := groups[index]
			groupID := group.ID
			models := h.gatewayService.GetAvailableModels(c.Request.Context(), &groupID, "")
			sort.Strings(models)

			result[index] = GroupModels{
				Group:  *dto.GroupFromService(&group),
				Models: models,
			}
		}(i)
	}

	wg.Wait()

	response.Success(c, result)
}
