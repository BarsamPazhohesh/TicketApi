package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type memoryItem struct {
	data      []byte
	expiresAt time.Time
}

type CacheService struct {
	redis     *redis.Client
	memory    map[string]memoryItem
	memoryMu  sync.RWMutex
}

// NewCacheService creates a new cache service
func NewCacheService(redis *redis.Client) *CacheService {
	return &CacheService{
		redis:  redis,
		memory: make(map[string]memoryItem),
	}
}

// Set any struct as JSON
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if c == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if c.redis != nil {
		return c.redis.Set(ctx, key, data, ttl).Err()
	}

	// In-memory fallback (when redis is nil, e.g. tests / local no-redis mode)
	c.memoryMu.Lock()
	defer c.memoryMu.Unlock()
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	c.memory[key] = memoryItem{
		data:      data,
		expiresAt: exp,
	}
	return nil
}

// Get JSON into a struct
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	if c == nil {
		return false, nil
	}

	if c.redis != nil {
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

	// In-memory fallback
	c.memoryMu.RLock()
	item, ok := c.memory[key]
	c.memoryMu.RUnlock()

	if !ok {
		return false, nil
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		c.memoryMu.Lock()
		delete(c.memory, key)
		c.memoryMu.Unlock()
		return false, nil
	}

	err := json.Unmarshal(item.data, dest)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Delete key (for invalidation)
func (c *CacheService) Delete(ctx context.Context, key string) error {
	if c == nil {
		return nil
	}

	if c.redis != nil {
		return c.redis.Del(ctx, key).Err()
	}

	c.memoryMu.Lock()
	delete(c.memory, key)
	c.memoryMu.Unlock()
	return nil
}
