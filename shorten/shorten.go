package shorten

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type ShortenURL struct {
	URL string `json:"url"`
}

type shorter interface {
	generate() (string, error)
}

type storer interface {
	Save(ctx context.Context, shortenURL, rawURL string) error
}

type Handler struct {
	// db    *sql.DB
	short shorter
	store storer
}

func NewHandler(short shorter, store storer) *Handler {
	return &Handler{short: short, store: store}
}

func (handler *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var url ShortenURL
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shorten, err := handler.short.generate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := handler.store.Save(ctx, shorten, url.URL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"shorten": shorten,
	})
}

func NewRawURLHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shorten := r.PathValue("shorten")

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		row := db.QueryRowContext(ctx, "SELECT url FROM links WHERE code = ?", shorten)

		var url string
		if err := row.Scan(&url); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "not found shorten", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
