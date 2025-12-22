package server

import (
	"database/sql"
	"net/http"

	"github.com/Maltide/link-shortener/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func StartServer(db *sql.DB, log *zap.SugaredLogger) {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/web.html")
	})
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Post("/short", handlers.ShortHandler(db, log))
	r.Get("/{hash}", handlers.RedirectHandler(db, log))

	http.ListenAndServe(":8080", r)
}
