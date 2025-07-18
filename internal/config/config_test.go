package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `
organization: "test-org"
port: 9090
webhook_secret: "test-secret"
notifications:
  - event_type: "issues"
    actions: ["opened", "closed"]
  - event_type: "pull_request"
    actions: ["opened"]
`

	tmpfile, err := os.CreateTemp("", "config_test_*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load the config
	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify the loaded values
	if config.Organization != "test-org" {
		t.Errorf("Expected organization 'test-org', got '%s'", config.Organization)
	}

	if config.Port != 9090 {
		t.Errorf("Expected port 9090, got %d", config.Port)
	}

	if config.WebhookSecret != "test-secret" {
		t.Errorf("Expected webhook secret 'test-secret', got '%s'", config.WebhookSecret)
	}

	if len(config.Notifications) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(config.Notifications))
	}

	// Test validation
	if err := config.ValidateConfig(); err != nil {
		t.Errorf("Config validation failed: %v", err)
	}
}

func TestLoadConfigDefaultPort(t *testing.T) {
	// Create a config file without port
	configContent := `
organization: "test-org"
webhook_secret: "test-secret"
`

	tmpfile, err := os.CreateTemp("", "config_test_*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load the config
	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify default port is set
	if config.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", config.Port)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Organization:  "test-org",
				Port:          8080,
				WebhookSecret: "test-secret",
			},
			wantErr: false,
		},
		{
			name: "missing organization",
			config: &Config{
				Port:          8080,
				WebhookSecret: "test-secret",
			},
			wantErr: true,
		},
		{
			name: "missing webhook secret",
			config: &Config{
				Organization: "test-org",
				Port:         8080,
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			config: &Config{
				Organization:  "test-org",
				Port:          0,
				WebhookSecret: "test-secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.ValidateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
