package middleware

import (
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyOrganization is the gin context key for the organization
	ContextKeyOrganization ContextKey = "organization"
	// ContextKeyOrgMember is the gin context key for the org member
	ContextKeyOrgMember ContextKey = "org_member"
	// ContextKeyOrgProject is the gin context key for the org project
	ContextKeyOrgProject ContextKey = "org_project"
)

// NewOrgAuthMiddleware creates an org admin authentication middleware
func NewOrgAuthMiddleware(authService *service.AuthService, userService *service.UserService, orgService *service.OrganizationService) OrgAuthMiddleware {
	return OrgAuthMiddleware(orgAuth(authService, userService, orgService))
}

// orgAuth creates the org admin middleware handler
func orgAuth(authService *service.AuthService, userService *service.UserService, orgService *service.OrganizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Reuse JWT auth subject from context (set by jwtAuth middleware)
		subject, ok := GetAuthSubjectFromContext(c)
		if !ok {
			AbortWithError(c, 401, "UNAUTHORIZED", "Authentication required")
			return
		}

		// Get user from DB
		user, err := userService.GetByID(c.Request.Context(), subject.UserID)
		if err != nil {
			AbortWithError(c, 401, "USER_NOT_FOUND", "User not found")
			return
		}

		// Check user is org_admin
		if !user.IsOrgAdmin() {
			AbortWithError(c, 403, "NOT_ORG_ADMIN", "Organization admin access required")
			return
		}

		// Load the organization owned by this user
		org, err := orgService.GetByOwnerID(c.Request.Context(), user.ID)
		if err != nil {
			AbortWithError(c, 403, "ORG_NOT_FOUND", "No organization found for this user")
			return
		}

		if !org.IsActive() {
			AbortWithError(c, 403, "ORG_INACTIVE", "Organization is not active")
			return
		}

		// Set org and member in context
		c.Set(string(ContextKeyOrganization), org)

		c.Next()
	}
}

// GetOrganizationFromContext retrieves the organization from gin context
func GetOrganizationFromContext(c *gin.Context) (*service.Organization, bool) {
	value, exists := c.Get(string(ContextKeyOrganization))
	if !exists {
		return nil, false
	}
	org, ok := value.(*service.Organization)
	return org, ok
}

// GetOrgMemberFromContext retrieves the org member from gin context
func GetOrgMemberFromContext(c *gin.Context) (*service.OrgMember, bool) {
	value, exists := c.Get(string(ContextKeyOrgMember))
	if !exists {
		return nil, false
	}
	member, ok := value.(*service.OrgMember)
	return member, ok
}

// GetOrgProjectFromContext retrieves the org project from gin context
func GetOrgProjectFromContext(c *gin.Context) (*service.OrgProject, bool) {
	value, exists := c.Get(string(ContextKeyOrgProject))
	if !exists {
		return nil, false
	}
	project, ok := value.(*service.OrgProject)
	return project, ok
}
