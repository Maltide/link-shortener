// Package handlers contains HTTP handlers for URL shortening and redirection.
package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/Maltide/link-shortener/pkg/hash"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// ShortHandler handles requests for URL shortening, creates short link, and returns it to the client.
func ShortHandler(db *sql.DB, log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Errorf("parseForm issue: %v", err)
			http.Error(w, "parse link error", http.StatusInternalServerError)
			return
		}

		linkIn := r.FormValue("url")

		if strings.HasPrefix(linkIn, "http://localhost:8080/") || strings.HasPrefix(linkIn, "https://localhost:8080/") {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte("Wrong request!"))
			if err != nil {
				log.Errorf("w.Write(after strings.HasPrefix) issue: %v", err)
				http.Error(w, "can't answer to you right now", http.StatusInternalServerError)
				return
			}
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
			_, err = w.Write([]byte(shortLink))
			if err != nil {
				log.Errorf("w.Write(after db.QueryRow) issue: %v", err)
				http.Error(w, "can't answer to you right now", http.StatusInternalServerError)
				return
			}
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

		shortLink, err = hash.Hash(id, log, db)
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

		_, err = w.Write([]byte(shortLink))
		if err != nil {
			log.Errorf("w.Write(when try to send shortLink) issue: %v", err)
			http.Error(w, "can't answer to you right now", http.StatusInternalServerError)
			return
		}
	}
}

// RedirectHandler handles redirection from short link to the original URL.
func RedirectHandler(db *sql.DB, log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
