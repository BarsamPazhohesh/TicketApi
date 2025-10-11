package token

import (
	"errors"
	"sync"

	"ticket-api/internal/env"
	"ticket-api/internal/errx"
)

var (
	secretKey []byte
	once      sync.Once
)

// secretKeyBytes loads and caches the secret key used to sign tokens
func secretKeyBytes() ([]byte, *errx.APIError) {
	once.Do(func() {
		secret := env.GetEnvString("JWT_SECRET", "")
		if len(secret) < 32 {
			secretKey = nil
		} else {
			secretKey = []byte(secret)
		}
	})

	if len(secretKey) == 0 {
		return nil, errx.Respond(
			errx.ErrWeakJWTSecret,
			errors.New("token secret missing or too short"),
		)
	}
	return secretKey, nil
}

