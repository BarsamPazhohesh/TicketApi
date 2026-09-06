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
// If requirePhone is true, the captcha token must contain a verified non-empty phone number.
func CaptchaMiddleware(tokenService *token.TokenService, requirePhone bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieService := cookie.NewCaptchaCookieService()
		authService := cookie.NewAuthCookieService()
		userIP := c.ClientIP()

		// Check for user auth token
		authToken, errCookie := authService.Get(c)
		if errCookie == nil {
			// Validate auth token
			claims, err := tokenService.ParseAuthToken(c.Request.Context(), authToken)
			if err == nil {
				// Auth token is valid, skip captcha
				c.Set("user", claims)
				c.Next()
				return
			}

			if err.Err.Code != errx.ErrUnauthorized {
				c.AbortWithStatusJSON(err.HTTPStatus, err)
				return
			}
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

		// Check verified phone requirement for Level 2 actions
		if requirePhone && parsedCaptchaToken.PhoneNumber == "" {
			appErr := errx.Respond(errx.ErrUnauthorized, errors.New("phone verification required"))
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr)
			return
		}

		if parsedCaptchaToken.PhoneNumber != "" {
			c.Set("guest_phone", parsedCaptchaToken.PhoneNumber)
		}

		c.Next()
	}
}
