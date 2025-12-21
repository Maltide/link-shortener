package helpers

import (
	"fmt"

	"github.com/Maltide/link-shortener/pkg/config"
)

func ConnStr(cfg config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
}
