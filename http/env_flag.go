package http_client

import "os"

// getEnv retrieves variables from the environment and falls back
// to a passed fallback variable if it isn't set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
