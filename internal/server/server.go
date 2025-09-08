package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/cunnati/store-service/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	mux    *http.ServeMux
	logger *slog.Logger
}

func New(logger *slog.Logger) *Server {
	mux := http.NewServeMux()
	s := &Server{mux: mux, logger: logger}
	// Only metrics by default; health should be set by the app to include dependencies
	mux.Handle("/metrics", promhttp.Handler())
	return s
}

func (s *Server) Start(port string) *http.Server {
	srv := &http.Server{Addr: ":" + port, Handler: s.mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = srv.ListenAndServe() }()
	return srv
}

func (s *Server) Shutdown(ctx context.Context, srv *http.Server) error { return srv.Shutdown(ctx) }

func (s *Server) MustStart(port string) *http.Server {
	srv := s.Start(port)
	fmt.Printf("listening on :%s\n", port)
	return srv
}

// Handle registers a handler with logging and metrics.
func (s *Server) Handle(route string, h http.Handler) {
	s.mux.Handle(route, s.wrap(route, h))
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (w *statusRecorder) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (s *Server) wrap(route string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rec, r)
		dur := time.Since(start)
		if s.logger != nil {
			s.logger.Info("http_request",
				slog.String("route", route),
				slog.String("method", r.Method),
				slog.Int("status", rec.status),
				slog.String("remote", r.RemoteAddr),
				slog.Int64("dur_ms", dur.Milliseconds()),
			)
		}
		metrics.HTTPRequests.WithLabelValues(route, strconv.Itoa(rec.status)).Inc()
	})
}
