package server

import (
	"net/http"

	"urlshortener/internal/router"

	"github.com/redis/go-redis/v9"
)

type State struct {
	redisClient *redis.Client
	newKey      string
}

func (s *State) RegisterRouter() *router.Router {
	r := router.NewRouter()

	r.Group(func(r *router.Router) {
		r.Use(Recoverer)
		r.HandleFunc("GET /{identifier}", s.GetLongUrl)
		r.HandleFunc("POST /", s.ShortenUrl)
		r.HandleFunc("GET /latest", s.Latest)
	})

	return r
}

func NewServer() *http.Server {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	state := &State{
		redisClient: client,
		newKey:      "",
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: state.RegisterRouter(),
	}

	return server
}
