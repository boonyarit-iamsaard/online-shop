package server

import (
	"net/http"
	"time"

	"github.com/boonyarit-iamsaard/finance-manager/internal/health"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Router struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func NewRouter(log *zap.Logger, db *pgxpool.Pool) *Router {
	return &Router{
		log: log,
		db:  db,
	}
}

func (r *Router) Engine() *gin.Engine {
	engine := gin.New()

	engine.Use(r.requestLogger())
	engine.Use(r.recovery())

	healthHandler := health.NewHandler(r.db)

	engine.GET("/healthz", healthHandler.Liveness)
	engine.GET("/readyz", healthHandler.Readiness)

	return engine
}

func (r *Router) requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		r.log.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("duration", time.Since(start)),
		)
	}
}

func (r *Router) recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		r.log.Error("panic recovered",
			zap.Any("error", recovered),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
