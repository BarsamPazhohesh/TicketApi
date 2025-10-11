package services

import (
	"ticket-api/internal/services/captcha"
	"ticket-api/internal/services/token"
)

type AppServices struct {
	Captcha *captcha.CaptchaService
	Token   *token.TokenService
}

func NewAppService() *AppServices {
	return &AppServices{
		Captcha: captcha.NewCaptchaService(),
		Token:   token.NewTokenService(),
	}
}
