package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// AnnouncementHandler handles public announcement endpoints
type AnnouncementHandler struct {
	svc *service.AnnouncementService
}

func NewAnnouncementHandler(svc *service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{svc: svc}
}

// GetActive returns all active announcements (public, no auth)
func (h *AnnouncementHandler) GetActive(c *gin.Context) {
	items, err := h.svc.ListActive(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]dto.Announcement, 0, len(items))
	for i := range items {
		out = append(out, *dto.AnnouncementFromService(&items[i]))
	}
	response.Success(c, out)
}
