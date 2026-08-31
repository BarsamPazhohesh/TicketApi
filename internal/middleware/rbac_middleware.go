package middleware

import (
	"errors"
	"strings"
	"ticket-api/internal/errx"
	"ticket-api/internal/security"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
)

// DynamicRBACGuardMiddleware verifies user roles against the in-memory PBAC registry
func DynamicRBACGuardMiddleware(registry *security.SecurityRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("user")
		if !exists {
			appErr := errx.Respond(errx.ErrUnauthorized, errors.New("missing user in context"))
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr)
			return
		}

		user, ok := val.(*token.AuthClaims)
		if !ok || user == nil {
			appErr := errx.Respond(errx.ErrUnauthorized, errors.New("invalid user claims"))
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr)
			return
		}

		// Normalize path relative to /api/v1/
		fullPath := c.FullPath()
		trimmedPath := strings.TrimPrefix(fullPath, "/api/v1/")

		allowed, routeEnabled := registry.CanAccessRoute(user.RoleIDs, c.Request.Method, trimmedPath)
		if !routeEnabled {
			appErr := errx.Respond(errx.ErrNotFound, errors.New("route is disabled"))
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr)
			return
		}

		if !allowed {
			appErr := errx.Respond(errx.ErrForbidden, errors.New("insufficient permissions for this resource"))
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr)
			return
		}

		c.Next()
	}
}
