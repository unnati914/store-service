package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/cunnati/store-service/internal/config"
	"github.com/cunnati/store-service/internal/db"
	"github.com/cunnati/store-service/internal/handlers"
	"github.com/cunnati/store-service/internal/logging"
	"github.com/cunnati/store-service/internal/metrics"
	"github.com/cunnati/store-service/internal/models"
	"github.com/cunnati/store-service/internal/server"
)

func run() int {
	cfg := config.FromEnv()
	log := logging.New()
	metrics.MustRegisterAll()

	pg, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Error("db_connect_error", "err", err)
		return 1
	}
	defer pg.Close()

	rds, err := db.ConnectRedis(cfg.RedisAddr)
	if err != nil {
		log.Error("redis_connect_error", "err", err)
		return 1
	}
	defer rds.Close()

	s := server.New(log)
	// Health handler checks DB and Redis
	s.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := pg.SQL.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"ok":false,"db":"down"}`))
			return
		}
		if err := rds.Client.Ping(ctx).Err(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"ok":false,"redis":"down"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))

	// Auth routes
	authHandler := &handlers.AuthHandler{Users: &models.UserStore{DB: pg.SQL}, JWTSecret: cfg.JWTSecret}
	s.Handle("/signup", http.HandlerFunc(authHandler.Signup))
	s.Handle("/login", http.HandlerFunc(authHandler.Login))

	srv := s.MustStart(cfg.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = s.Shutdown(ctxShutdown, srv)
	return 0
}
