package admin

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// AnnouncementHandler handles admin announcement CRUD
type AnnouncementHandler struct {
	svc *service.AnnouncementService
}

func NewAnnouncementHandler(svc *service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{svc: svc}
}

// List handles listing announcements with pagination
func (h *AnnouncementHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	items, pag, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]dto.Announcement, 0, len(items))
	for i := range items {
		out = append(out, *dto.AnnouncementFromService(&items[i]))
	}
	response.Paginated(c, out, pag.Total, pag.Page, pag.PageSize)
}

// Create handles creating a new announcement
func (h *AnnouncementHandler) Create(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content"`
		Status   string `json:"status"`
		Priority int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Status == "" {
		req.Status = "active"
	}
	a := &service.Announcement{
		Title:    req.Title,
		Content:  req.Content,
		Status:   req.Status,
		Priority: req.Priority,
	}
	if err := h.svc.Create(c.Request.Context(), a); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.AnnouncementFromService(a))
}

// Update handles updating an announcement
func (h *AnnouncementHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	existing, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	var req struct {
		Title    *string `json:"title"`
		Content  *string `json:"content"`
		Status   *string `json:"status"`
		Priority *int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Content != nil {
		existing.Content = *req.Content
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.Priority != nil {
		existing.Priority = *req.Priority
	}
	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.AnnouncementFromService(existing))
}

// Delete handles deleting an announcement
func (h *AnnouncementHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, nil)
}
