package options

import (
	"go.uber.org/zap"
	"os"
	"strconv"
)

// getStringEnv gets the specific environment variables with default value, returns default value if variable not set
func getStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// getIntEnv gets the specific environment variables with default value, returns default value if variable not set
func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return ConvertStringToInt(value)
}

// ConvertStringToInt converts string environment variables to integer values
func ConvertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Warn("an error occurred while converting from string to int. Setting it as zero",
			zap.String("error", err.Error()))
		i = 0
	}
	return i
}
