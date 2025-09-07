package util

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

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

