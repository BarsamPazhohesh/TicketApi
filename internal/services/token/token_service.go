package token

import (
	"ticket-api/internal/config"
	"time"

	"github.com/patrickmn/go-cache"
)

// TokenService handles JWT generation and parsing
type TokenService struct {
	cache *cache.Cache
}

// NewTokenService creates a new service instance and loads the secret
func NewTokenService() *TokenService {
	cfg := config.Get().OneTimeToken
	timeout := time.Duration(cfg.ExpiredTimeToken) * time.Minute
	cleanupInterval := time.Duration(cfg.CleanupInterval) * time.Minute

	return &TokenService{
		cache: cache.New(timeout, cleanupInterval),
	}
}
