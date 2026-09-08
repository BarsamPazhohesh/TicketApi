package services

import (
	"ticket-api/internal/repository"
	"ticket-api/internal/services/auth"
	"ticket-api/internal/services/cache"
	"ticket-api/internal/services/captcha"
	"ticket-api/internal/services/otp"
	"ticket-api/internal/services/smsloader"
	"ticket-api/internal/services/storage"
	"ticket-api/internal/services/ticket"
	"ticket-api/internal/services/token"
	"ticket-api/internal/services/user"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type AppServices struct {
	Captcha     *captcha.CaptchaService
	Token       *token.TokenService
	Cache       *cache.CacheService
	FileStorage *storage.StorageService
	Auth        *auth.AuthService
	Ticket      *ticket.TicketService
	User        *user.UserService
	OTP         *otp.OTPService
	SMSLoader   *smsloader.SMSLoaderService
}

func NewAppService(redis *redis.Client, minio *minio.Client, repos *repository.AppRepositories) *AppServices {
	captchaSvc := captcha.NewCaptchaService()
	tokenSvc := token.NewTokenService(redis)
	cacheSvc := cache.NewCacheService(redis)
	storageSvc := storage.NewStorageService(minio)
	smsLoaderSvc := smsloader.NewSMSLoaderService(repos.SMSTypeMessages)
	otpSvc := otp.NewOTPService(repos.SMSWarehouse, cacheSvc, smsLoaderSvc)

	authSvc := auth.NewAuthService(repos.Users, repos.RolesRelations, tokenSvc)
	ticketSvc := ticket.NewTicketService(
		repos.Ticket,
		repos.ChatRepository,
		repos.TicketTypes,
		repos.TicketPriorities,
		repos.TicketStatus,
		repos.Users,
		repos.Departments,
		storageSvc,
	)
	userSvc := user.NewUserService(repos.Users)

	return &AppServices{
		Captcha:     captchaSvc,
		Token:       tokenSvc,
		Cache:       cacheSvc,
		FileStorage: storageSvc,
		Auth:        authSvc,
		Ticket:      ticketSvc,
		User:        userSvc,
		OTP:         otpSvc,
		SMSLoader:   smsLoaderSvc,
	}
}
