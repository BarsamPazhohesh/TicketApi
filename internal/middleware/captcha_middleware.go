package middleware

import (
	"errors"
	"ticket-api/internal/config"
	"ticket-api/internal/errx"
	"ticket-api/internal/services/cookie"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
)

// CaptchaMiddleware ensures that either a valid auth token or a captcha token is present.
func CaptchaMiddleware(tokenService *token.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieService := cookie.NewCaptchaCookieService()
		authService := cookie.NewAuthCookieService()
		userIP := c.ClientIP()

		// Check for user auth token
		authToken, errCookie := authService.Get(c)
		if errCookie != nil {
			appErr := errx.Respond(errx.ErrUnauthorized, errCookie)
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr)
			return
		}

		// Validate auth token
		_, err := tokenService.ParseAuthToken(authToken)
		if err == nil {
			// Auth token is valid, skip captcha
			c.Next()
			return
		}

		if err.Err.Code != errx.ErrUnauthorized {
			c.AbortWithStatusJSON(err.HTTPStatus, err)
			return
		}

		// Check for captcha token
		captchaToken, errToken := cookieService.Get(c)
		if errToken != nil {
			appErr := errx.Respond(errx.ErrUnauthorized, errToken)
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr)
			return
		}

		// Validate captcha token
		parsedCaptchaToken, err := tokenService.ParseCaptchaToken(captchaToken)
		if err != nil {
			c.AbortWithStatusJSON(err.HTTPStatus, err)
			return
		}

		// Optional IP validation
		if config.Get().Captcha.ValidateIP && parsedCaptchaToken.IP != userIP {
			appErr := errx.Respond(errx.ErrUnauthorized, errors.New("user IP does not match captcha token IP"))
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr)
			return
		}

		c.Next()
	}
}
