package handler

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// SubSiteHandler handles user-facing sub-site requests.
type SubSiteHandler struct {
	subSiteService *service.SubSiteService
}

func NewSubSiteHandler(subSiteService *service.SubSiteService) *SubSiteHandler {
	return &SubSiteHandler{subSiteService: subSiteService}
}

// GetOpenInfo returns self-service opening information for the current site scope.
func (h *SubSiteHandler) GetOpenInfo(c *gin.Context) {
	info, err := h.subSiteService.GetOpenInfo(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, info)
}

// ListOwned returns all sub-sites owned by the current user.
func (h *SubSiteHandler) ListOwned(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	items, err := h.subSiteService.ListOwnedSites(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, items)
}

// GetOwned returns one owned sub-site detail.
func (h *SubSiteHandler) GetOwned(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	siteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || siteID <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	item, err := h.subSiteService.GetOwnedSite(c.Request.Context(), subject.UserID, siteID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

// UpdateOwned updates a sub-site owned by the current user.
func (h *SubSiteHandler) UpdateOwned(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	siteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || siteID <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	var req service.UpdateOwnedSubSiteInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	req.ID = siteID
	item, err := h.subSiteService.UpdateOwnedSite(c.Request.Context(), subject.UserID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}
