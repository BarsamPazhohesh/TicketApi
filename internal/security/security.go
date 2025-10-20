package security

import (
	"errors"
	"ticket-api/internal/errx"

	"golang.org/x/crypto/bcrypt"
)

// Returns nil if match, *errx.APIError if invalid.
func CompareHashPassword(hashedPassword, plainPassword string) *errx.APIError {
	if plainPassword == "" {
		return errx.Respond(errx.ErrInvalidCredentials, errors.New("invalid credentials"))
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		// bcrypt returns error on mismatch or bad hash
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errx.Respond(errx.ErrInvalidCredentials, errors.New("invalid credentials"))
		}
		// any other bcrypt error
		return errx.Respond(errx.ErrInternalServerError, err)
	}
	return nil
}
