package handler

import (
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
	pricingService *service.PricingService
}

// NewModelPlazaHandler creates a new ModelPlazaHandler
func NewModelPlazaHandler(apiKeyService *service.APIKeyService, gatewayService *service.GatewayService, pricingService *service.PricingService) *ModelPlazaHandler {
	return &ModelPlazaHandler{
		apiKeyService:  apiKeyService,
		gatewayService: gatewayService,
		pricingService: pricingService,
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

	groupIDs := make([]int64, 0, len(groups))
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID)
	}

	modelsByGroup := h.gatewayService.GetAvailableModelsByGroups(c.Request.Context(), groupIDs, "")

	result := make([]GroupModels, len(groups))
	for i, group := range groups {
		result[i] = GroupModels{
			Group:  *dto.GroupFromService(&group),
			Models: modelsByGroup[group.ID],
		}
	}

	response.Success(c, result)
}
