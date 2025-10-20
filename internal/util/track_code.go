// Package util
package util

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	// Allowed characters: uppercase letters + digits
	ticketChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ticketLength = 8  // track code Length
	maxAttempts  = 10 // max retries for uniqueness
)

// generateTrackCode creates a random ticket code of fixed length.
func generateTrackCode() (string, error) {
	b := make([]byte, ticketLength)
	for i := range ticketLength {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(ticketChars))))
		if err != nil {
			return "", err
		}
		b[i] = ticketChars[n.Int64()]
	}
	return string(b), nil
}

// GenerateUniqueTrackCode generates a unique ticket code by checking MongoDB.
func GenerateUniqueTrackCode(ctx context.Context, coll *mongo.Collection) (string, error) {
	for range maxAttempts {
		code, err := generateTrackCode()
		if err != nil {
			return "", err
		}

		// check if code already exists
		count, err := coll.CountDocuments(ctx, bson.M{"trackCode": code})
		if err != nil {
			return "", err
		}

		if count == 0 {
			return code, nil
		}
	}
	return "", errors.New("failed to generate unique track code code after max attempts")
}

// ParsTrackCode validates and returns a cleaned track code.
// It ensures the code has the correct length and allowed characters (A-Z, 0-9).
func ParsTrackCode(code string) (string, error) {
	code = strings.ToUpper(code)

	if len(code) != ticketLength {
		return "", errors.New("trackCode length is invalid")
	}

	match, err := regexp.MatchString("^[A-Z0-9]+$", code)
	if err != nil {
		return "", err
	}

	if !match {
		return "", errors.New("trackCode contains invalid characters")
	}

	return code, nil
}
