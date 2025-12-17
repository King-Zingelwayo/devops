package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portfolio-frontend-service/internal/handlers"
	"portfolio-frontend-service/pkg/config"
	"portfolio-frontend-service/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	handler := handlers.NewFrontendHandler(cfg.GameServiceURL, log)

	r := mux.NewRouter()
	
	// Health endpoint
	r.HandleFunc("/health", handler.Health).Methods("GET")
	
	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	
	// Portfolio UI
	r.HandleFunc("/", handler.Index).Methods("GET")
	r.HandleFunc("/api/game/start", handler.ProxyStartGame).Methods("POST")
	r.HandleFunc("/api/game/move", handler.ProxyMove).Methods("POST")
	r.HandleFunc("/api/game/status", handler.ProxyStatus).Methods("GET")

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.WithField("port", cfg.Port).Info("Starting frontend service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}
	log.Info("Server exited")
}