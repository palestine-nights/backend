package main

import (
	"os"
)

// Get environment variable with ability to specify default value.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		value = fallback
	}

	return value
}
