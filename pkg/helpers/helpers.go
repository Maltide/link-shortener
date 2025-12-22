package helpers

import (
	"database/sql"
	"fmt"

	"github.com/Maltide/link-shortener/pkg/config"
)

func ConnStr(cfg config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
}

func InitDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS links (
	id SERIAL PRIMARY KEY,
	original_link TEXT UNIQUE NOT NULL,
	hash_link TEXT UNIQUE
	);`

	_, err := db.Exec(query)

	return err
}
