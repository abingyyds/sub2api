package admin

import (
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminInviteCodeHandler struct {
	service *service.AdminInviteCodeService
}

func NewAdminInviteCodeHandler(service *service.AdminInviteCodeService) *AdminInviteCodeHandler {
	return &AdminInviteCodeHandler{service: service}
}

type CreateAdminInviteCodeRequest struct {
	SourceName string `json:"source_name" binding:"required,max=100"`
	MaxUses    *int   `json:"max_uses"`
	Notes      string `json:"notes"`
}

type UpdateAdminInviteCodeRequest struct {
	SourceName *string `json:"source_name"`
	MaxUses    *int    `json:"max_uses"`
	Enabled    *bool   `json:"enabled"`
	Notes      *string `json:"notes"`
}

func (h *AdminInviteCodeHandler) Create(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	var req CreateAdminInviteCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	code, err := h.service.Create(c.Request.Context(), req.SourceName, subject.UserID, req.MaxUses, req.Notes)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toDTO(code))
}

func (h *AdminInviteCodeHandler) List(c *gin.Context) {
	page, pageSize := response.GetPaginationParams(c)
	codes, pagination, err := h.service.List(c.Request.Context(), response.ToPaginationParams(page, pageSize))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	dtos := make([]dto.AdminInviteCode, len(codes))
	for i, code := range codes {
		dtos[i] = *toDTO(&code)
	}
	response.SuccessWithPagination(c, dtos, pagination)
}

func (h *AdminInviteCodeHandler) Update(c *gin.Context) {
	id, err := response.GetIDParam(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	var req UpdateAdminInviteCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	code, err := h.service.Update(c.Request.Context(), id, req.SourceName, req.MaxUses, req.Enabled, req.Notes)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toDTO(code))
}

func (h *AdminInviteCodeHandler) Delete(c *gin.Context) {
	id, err := response.GetIDParam(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Invite code deleted"})
}

func toDTO(code *service.AdminInviteCode) *dto.AdminInviteCode {
	return &dto.AdminInviteCode{
		ID:         code.ID,
		Code:       code.Code,
		SourceName: code.SourceName,
		CreatedBy:  code.CreatedBy,
		UsedCount:  code.UsedCount,
		MaxUses:    code.MaxUses,
		Enabled:    code.Enabled,
		Notes:      code.Notes,
		CreatedAt:  code.CreatedAt,
		UpdatedAt:  code.UpdatedAt,
	}
}
