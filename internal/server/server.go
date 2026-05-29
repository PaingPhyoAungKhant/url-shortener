// Package server provides the implementation of the HTTP server for the URL Shortener application.
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/config"
	"github.com/PaingPhyoAungKhant/url-shortener/internal/handler"
	"github.com/PaingPhyoAungKhant/url-shortener/internal/logger"
	"github.com/PaingPhyoAungKhant/url-shortener/internal/middleware"
)

// Server represents the HTTP server
type Server struct {
	http *http.Server
	cfg  *config.Config
	log  *logger.Logger
	mux  *http.ServeMux
}

// Group represents a route group
type Group struct {
	server      *Server
	prefix      string
	middlewares []func(http.Handler) http.Handler
}

// New creates a new instance of the Server with the necessary routes and handlers.
func New(cfg *config.Config, log *logger.Logger) *Server {
	mux := http.NewServeMux()
	h := handler.New()

	wrappedMux := chain(
		mux,
		middleware.RequestID,
		middleware.Error,
		middleware.Logging(log),
		middleware.Recovery(log),
	)

	s := &Server{
		http: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      wrappedMux,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		mux: mux,
		cfg: cfg,
		log: log,
	}

	s.Handle("GET /health", http.HandlerFunc(h.Health))

	return s
}

// Group creates a new route group with given prefix name and middlewares
func (s *Server) Group(prefix string, middlewares ...func(http.Handler) http.Handler) *Group {
	return &Group{
		server:      s,
		prefix:      prefix,
		middlewares: middlewares,
	}
}

// Handle register a handler to the route group with given pattern and optional additional middlewares provided
func (g *Group) Handle(pattern string, h http.Handler, middlewares ...func(http.Handler) http.Handler) {
	combined := append(g.middlewares, middlewares...)
	g.server.Handle(g.prefix+pattern, h, combined...)
}

// chain is a helper function for middleware chaining
func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// Handle register a route handler with provided middlewares to a given route pattern
func (s *Server) Handle(pattern string, h http.Handler, middlewares ...func(http.Handler) http.Handler) {
	s.mux.Handle(pattern, chain(h, middlewares...))
}

// Start starts the HTTP server and listens for incoming requests.
func (s *Server) Start(ctx context.Context) error {
	listenErr := make(chan error, 1)

	go func() {
		log.Printf("Server is listening on %s \n", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			listenErr <- err
		}
		close(listenErr)
	}()

	select {
	case err := <-listenErr:
		return err
	case <-ctx.Done():
	}

	shutDownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.cfg.Server.ShutdownTimeout)
	defer cancel()
	if err := s.http.Shutdown(shutDownCtx); err != nil {
		return fmt.Errorf("server shudtdown failed: %w", err)
	}
	if err := <-listenErr; err != nil {
		return err
	}

	log.Println("server gracefully stopped")

	return nil
}

// Addr returns the address the server is listening on.
func (s *Server) Addr() string {
	return s.http.Addr
}

// ReadTimeout returns the read timeout duration for the server.
func (s *Server) ReadTimeout() time.Duration {
	return s.http.ReadTimeout
}

// WriteTimeout returns the write timeout duration for the server.
func (s *Server) WriteTimeout() time.Duration {
	return s.http.WriteTimeout
}

func (s *Server) IdleTimeout() time.Duration {
	return s.http.IdleTimeout
}
