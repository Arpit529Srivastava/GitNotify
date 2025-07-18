package server

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	"gitnotify/internal/config"
	"gitnotify/internal/github"
	"os"

	"gopkg.in/yaml.v3"
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

	// Add config management API endpoints
	http.HandleFunc("/api/config", s.handleConfigAPI)

	// Start the server
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("Starting GitNotify server on port %d", s.config.Port)
	log.Printf("Webhook endpoint: http://localhost%s/webhook", addr)
	log.Printf("Health check: http://localhost%s/health", addr)
	log.Printf("Config API: http://localhost%s/api/config", addr)

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

// handleConfigAPI handles GET and PUT for /api/config with token authentication
func (s *Server) handleConfigAPI(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("GITNOTIFY_CONFIG_TOKEN")
	if token == "" {
		http.Error(w, "Config API token not set", http.StatusInternalServerError)
		return
	}
	if !checkTokenAuth(r, token) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.config)
	case http.MethodPut:
		var newConfig config.Config
		if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid JSON"}`))
			return
		}
		if err := newConfig.ValidateConfig(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid config: ` + err.Error() + `"}`))
			return
		}
		// Save config to file (assume config.yml)
		if err := saveConfigToFile("config.yml", &newConfig); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to save config: ` + err.Error() + `"}`))
			return
		}
		// Reload config in server
		s.config = &newConfig
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	default:
		w.Header().Set("Allow", "GET, PUT")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// checkTokenAuth checks for Bearer token in Authorization header
func checkTokenAuth(r *http.Request, token string) bool {
	header := r.Header.Get("Authorization")
	if header == "Bearer "+token {
		return true
	}
	return false
}

// saveConfigToFile saves the config as YAML to the given path
func saveConfigToFile(path string, cfg *config.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
