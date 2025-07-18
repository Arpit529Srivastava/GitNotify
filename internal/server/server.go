package server

import (
	"fmt"
	"log"
	"net/http"

	"gitnotify/internal/config"
	"gitnotify/internal/github"
)

// Server represents the HTTP server
type Server struct {
	config  *config.Config
	handler *github.Handler
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config) *Server {
	return &Server{
		config:  cfg,
		handler: github.NewHandler(cfg),
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Set up routes
	http.HandleFunc("/webhook", s.handler.HandleWebhook)

	// Add a health check endpoint
	http.HandleFunc("/health", s.handleHealth)

	// Start the server
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("Starting GitNotify server on port %d", s.config.Port)
	log.Printf("Webhook endpoint: http://localhost%s/webhook", addr)
	log.Printf("Health check: http://localhost%s/health", addr)

	return http.ListenAndServe(addr, nil)
}

// handleHealth provides a simple health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"gitnotify"}`))
}
