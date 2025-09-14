package util

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidatePhoneNumber is a custom validator function
func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	// The regex to validate a string starting with "09" followed by 9 digits
	re := regexp.MustCompile(`^09\d{9}$`)
	return re.MatchString(fl.Field().String())
}
