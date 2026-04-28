package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost/finance_manager?sslmode=disable"
	}

	startupCtx, startupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startupCancel()

	db, err := database.NewPostgresPool(startupCtx, dsn)
	if err != nil {
		log.Fatal("database connection failed", zap.Error(err))
	}
	defer db.Close()

	router := server.NewRouter(log, db).Engine()

	s := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("starting server", zap.String("port", port))

		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-stopCtx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
		if err := s.Close(); err != nil {
			log.Fatal("force shutdown failed", zap.Error(err))
		}
		return
	}

	log.Info("server stopped")
}
