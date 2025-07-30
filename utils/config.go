package utils

import (
	"os"
	"strconv"
)

// GetEnv membaca environment variable sebagai string, dengan fallback ke nilai default.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvAsInt membaca environment variable sebagai integer, dengan fallback.
func GetEnvAsInt(key string, fallback int) int {
	valueStr := GetEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

