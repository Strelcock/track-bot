package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN     string
	AdminId string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	return &Config{

		DSN:     os.Getenv("DSN"),
		AdminId: os.Getenv("ADMIN"),
	}
}
