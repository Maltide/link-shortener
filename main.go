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
		log.Fatalf("InitDB: fail to inizialize table: %v", err)
	}

	log.Infof("db was inizialised with params")

	server.StartServer(db, log)

	log.Info("Server started")
}
