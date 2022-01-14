package zutils

import "os"

// GetEnv returns the value of the environment variable named by the key and returns the default value if the environment variable is not set.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
