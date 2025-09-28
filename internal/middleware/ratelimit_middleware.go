package middleware

import (
	"context"
	"errors"
	"fmt"
	"net"
	"ticket-api/internal/errx"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(rdb *redis.Client, limitPerMinute int) gin.HandlerFunc {
	limiter := redis_rate.NewLimiter(rdb)

	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			ip = c.ClientIP() // fallback
		}

		// build key using IP + endpoint path
		key := fmt.Sprintf("rate:%s:%s", ip, c.FullPath())

		// allow limitPerMinute per minute
		res, err := limiter.Allow(
			context.Background(),
			key,
			redis_rate.PerMinute(limitPerMinute),
		)
		if err != nil {
			appErr := errx.Respond(errx.ErrInternalServerError, err)
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr)
			return
		}

		if res.Allowed == 0 {
			msg := fmt.Sprintf(
				"error=%s retry_after=%v remaining=%v reset_in_sec=%v",
				"Too Many Requests",
				res.RetryAfter/time.Second,
				res.Remaining,
				res.ResetAfter.Seconds(),
			)
			appErr := errx.Respond(errx.ErrTooManyRequest, errors.New(msg))
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr)
			return
		}

		c.Next()
	}
}
