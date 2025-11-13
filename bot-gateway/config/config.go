package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Broker string
	Core   string
	Token  string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env, using only env variables")
	}

	return &Config{
		Broker: os.Getenv("BROKER"),
		Core:   os.Getenv("CORE"),
		Token:  os.Getenv("TOKEN"),
	}
}
