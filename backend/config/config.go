package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MONGO_URI  string
	PORT       string
	DB_NAME    string
	JWT_SECRET string
}

func LoadConfig() Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return Config{
		MONGO_URI:  os.Getenv("MONGO_URI"),
		PORT:       os.Getenv("PORT"),
		DB_NAME:    os.Getenv("DB_NAME"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
	}
}
