package tools

import (
	"os"
	"regexp"
)

// GetEnv returns environment variable with ability to specify default value.
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		value = fallback
	}

	return value
}

// ValidateEmail validates email.
func ValidateEmail(email string) bool {
	result, err := regexp.MatchString(`.+@.+\..+`, email)

	if err != nil {
		panic(err)
	}

	return result
}
