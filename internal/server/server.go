package server

import (
	"net/http"

	"urlshortener/internal/handlers"
	"urlshortener/internal/middleware"
	"urlshortener/internal/router"
)

func RegisterRouter() *router.Router {
	r := router.NewRouter()

	r.Group(func(r *router.Router) {
		r.Use(middleware.Recoverer)
		r.HandleFunc("GET /", handlers.HelloHandler)
		r.HandleFunc("GET /test", nil)
	})

	return r
}

func NewServer() *http.Server {
	server := &http.Server{
		Addr:    ":8080",
		Handler: RegisterRouter(),
	}

	return server
}
