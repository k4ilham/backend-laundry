package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		// handle error or ignore if using system envs
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
