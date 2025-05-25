package main

import (
	"log/slog"
	"net/http"

	"urlshortener/internal/router"
)

func RegisterRouter() *router.Router {
	r := router.NewRouter()

	r.Group(func(r *router.Router) {
		r.HandleFunc("GET /", helloHandler)
	})

	return r
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: RegisterRouter(),
	}

	slog.Info("Server starting on: http://localhost:8080")
	err := server.ListenAndServe()
	slog.Error(err.Error())
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, World"))
}
