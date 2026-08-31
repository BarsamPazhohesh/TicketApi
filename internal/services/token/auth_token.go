package token

import (
	"context"
	"errors"
	"ticket-api/internal/config"
	"ticket-api/internal/errx"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthClaims holds claims for authentication tokens
type AuthClaims struct {
	UserID      int64    `json:"user_id"`
	Username    string   `json:"username"`
	RoleIDs     []int64  `json:"role_ids"`
	Permissions []string `json:"permissions,omitempty"`

	jwt.RegisteredClaims
}

// NewAuthToken creates auth token with unique JTI and claims
func (s *TokenService) NewAuthToken(credential AuthClaims) (string, *errx.APIError) {
	secret, errSecret := secretKeyBytes()
	if errSecret != nil {
		return "", errSecret
	}

	cfg := config.Get().Auth
	now := time.Now()
	tokenID := uuid.New().String()

	claims := AuthClaims{
		UserID:      credential.UserID,
		Username:    credential.Username,
		RoleIDs:     credential.RoleIDs,
		Permissions: credential.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.ExpiredTimeToken) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
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

// ParseAuthToken parses, validates, and checks revocation of auth token
func (s *TokenService) ParseAuthToken(ctx context.Context, tokenString string) (*AuthClaims, *errx.APIError) {
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

	if claims.ID != "" {
		revoked, revErr := s.IsTokenRevoked(ctx, claims.ID)
		if revErr != nil {
			return nil, errx.Respond(errx.ErrServiceUnavailable, revErr)
		}
		if revoked {
			return nil, errx.Respond(errx.ErrUnauthorized, errors.New("token has been revoked"))
		}
	}

	return claims, nil
}
