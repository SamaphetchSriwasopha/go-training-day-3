package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := sql.Open("sqlite3", "./urlshorten.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.Handle("/ping", http.HandlerFunc(pingpongHandler))
	mux.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println((r.PathValue("id")))
	})

	// handler := shorten.NewHandler(db)
	// mux.HandleFunc("/shorten", http.HandlerFunc(handler.ShortenDb))

	// mux.HandleFunc("/{short}", http.HandlerFunc(shortenRedirectDb))

	mux.Handle("/token", http.HandlerFunc(tokenHandler))

	http.ListenAndServe(":8080", mux)
}

// type handler struct{}
func pingpongHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "hello, world")
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

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	mySigningKey := []byte("AllYourBase")
	// Create the Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Issuer:    "local",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	// w.Write(b)
	// w.WriteHeader(http.StatusOK)
	fmt.Println(ss, err)

}

func chanel() {
	ch := make(chan int)
	go fibonacci(ch)
	for range 20 {
		fmt.Print(<-ch, ",")
	}
}

func fibonacci(ch chan int) {
	a, b := 0, 1
	for {
		ch <- a
		a, b = b, a+b
	}
}
