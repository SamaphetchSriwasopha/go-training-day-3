package shorten

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

type ShortenURL struct {
	URL string `json:"url"`
}
type database struct {
	mux  sync.Mutex
	data map[string]string
}

var db = database{
	mux:  sync.Mutex{},
	data: make(map[string]string),
}

const schema = `CREATE TABLE IF NOT EXISTS links ( code TEXT PRIMARY KEY, url TEXT NOT NULL );`

func (handler *Handler) ShortenDb(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := handler.db.ExecContext(ctx, schema); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var url ShortenURL
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println((r.PathValue("id")))
	if err := json.Unmarshal(b, &url.URL); err != nil {
		http.Error(w, err.Error(), http.StatusInsufficientStorage)
		return
	}

	shorten, err := GetShortenURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInsufficientStorage)
		return
	}

	ctx, cancel = context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if _, err := handler.db.ExecContext(ctx, `INSERT INTO links (code, url) VALUES(?,?)`, shorten, url.URL); err != nil {

	}

	// db.mux.Lock()
	// db.data[shorten] = url.URL
	// db.mux.Unlock()

	w.Header().Set("Content-Type", "application/json")
	// jsonStr := fmt.Sprintf(`{"message":"hello %s"}`, name)
	// w.Write([]byte("{\"a\":\"1\"}"))
	json.NewEncoder(w).Encode(map[string]string{
		"shorthen": shorten,
	})
}

// func shortenRedirectDb(w http.Response, r *http.Request) {
// 	shorten := r.PathValue(("shorten"))
// 	if row, ok := db.data[shorten]; ok {
// 		delete(db.data, shorten)
// 		http.Redirect(w, r, row, http.StatusNotFound)
// 		return
// 	}
// 	w.WriteHeader(http.StatusFound)
// }

func ShortenMap(w http.ResponseWriter, r *http.Request) {
	var url ShortenURL
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println((r.PathValue("id")))
	if err := json.Unmarshal(b, &url.URL); err != nil {
		http.Error(w, err.Error(), http.StatusInsufficientStorage)
		return
	}

	shorten, err := GetShortenURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInsufficientStorage)
		return
	}
	db.mux.Lock()
	db.data[shorten] = url.URL
	db.mux.Unlock()

	w.Header().Set("Content-Type", "application/json")
	// jsonStr := fmt.Sprintf(`{"message":"hello %s"}`, name)
	// w.Write([]byte("{\"a\":\"1\"}"))
	json.NewEncoder(w).Encode(map[string]string{
		"shorthen": shorten,
	})
}

func ShortenCode(w http.ResponseWriter, r *http.Request) {
	shorten := r.PathValue(("shorten"))
	if row, ok := db.data[shorten]; ok {
		delete(db.data, shorten)
		http.Redirect(w, r, row, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusFound)
}

func GetShortenURL() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil // e.g. "aZ3xQ1"
}

// เก็บ mapping short -> url ]' map [string] string ก่อนยังไม่มี db
