package main

import (
	"log/slog"

	"urlshortener/internal/server"
)

func main() {
	httpServer := server.NewServer()

	slog.Info("Server starting on: http://localhost:8080")
	err := httpServer.ListenAndServe()
	slog.Error(err.Error())
}
