package token

import (
	"context"
	"fmt"
	"ticket-api/internal/config"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

// TokenService handles JWT generation, parsing, and revocation
type TokenService struct {
	cache *cache.Cache
	redis *redis.Client
}

// NewTokenService creates a new service instance
func NewTokenService(redis *redis.Client) *TokenService {
	cfg := config.Get().OneTimeToken
	timeout := time.Duration(cfg.ExpiredTimeToken) * time.Minute
	cleanupInterval := time.Duration(cfg.CleanupInterval) * time.Minute

	return &TokenService{
		cache: cache.New(timeout, cleanupInterval),
		redis: redis,
	}
}

// RevokeToken marks a JWT ID (JTI) as revoked in Redis until its expiration time
func (s *TokenService) RevokeToken(ctx context.Context, jti string, duration time.Duration) error {
	if s.redis == nil || jti == "" {
		return nil
	}
	key := fmt.Sprintf("revoked:token:%s", jti)
	return s.redis.Set(ctx, key, "revoked", duration).Err()
}

// IsTokenRevoked checks if a token JTI is in the revocation blocklist. Fail closed on error.
func (s *TokenService) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	if s.redis == nil || jti == "" {
		return false, nil
	}
	key := fmt.Sprintf("revoked:token:%s", jti)
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return true, fmt.Errorf("failed to verify token revocation with redis: %w", err)
	}
	return exists > 0, nil
}
