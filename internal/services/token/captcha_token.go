package token

import (
	"errors"
	"ticket-api/internal/config"
	"ticket-api/internal/errx"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CaptchaClaims defines the claims inside a captcha token
type CaptchaClaims struct {
	jwt.RegisteredClaims
}

// NewCaptchaToken creates a captcha JWT token using config
func (s *TokenService) NewCaptchaToken() (string, *errx.APIError) {
	secret, errSecret := secretKeyBytes()
	if errSecret != nil {
		return "", errSecret
	}
	cfg := config.Get().Captcha
	claims := CaptchaClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(cfg.ExpiredTimeToken))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ticket-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", errx.Respond(errx.ErrInternalServerError, err)
	}
	return signed, nil
}

// ParseCaptchaToken parses captcha token
func (s *TokenService) ParseCaptchaToken(tokenString string) (*CaptchaClaims, *errx.APIError) {
	claims := &CaptchaClaims{}
	secret, errSecret := secretKeyBytes()
	if errSecret != nil {
		return nil, errSecret
	}

	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, errx.Respond(errx.ErrUnauthorized, err)
	}
	if !parsed.Valid {
		return nil, errx.Respond(errx.ErrUnauthorized, errors.New("invalid or expired token"))
	}
	return claims, nil
}
