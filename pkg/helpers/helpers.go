// Package helpers provides utility functions for database: connection strings, initialization, and maintenance.
package helpers

import (
	"database/sql"
	"fmt"

	"github.com/Maltide/link-shortener/pkg/config"
	"go.uber.org/zap"
)

// ConnStr forms a database connection string based on the provided configuration.
func ConnStr(cfg config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
}

// InitDB creates the 'links' table in the database if it does not exist.
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

// ClearDB deletes all records from the 'links' table.
func ClearDB(db *sql.DB, log *zap.SugaredLogger) error {
	_, err := db.Exec("DELETE ALL FROM links")
	if err != nil {
		log.Errorf("clearDB fail: %v", err)
		return err
	}

	return err
}
