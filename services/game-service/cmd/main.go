package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portfolio-game-service/internal/handlers"
	"portfolio-game-service/internal/services"
	"portfolio-game-service/pkg/config"
	"portfolio-game-service/pkg/logger"
	"portfolio-game-service/pkg/metrics"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)
	
	metrics.Init()
	gameService := services.NewGameService(log)
	gameHandler := handlers.NewGameHandler(gameService, log)

	r := mux.NewRouter()
	
	// Health endpoint
	r.HandleFunc("/health", gameHandler.Health).Methods("GET")
	
	// Game endpoints
	r.HandleFunc("/game/start", gameHandler.StartGame).Methods("POST")
	r.HandleFunc("/game/move", gameHandler.MakeMove).Methods("POST")
	r.HandleFunc("/game/status", gameHandler.GetStatus).Methods("GET")
	
	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.WithField("port", cfg.Port).Info("Starting game service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Server failed to start")
		}
	}()

	// Graceful shutdown
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