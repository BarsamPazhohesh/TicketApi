package services

import (
	"ticket-api/internal/services/cache"
	"ticket-api/internal/services/captcha"
	"ticket-api/internal/services/storage"
	"ticket-api/internal/services/token"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type AppServices struct {
	Captcha     *captcha.CaptchaService
	Token       *token.TokenService
	Cache       *cache.CacheService
	FileStorage *storage.StorageService
}

func NewAppService(redis *redis.Client, minio *minio.Client) *AppServices {
	return &AppServices{
		Captcha:     captcha.NewCaptchaService(),
		Token:       token.NewTokenService(),
		Cache:       cache.NewCacheService(redis),
		FileStorage: storage.NewStorageService(minio),
	}
}
