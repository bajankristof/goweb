package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
	ShutdownTimeout time.Duration
}

type ServerOption func(*Server)

func New(handler http.Handler, opts ...ServerOption) *Server {
	srv := &Server{
		Server: &http.Server{
			Handler: handler,
		},
		ShutdownTimeout: 15 * time.Second,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	errs := make(chan error, 1)
	go func() { errs <- s.Server.ListenAndServe() }()

	slog.InfoContext(ctx, "server started", "addr", s.Addr)

	select {
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	case err := <-errs:
		slog.ErrorContext(ctx, "server error", "addr", s.Addr, "err", err)
		return err
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.ShutdownTimeout)
	defer cancel()
	return s.Server.Shutdown(ctx)
}

func WithAddr(addr string) ServerOption {
	return func(s *Server) {
		s.Addr = addr
	}
}

func WithReadTimeout(d time.Duration) ServerOption {
	return func(s *Server) {
		s.ReadTimeout = d
	}
}

func WithWriteTimeout(d time.Duration) ServerOption {
	return func(s *Server) {
		s.WriteTimeout = d
	}
}

func WithIdleTimeout(d time.Duration) ServerOption {
	return func(s *Server) {
		s.IdleTimeout = d
	}
}

func WithShutdownTimeout(d time.Duration) ServerOption {
	return func(s *Server) {
		s.ShutdownTimeout = d
	}
}
