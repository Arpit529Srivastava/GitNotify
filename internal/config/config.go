package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure for GitNotify
type Config struct {
	Organization  string         `yaml:"organization"`
	Port          int            `yaml:"port"`
	WebhookSecret string         `yaml:"webhook_secret"`
	Notifications []Notification `yaml:"notifications"`
	GitHubApp     GitHubApp      `yaml:"github_app"`
}

// Notification represents a notification rule
type Notification struct {
	EventType string   `yaml:"event_type"` // issues, pull_request, push
	Actions   []string `yaml:"actions"`    // opened, closed, reopened, etc.
	Repos     []string `yaml:"repos"`      // specific repos to monitor (optional)
}

// GitHubApp represents GitHub App authentication settings
type GitHubApp struct {
	AppID          int    `yaml:"app_id"`
	InstallationID int    `yaml:"installation_id"`
	PrivateKeyPath string `yaml:"private_key_path"`
}

// LoadConfig reads the configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set default port if not specified
	if config.Port == 0 {
		config.Port = 8080
	}

	return &config, nil
}

// ValidateConfig checks if the configuration is valid
func (c *Config) ValidateConfig() error {
	if c.Organization == "" {
		return fmt.Errorf("organization is required")
	}
	if c.WebhookSecret == "" {
		return fmt.Errorf("webhook_secret is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}
