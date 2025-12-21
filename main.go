package main

import (
	"database/sql"

	"github.com/Maltide/link-shortener/helpers"
	"github.com/Maltide/link-shortener/pkg/config"
	"github.com/Maltide/link-shortener/pkg/logger"
	_ "github.com/lib/pq"
)

// import (
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Custom server configuration!")
// }

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

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	log.Infof("db was inizialised with params:")

	// server := &http.Server{
	// 	Addr:         ":8080",
	// 	Handler:      http.HandlerFunc(handler),
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// }

	// fmt.Println("Starting custom server at port 8080")
	// err := server.ListenAndServe()
	// if err != nil {
	// 	fmt.Println("Error starting the server:", err)
	// }
}
