// Package main is the entry point for the link shortener application.
package main

import (
	"database/sql"
	"time"

	"github.com/Maltide/link-shortener/pkg/config"
	"github.com/Maltide/link-shortener/pkg/helpers"
	"github.com/Maltide/link-shortener/pkg/logger"
	"github.com/Maltide/link-shortener/pkg/server"
	_ "github.com/lib/pq"
)

// main initializes configuration, logger, database connection, and starts the server.
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	log, err := logger.Logger("debug")
	if err != nil {
		return
	}

	connStr := helpers.ConnStr(cfg)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	time.Sleep(5 * time.Second)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err := helpers.InitDB(db); err != nil {
		log.Fatalf("InitDB: fail to initialize table: %v", err)
	}

	log.Infof("db was initialized with params")

	server.StartServer(db, log)

	log.Info("Server started")
}
