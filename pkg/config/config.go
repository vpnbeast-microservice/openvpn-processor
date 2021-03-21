package config

import (
	"log"
	"os"
	"strconv"
)

func ConvertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("An error occured while converting %s to int. Setting it as zero.", s)
		i = 0
	}
	return i
}

func GetStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return ConvertStringToInt(value)
}