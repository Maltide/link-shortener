package config

import (
	"os"
)

type Config struct {
	LogLevel string
}

func LoadConfig() (Config, error) {

	loglevel := os.Getenv("LOG_LEVEL")

	return Config{
		LogLevel: loglevel,
	}, nil
}
