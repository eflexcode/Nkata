package evn

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func InitEvn() error {
	err := godotenv.Load()
	return err
}

func GetString(fallback string, key string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func GetInt(fallback int, key string) int {

	value := os.Getenv(key)

	intVal, err := strconv.Atoi(value)

	if err != nil {
		return fallback
	}

	return intVal
}
