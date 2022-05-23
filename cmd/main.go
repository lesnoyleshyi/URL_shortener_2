package main

import (
	"URL_shortener_2/internal/handlers"
	"URL_shortener_2/internal/repository"
	"URL_shortener_2/internal/services"
	"errors"
	"flag"
	"log"
	"net/http"
)

func main() {
	var storageType string
	flag.StringVar(&storageType, "storageType", "pg",
		"choose storageType to save URLs: pg for postgres, cache for in-memory")
	flag.Parse()

	repo := repository.New(storageType)
	service := services.New(&repo)
	handler := handlers.New(service)
	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed unexpectidly with error: %s", err)
	}
}