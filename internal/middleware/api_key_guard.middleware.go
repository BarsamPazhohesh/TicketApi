package middleware

import (
	"errors"
	"ticket-api/internal/config"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
)

// ApiKeyGuardMiddleware validates the presence and format of the API key in requests.
func ApiKeyGuardMiddleware(tokenSvc *token.TokenService, apiRepo *repository.APIKeysRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get("x-api-key")

		// Check if the API key length is valid
		if len(key) < config.Get().APIKey.Size {
			appErr := errx.Respond(errx.ErrUnauthorized, errors.New("invalid API key format"))
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr.Err)
			return
		}

		// Hash the API key before querying the repository
		hashKey := tokenSvc.Hash(key)
		_, appErr := apiRepo.GetApiKeyIDByKey(c.Request.Context(), hashKey)
		if appErr != nil {
			c.AbortWithStatusJSON(appErr.HTTPStatus, appErr.Err)
			return
		}

		// Continue to the next middleware/handler
		c.Next()

	}
}
