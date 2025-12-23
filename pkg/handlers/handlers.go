package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/Maltide/link-shortener/pkg/hash"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func ShortHandler(db *sql.DB, log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		linkIn := r.FormValue("url")

		if strings.HasPrefix(linkIn, "http://localhost:8080/") || strings.HasPrefix(linkIn, "https://localhost:8080/") {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Wrong request!"))
			return
		}

		resp, err := http.Head(linkIn)
		if err != nil {
			log.Errorf("status check issue: %v", err)
			http.Error(w, "link is not valid", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Errorf("url not exist in IANA")
			http.Error(w, "link is not valid", http.StatusInternalServerError)
			resp.Body.Close()
			return
		} else {
			resp.Body.Close()
		}

		var id int
		var shortLink string

		err = db.QueryRow("SELECT hash_link FROM links WHERE original_link = $1", linkIn).Scan(&shortLink)
		if err == nil {
			w.Write([]byte(shortLink))
			return
		} else if err != sql.ErrNoRows {
			log.Errorf("dbQuery row issue: %v", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		err = db.QueryRow("INSERT INTO links (original_link) VALUES ($1) RETURNING id", linkIn).Scan(&id)
		if err != nil {
			log.Errorf("dbQuery insert issue: %v", err)
			http.Error(w, "DB insert error", http.StatusInternalServerError)
			return
		}

		shortLink, err = hash.Hash(id, log)
		if err != nil {
			log.Errorf("hash issue: %v", err)
			http.Error(w, "DB update error", http.StatusInternalServerError)
		}

		_, err = db.Exec("UPDATE links SET hash_link = $1 WHERE id = $2", shortLink, id)
		if err != nil {
			log.Errorf("dbQuery update issue: %v", err)
			http.Error(w, "DB update error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(shortLink))
	}
}

func RedirectHandler(db *sql.DB, log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// r.ParseForm()

		// linkIn := r.FormValue("url")

		hash := chi.URLParam(r, "hash")

		var origianlLink string

		err := db.QueryRow("SELECT original_link FROM links WHERE hash_link = $1", hash).Scan(&origianlLink)
		if err != nil {
			log.Errorf("dbQuery search issue: %v", err)
			http.Error(w, "this address doesn't exist", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, origianlLink, http.StatusFound)
	}
}
