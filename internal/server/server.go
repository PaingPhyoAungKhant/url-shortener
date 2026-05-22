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
)

// Server represents the HTTP server
type Server struct {
	http *http.Server
	cfg  *config.Config
}

// New creates a new instance of the Server with the necessary routes and handlers.
func New(cfg *config.Config) *Server {
	mux := http.NewServeMux()
	h := handler.New()
	mux.HandleFunc("GET /health", h.Health)

	return &Server{
		http: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      mux,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		cfg: cfg,
	}
}

func (s *Server) Router() *http.ServeMux {
	return s.http.Handler.(*http.ServeMux)
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
