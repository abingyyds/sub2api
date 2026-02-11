package admin

import (
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// ReferralHandler handles admin referral management
type ReferralHandler struct {
	referralService *service.ReferralService
}

// NewReferralHandler creates a new admin referral handler
func NewReferralHandler(referralService *service.ReferralService) *ReferralHandler {
	return &ReferralHandler{referralService: referralService}
}

// List handles listing all referral records with pagination
// GET /api/v1/admin/referrals
func (h *ReferralHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	search := strings.TrimSpace(c.Query("search"))
	if len(search) > 100 {
		search = search[:100]
	}

	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	referrals, pag, err := h.referralService.ListAll(c.Request.Context(), params, search)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Paginated(c, referrals, pag.Total, pag.Page, pag.PageSize)
}

// GetStats returns global referral statistics
// GET /api/v1/admin/referrals/stats
func (h *ReferralHandler) GetStats(c *gin.Context) {
	stats, err := h.referralService.GetStats(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, stats)
}
