package zutils

import (
	"os"
	"strconv"
	"strings"
)

// GetEnv returns the value of the environment variable named by the key and returns the default value if the environment variable is not set.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvToLower(key, fallback string) string {
	return strings.ToLower(GetEnv(key, fallback))
}

func GetEnvToUpper(key, fallback string) string {
	return strings.ToUpper(GetEnv(key, fallback))
}

func GetEnvToBool(key, fallback string) bool {
	enabled := false
	if v := GetEnv(key, fallback); v != "" {
		enabled, _ = strconv.ParseBool(v)
	}
	return enabled
}
