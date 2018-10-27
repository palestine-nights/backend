package tools

import (
	"os"
)

// GetEnv returns environment variable with ability to specify default value.
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		value = fallback
	}

	return value
}
