package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boonyarit-iamsaard/online-shop/internal/config"
	"github.com/boonyarit-iamsaard/online-shop/internal/database"
	"github.com/boonyarit-iamsaard/online-shop/internal/logger"
	"github.com/boonyarit-iamsaard/online-shop/internal/server"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer func(log *zap.Logger) {
		_ = log.Sync()
	}(log)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config failed", zap.Error(err))
	}

	startupCtx, startupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startupCancel()

	db, err := database.NewPostgresPool(startupCtx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("database connection failed", zap.Error(err))
	}
	defer db.Close()

	router := server.NewRouter(log, db).Engine()

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("starting server", zap.String("port", cfg.Port))

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-stopCtx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
		if err := httpServer.Close(); err != nil {
			log.Fatal("force shutdown failed", zap.Error(err))
		}
		return
	}

	log.Info("server stopped")
}
