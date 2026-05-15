package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ChangePasswordRequest represents the change password request payload
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdateProfileRequest represents the update profile request payload
type UpdateProfileRequest struct {
	Username *string `json:"username"`
}

// AcceptLegalAgreementRequest represents the legal agreement acceptance payload.
type AcceptLegalAgreementRequest struct {
	TermsAccepted    bool `json:"terms_accepted"`
	PrivacyAccepted  bool `json:"privacy_accepted"`
	APITermsAccepted bool `json:"api_terms_accepted"`
}

// GetProfile handles getting user profile
// GET /api/v1/users/me
func (h *UserHandler) GetProfile(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userData, err := h.userService.GetByID(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.UserFromService(userData))
}

// ChangePassword handles changing user password
// POST /api/v1/users/me/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	svcReq := service.ChangePasswordRequest{
		CurrentPassword: req.OldPassword,
		NewPassword:     req.NewPassword,
	}
	err := h.userService.ChangePassword(c.Request.Context(), subject.UserID, svcReq)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Password changed successfully"})
}

// UpdateProfile handles updating user profile
// PUT /api/v1/users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	svcReq := service.UpdateProfileRequest{
		Username: req.Username,
	}
	updatedUser, err := h.userService.UpdateProfile(c.Request.Context(), subject.UserID, svcReq)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.UserFromService(updatedUser))
}

// AcceptLegalAgreement records that the current user accepted the latest legal documents.
// POST /api/v1/user/legal-agreement
func (h *UserHandler) AcceptLegalAgreement(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req AcceptLegalAgreementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if !req.TermsAccepted || !req.PrivacyAccepted || !req.APITermsAccepted {
		response.BadRequest(c, "You must read and accept the User Agreement, Privacy Policy, and API usage terms before using the API.")
		return
	}

	if err := h.userService.AcceptLegalAgreement(c.Request.Context(), subject.UserID); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	userData, err := h.userService.GetProfile(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.UserFromService(userData))
}

// UpdateDiscoverySourceRequest 更新用户来源请求
type UpdateDiscoverySourceRequest struct {
	Source string `json:"source" binding:"required,max=50"`
}

// UpdateDiscoverySource 更新用户来源渠道
// PUT /api/v1/user/discovery-source
func (h *UserHandler) UpdateDiscoverySource(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req UpdateDiscoverySourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	updatedUser, err := h.userService.UpdateDiscoverySource(c.Request.Context(), subject.UserID, req.Source)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.UserFromService(updatedUser))
}
