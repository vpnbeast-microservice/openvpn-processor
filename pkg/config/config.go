package config

import (
	"os"
	"strconv"
)

// ConvertStringToInt converts string from environment variable to integer value
func ConvertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		i = 0
	}
	return i
}

// GetStringEnv gets the environment variable if provided, returns the default value if not provided
func GetStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// GetIntEnv gets the environment variable if provided, returns the default value if not provided
func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return ConvertStringToInt(value)
}
