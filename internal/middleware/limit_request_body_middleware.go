package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LimitRequestBody(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit<<10)
		c.Next()
	}
}
