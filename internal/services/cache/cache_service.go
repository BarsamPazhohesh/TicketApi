package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	redis *redis.Client
}

// NewCacheService creates a new cache service
func NewCacheService(redis *redis.Client) *CacheService {
	return &CacheService{
		redis: redis,
	}
}

// Set any struct as JSON
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, data, ttl).Err()
}

// Get JSON into a struct
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	val, err := c.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // not found
	} else if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Delete key (for invalidation)
func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}
