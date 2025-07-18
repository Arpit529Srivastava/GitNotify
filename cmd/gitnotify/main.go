package main

import (
	"flag"
	"log"

	"gitnotify/internal/config"
	"gitnotify/internal/server"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	log.Printf("Loading configuration from: %s", *configPath)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.ValidateConfig(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully")
	log.Printf("Organization: %s", cfg.Organization)
	log.Printf("Port: %d", cfg.Port)
	log.Printf("Notification rules: %d", len(cfg.Notifications))

	// Create and start server
	srv := server.NewServer(cfg)

	// Handle graceful shutdown
	go func() {
		// TODO: Add signal handling for graceful shutdown
		// For now, just log that the server is running
		log.Printf("GitNotify server is running. Press Ctrl+C to stop.")
	}()

	// Start the server
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
