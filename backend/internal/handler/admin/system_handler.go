package admin

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// SystemHandler handles system-related operations
type SystemHandler struct {
	version string
}

// NewSystemHandler creates a new SystemHandler
func NewSystemHandler(version string) *SystemHandler {
	return &SystemHandler{
		version: version,
	}
}

// GetVersion returns the current version
// GET /api/v1/admin/system/version
func (h *SystemHandler) GetVersion(c *gin.Context) {
	response.Success(c, gin.H{
		"version": h.version,
	})
}
