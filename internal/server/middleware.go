package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Recoverer function handles goroutines panics for each request.
// Prevents the shutdown of the whole server from a panic from any goroutine's request.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				log.Printf("ERROR %v\n", err)

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// Logger logs the incoming HTTP request and its duration
func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request details
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(w, r)

		// Log response details
		log.Printf("Response: %s %s - %v", r.Method, r.URL.Path, time.Since(start))
	}

	return http.HandlerFunc(fn)
}
