package main

import (
	"URL_shortener_2/internal/handlers"
	"URL_shortener_2/internal/repository"
	"URL_shortener_2/internal/services"
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var storageType string
	flag.StringVar(&storageType, "storage", "pg",
		"choose storage to save URLs: pg for postgres, cache for in-memory")
	flag.Parse()

	repo := repository.New(storageType)
	service := services.New(&repo)
	handler := handlers.New(service)
	server := http.Server{Addr: ":8080", Handler: handler}

	done := make(chan struct{})

	go start(&server)
	shutdown(&server, done)

	<-done
	log.Println("Server stopped gracefully")
}

func start(srv *http.Server) {
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed unexpectidly with error: %s", err)
	}
}

func shutdown(srv *http.Server, done chan struct{}) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-ch
	log.Printf("%s signal caught, shutting down gracefully", sig)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		cancel()
		close(done)
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %s", err)
	}
}
