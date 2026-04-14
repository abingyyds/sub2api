package admin

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type SubSiteHandler struct {
	subSiteService *service.SubSiteService
}

func NewSubSiteHandler(subSiteService *service.SubSiteService) *SubSiteHandler {
	return &SubSiteHandler{subSiteService: subSiteService}
}

func (h *SubSiteHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	items, pag, err := h.subSiteService.List(c.Request.Context(), params, c.Query("search"), c.Query("status"))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}

func (h *SubSiteHandler) Create(c *gin.Context) {
	var req service.CreateSubSiteInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	site, err := h.subSiteService.Create(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, site)
}

func (h *SubSiteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	var req service.UpdateSubSiteInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	req.ID = id
	site, err := h.subSiteService.Update(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, site)
}

func (h *SubSiteHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	if err := h.subSiteService.Delete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "sub-site deleted"})
}
