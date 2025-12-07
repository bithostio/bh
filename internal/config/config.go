package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"api_base_url"`
}

// DefaultConfigPath returns the default path for the config file
func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".bh", "config.json"), nil
}

func New() *Config {
	return &Config{
		BaseURL: ensureTrailingSlash(getBaseURL()),
	}
}

// Load reads the configuration from file or environment variables
func Load() (*Config, error) {
	// Check environment variable first
	if apiKey := os.Getenv("BH_API_KEY"); apiKey != "" {
		return &Config{
			APIKey:  apiKey,
			BaseURL: ensureTrailingSlash(getBaseURL()),
		}, nil
	}

	// Load from config file
	configPath, err := DefaultConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found. Run 'bh auth' to set up your API key")
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = getBaseURL()
	}

	cfg.BaseURL = ensureTrailingSlash(cfg.BaseURL)

	return &cfg, nil
}

// Save writes the configuration to the config file
func Save(cfg *Config) error {
	configPath, err := DefaultConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Write with secure permissions (owner read/write only)
	return os.WriteFile(configPath, data, 0600)
}

func getBaseURL() string {
	if url := os.Getenv("BH_API_URL"); url != "" {
		return url
	}
	return "https://dashboard.bithost.io/api/v1/"
}

// ensureTrailingSlash ensures the URL ends with a trailing slash
func ensureTrailingSlash(url string) string {
	if !strings.HasSuffix(url, "/") {
		return url + "/"
	}
	return url
}
