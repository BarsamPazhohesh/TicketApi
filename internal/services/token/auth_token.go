package token

import (
	"errors"
	"ticket-api/internal/config"
	"ticket-api/internal/errx"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthClaims holds claims for authentication tokens
type AuthClaims struct {
	UserID   int64   `json:"user_id"`
	Username string  `json:"username"`
	RoleIDs  []int64 `json:"role_ids"`

	jwt.RegisteredClaims
}

// NewAuthToken creates auth token using config
func (s *TokenService) NewAuthToken(credential AuthClaims) (string, *errx.APIError) {
	secret, errSecret := secretKeyBytes()
	if errSecret != nil {
		return "", errSecret
	}

	cfg := config.Get().Auth
	claims := AuthClaims{
		UserID:   credential.UserID,
		Username: credential.Username,
		RoleIDs:  credential.RoleIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpiredTimeToken) * time.Minute)),
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

// ParseAuthToken parses auth token
func (s *TokenService) ParseAuthToken(tokenString string) (*AuthClaims, *errx.APIError) {
	secret, errSecret := secretKeyBytes()
	if errSecret != nil {
		return nil, errSecret
	}

	claims := &AuthClaims{}
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
