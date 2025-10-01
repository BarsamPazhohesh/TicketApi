package middleware

import (
	"errors"
	"strings"
	"ticket-api/internal/errx"
	"ticket-api/internal/routes"

	"github.com/gin-gonic/gin"
)

func RouteStatusChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := strings.TrimPrefix(c.FullPath(), "/api/v1/")
		method := c.Request.Method

		// Lookup in APIRoutes
		if !routes.IsRouteEnabled(path, method) {
			err := errx.Respond(errx.ErrServiceUnavailable, errors.New("routes status is disable"))
			c.AbortWithStatusJSON(err.HTTPStatus, err)
			return
		}

		c.Next()
	}
}
