// Package server sets up and starts the HTTP server with routing.
package server

import (
	"database/sql"
	"net/http"

	"github.com/Maltide/link-shortener/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// StartServer launches the HTTP server and sets up routing for requests.
func StartServer(db *sql.DB, log *zap.SugaredLogger) {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/web.html")
	})
	r.Get("/{hash}", handlers.RedirectHandler(db, log))
	r.Post("/short", handlers.ShortHandler(db, log))
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Errorf("fail to start server: %v", err)
		return
	}
}
