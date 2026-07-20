package shorten

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ShortenURL struct {
	URL string `json:"url"`
}

type IfcShorter interface {
	generate() (string, error)
}

type IfdStorer interface {
	Save(ctx context.Context, shortenURL, rawURL string) error
}

type AbHandler struct {
	// db    *sql.DB
	fdShort IfcShorter
	fdStore IfdStorer
}

func CdNewHandler(short IfcShorter, store IfdStorer) *AbHandler {
	return &AbHandler{fdShort: short, fdStore: store}
}

func (rcHandler *AbHandler) ShortenX(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Shorten init")
	if r.Method != http.MethodPost {
		fmt.Println("Shorten check method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Shorten A")
	var url ShortenURL
	defer r.Body.Close()
	fmt.Println("Shorten before decode")
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		fmt.Println("Shorten error internal")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Shorten before generate ")
	shorten, err := rcHandler.fdShort.generate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Shorten with timeout")
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	fmt.Println("Shorten before save")
	if err := rcHandler.fdStore.Save(ctx, shorten, url.URL); err != nil {
		fmt.Println("Shorten save error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Shorten before response")
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
