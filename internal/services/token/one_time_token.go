package token

import (
	"ticket-api/internal/config"
	"ticket-api/internal/errx"
	"ticket-api/internal/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// OneTimeTokenClaims defines the claims inside a one-time token
type OneTimeTokenClaims struct {
	TokenID  string
	Username string
	jwt.RegisteredClaims
}

// NewOneTimeToken creates a one-time JWT token and caches its ID for single-use enforcement
func (s *TokenService) NewOneTimeToken(username string) (string, *errx.APIError) {
	secret, errSecret := secretKeyBytes()
	if errSecret != nil {
		return "", errSecret
	}

	id := util.GenerateUUID()
	cfg := config.Get().OneTimeToken

	claims := OneTimeTokenClaims{
		Username: username,
		TokenID:  id,
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

	// Only store the TokenID in cache; value can be anything (bool, struct{}, etc.)
	s.cache.Set(id, struct{}{}, time.Duration(cfg.ExpiredTimeToken)*time.Minute)

	return signed, nil
}

// ParseOneTimeToken parses a one-time JWT token, validates it, and enforces single-use
func (s *TokenService) ParseOneTimeToken(tokenString string) (*OneTimeTokenClaims, *errx.APIError) {
	secret, errSecret := secretKeyBytes()
	if errSecret != nil {
		return nil, errSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &OneTimeTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, errx.Respond(errx.ErrUnauthorized, err)
	}

	claims, ok := token.Claims.(*OneTimeTokenClaims)
	if !ok || !token.Valid {
		return nil, errx.Respond(errx.ErrUnauthorized, nil)
	}

	// Enforce single-use by checking cache
	if _, found := s.cache.Get(claims.TokenID); !found {
		return nil, errx.Respond(errx.ErrUnauthorized, nil)
	}
	s.cache.Delete(claims.TokenID)

	return claims, nil
}
