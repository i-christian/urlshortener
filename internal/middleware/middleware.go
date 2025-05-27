package middleware

import (
	"encoding/json"
	"log"
	"net/http"
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
