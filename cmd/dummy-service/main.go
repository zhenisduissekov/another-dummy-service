package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/zhenisduissekov/another-dummy-service/internal/config"
	"github.com/zhenisduissekov/another-dummy-service/internal/repository/inmem"
	"github.com/zhenisduissekov/another-dummy-service/internal/services"
	"github.com/zhenisduissekov/another-dummy-service/internal/transport"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Could not run app: %v\n", err)
	}
}

func run() error {
	// read config from env
	cfg := config.Read()

	// create port repository
	portStoreRepo := inmem.NewPortStore()

	// create port service
	portService := services.NewPortService(portStoreRepo)

	// create http server with application injected
	httpServer := transport.NewHttpServer(portService)

	// create http router
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("health OK")
	}).Methods(http.MethodGet)
	router.HandleFunc("/port", httpServer.GetPort).Methods(http.MethodGet)
	router.HandleFunc("/count", httpServer.CountPorts).Methods(http.MethodGet)
	router.HandleFunc("/ports", httpServer.UploadPorts).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	// listen to OS signals and gracefully shutdown HTTP server
	stopped := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
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
