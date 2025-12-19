package config

import (
	"os"
)

type Config struct {
	TGToken   string
	LogLevel  string
	TGTimeout int
}

func LoadConfig() (Config, error) {

	loglevel := os.Getenv("LOG_LEVEL")

	return Config{
		LogLevel: loglevel,
	}, nil
}
