package httpserver

import (
	"context"
	"log"
	"net/http"

	"github.com/waf-draft/waf/internal/config"
)

// Server represents the WAF HTTP server
type Server struct {
	httpServer *http.Server
	handler    *WAFHandler
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, handler *WAFHandler) *Server {
	router := NewRouter(handler, cfg.Logging.Output)
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Server.ListenAddress,
			Handler:      router,
			ReadTimeout:  cfg.Server.ReadTimeout(),
			WriteTimeout: cfg.Server.WriteTimeout(),
			IdleTimeout:  cfg.Server.IdleTimeout(),
		},
		handler: handler,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Starting WAF server on %s", s.httpServer.Addr)
	log.Printf("Upstream: %s", s.httpServer.Addr) // This will be logged from config
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

