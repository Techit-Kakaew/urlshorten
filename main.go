package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Techit-Kakaew/urlshorten/ginrouter"
	"github.com/Techit-Kakaew/urlshorten/shorten"
	"github.com/Techit-Kakaew/urlshorten/sqlite"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var mySigningKey = []byte("AllYourBase")

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	db, err := sql.Open("sqlite3", "./urlshortener.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		db.Close()
	}()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	storage := sqlite.NewStorage(db)
	handler := shorten.NewHandler(storage)
	redirectHandler := shorten.NewRedirectHandler(storage)

	r.POST("/shorten", ginrouter.NewHandler(handler.Handler))
	r.GET("/:shorturl", ginrouter.NewHandler(redirectHandler.Handler))

	// mux := http.NewServeMux()

	// mux.HandleFunc("/login", loginHandler)

	// redirectHandler := shorten.NewRedirectHandler(storage)

	// r.GET("/{shorturl}", ginrouter.NewHandler(redirectHandler.ServeHTTP))

	// mux.HandleFunc("/shorten", httprouter.NewHandler(handler.Handler))

	// redirectHandler := shorten.NewRedirectHandler(storage)
	// mux.HandleFunc("/{shorturl}", redirectHandler.ServeHTTP)

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + os.Getenv("PORT"),
	}

	log.Printf("api serving on %s...\n", os.Getenv("PORT"))
	srv.ListenAndServe()
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		Issuer:    "test",
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
}

func authenMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authen := r.Header.Get("Authentication")
		tokenString := authen[7:]

		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return mySigningKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}
