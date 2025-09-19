package middleware

import (
	"ticket-api/internal/errx"

	"github.com/gin-gonic/gin"
	"ticket-api/internal/services/cookie"
	"ticket-api/internal/services/token"
)

func CaptchaMiddleware(tokenService *token.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieService := cookie.NewCaptchaCookieService()

		// user is already authorized, skip captcha
		if _, exists := c.Get("user"); exists {
			c.Next()
			return
		}

		captchaToken, errToken := cookieService.Get(c)
		if errToken != nil {
			errApp := errx.Respond(errx.ErrUnauthorized, errToken)
			c.AbortWithStatusJSON(errApp.HTTPStatus, errApp)
			return
		}

		// parse and validate captcha token
		_, err := tokenService.ParseCaptchaToken(captchaToken)
		if err != nil {
			c.AbortWithStatusJSON(err.HTTPStatus, err)
			return
		}

		c.Next()
	}
}
