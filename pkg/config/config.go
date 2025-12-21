package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel string
	DBUser   string
	DBPass   string
	DBName   string
	DBHost   string
	DBPort   string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load() // loads .env file into os.Getenv
	if err != nil {
		return Config{}, fmt.Errorf("config: error loading .env file: %v", err)
	}

	return Config{
		LogLevel: os.Getenv("LOG_LEVEL"),
		DBUser:   os.Getenv("DB_USER"),
		DBPass:   os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		DBHost:   os.Getenv("DB_HOST"),
		DBPort:   os.Getenv("DB_PORT"),
	}, nil
}
