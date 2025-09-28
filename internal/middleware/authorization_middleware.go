// Package middleware
package middleware

import (
	"ticket-api/internal/errx"
	"ticket-api/internal/services/cookie"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware(tokenService *token.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieService := cookie.NewAuthCookieService() // returns *CookieService
		authToken, errCookie := cookieService.Get(c)
		if errCookie != nil {
			errApp := errx.Respond(errx.ErrUnauthorized, errCookie)
			c.AbortWithStatusJSON(errApp.HTTPStatus, errApp)
			return
		}

		// parse and validate the token
		user, err := tokenService.ParseAuthToken(authToken)
		if err != nil {
			c.AbortWithStatusJSON(err.HTTPStatus, err)
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
