package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN     string
	AdminId string
	Broker  string
	Tracker string
}

func Load() *Config {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env, using only env variables")
	}

	return &Config{
		DSN:     os.Getenv("DSN"),
		AdminId: os.Getenv("ADMIN"),
		Broker:  os.Getenv("BROKER"),
		Tracker: os.Getenv("TRACKER"),
	}
}
