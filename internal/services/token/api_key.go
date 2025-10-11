package token

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"ticket-api/internal/config"
	"ticket-api/internal/errx"
)

// Generate creates a secure random API key
func (token *TokenService) GenerateAPIKey() (string, *errx.APIError) {
	sizeBytes := config.Get().APIKey.Size
	bytes := make([]byte, sizeBytes)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", errx.Respond(errx.ErrInternalServerError, err)
	}
	return hex.EncodeToString(bytes), nil
}

// Hash computes SHA-256 hash of the API key for storage
func (token *TokenService) Hash(apiKey string) string {
	h := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(h[:])
}

// Validate compares the API key with stored hash securely
func (token *TokenService) ValidateAPIKey(apiKey string, hashedKey string) bool {
	hash := token.Hash(apiKey)
	if len(hash) != len(hashedKey) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(hash), []byte(hashedKey)) == 1
}
