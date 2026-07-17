package main

import (
	"context"
	"database/sql"
	"day3starter/shorten"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

// type database struct {
// 	mux  sync.Mutex
// 	data map[string]string
// }

// var db = database{
// 	mux:  sync.Mutex{},
// 	data: make(map[string]string),
// }

const schema = `CREATE TABLE IF NOT EXISTS links ( code TEXT PRIMARY KEY, url TEXT NOT NULL );`

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := sql.Open("sqlite", "./urlshorten.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if _, err := db.ExecContext(ctx, schema); err != nil {
			log.Fatal(err)

		}
	}

	mux := http.NewServeMux()
	mux.Handle("/ping", http.HandlerFunc(pingpongHandler))

	store := shorten.NewStore(db)
	short := shorten.NewShorter()

	handler := shorten.NewHandler(short, store)
	mySigningKey := []byte("AllYourBase")

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {

		// Create the Claims
		claims := &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString(mySigningKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"token": ss,
		})
	})

	mux.HandleFunc("/shorten", AuthenMiddleware(handler.Shorten, mySigningKey))
	mux.HandleFunc("/{shorten}", shorten.NewRawURLHandler(db))

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		fmt.Println("shutting down...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	fmt.Println("serve on :" + os.Getenv("PORT"))
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
	}

	fmt.Println("gracefully")
}

func AuthenMiddleware(handler http.HandlerFunc, signingKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) < 7 {
			http.Error(w, "ไม่บอก", http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[7:]

		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return signingKey, nil
		})
		if err != nil {
			http.Error(w, "ไม่บอก", http.StatusUnauthorized)
			return
		}

		handler(w, r)

		// after
	}
}

func pingpongHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{
		"message": "pong",
	}

	b, err := json.Marshal(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
	w.WriteHeader(http.StatusOK)
}
