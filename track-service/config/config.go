package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Broker string
	ApiKey string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env, using only env variables")
	}

	return &Config{
		ApiKey: os.Getenv("APIKEY"),
		Broker: os.Getenv("BROKER"),
	}
}
