package util

import (
	"errors"
	"regexp"
	"strings"
)

// parseObjectName validates an object name for MinIO
// Returns the object name itself if valid, or an error if invalid
func ParseObjectName(name string) (string, error) {
	if name == "" {
		return "", errors.New("object name is empty")
	}

	// Reject path traversal
	if strings.Contains(name, "..") {
		return "", errors.New("invalid object name: contains '..'")
	}

	// Reject backslashes
	if strings.Contains(name, "\\") {
		return "", errors.New("invalid object name: contains '\\'")
	}

	// Allow only safe characters: a-zA-Z0-9 - _ . /
	validKey := regexp.MustCompile(`^[a-zA-Z0-9\-_.\/]+$`)
	if !validKey.MatchString(name) {
		return "", errors.New("invalid object name: contains unsafe characters")
	}

	return name, nil
}

// parseObjectNames validates a slice of object names for MinIO
// Returns a slice of valid names, or an error if any name is invalid
func ParseObjectNames(names []string) ([]string, error) {
	validNames := make([]string, 0, len(names))

	for _, name := range names {
		parsed, err := ParseObjectName(name) // reuse single object validation
		if err != nil {
			return nil, errors.New("invalid object name '" + name + "': " + err.Error())
		}
		validNames = append(validNames, parsed)
	}

	return validNames, nil
}
