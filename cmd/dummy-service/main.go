package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zhenisduissekov/another-dummy-service/internal/config"
	"github.com/zhenisduissekov/another-dummy-service/internal/repository/inmem"
	"github.com/zhenisduissekov/another-dummy-service/internal/services"
	"github.com/zhenisduissekov/another-dummy-service/internal/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := app(); err != nil {
		fmt.Printf("Could not run app: %v\n", err)
	}
}

func app() error {
	cfg := config.New()

	repo := inmem.NewPortStore()

	service := services.NewService(repo)

	server := transport.NewHttpServer(service)

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("health OK")
	}).Methods(http.MethodGet)
	router.HandleFunc("/port", server.GetPort).Methods(http.MethodGet)
	router.HandleFunc("/count", server.CountPorts).Methods(http.MethodGet)
	router.HandleFunc("/ports", server.UploadPorts).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	stopped := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}

		close(stopped)
	}()

	log.Printf("Starting HTTP server on %s", cfg.Port)

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Println("Server has been stopped")

	return nil
}
