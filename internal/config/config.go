package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server ServerConfig `yaml:"server"`
	CORS   CORSConfig   `yaml:"cors"`
	Domain DomainConfig `yaml:"domain"`
	Log    LogConfig    `yaml:"logging"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port         string        `yaml:"port"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

// DomainConfig represents domain checking configuration
type DomainConfig struct {
	ExtensionsFile      string        `yaml:"extensions_file"`
	Timeout             time.Duration `yaml:"timeout"`
	MaxConcurrentChecks int           `yaml:"max_concurrent_checks"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

var globalConfig *Config

// Load loads configuration from file
func Load(configPath string) (*Config, error) {
	// Set default config path if empty
	if configPath == "" {
		configPath = "./configs/config.yaml"
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get returns the global configuration
func Get() *Config {
	if globalConfig == nil {
		// Load default config if not loaded
		defaultConfig := &Config{
			Server: ServerConfig{
				Port:         ":8080",
				Host:         "localhost",
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
			},
			CORS: CORSConfig{
				AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
			},
			Domain: DomainConfig{
				ExtensionsFile:      "./data/domain_extensions.txt",
				Timeout:             5 * time.Second,
				MaxConcurrentChecks: 10,
			},
			Log: LogConfig{
				Level:  "info",
				Format: "json",
			},
		}
		globalConfig = defaultConfig
	}
	return globalConfig
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	if cfg.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if cfg.Domain.ExtensionsFile == "" {
		return fmt.Errorf("domain extensions file path is required")
	}

	if cfg.Domain.Timeout <= 0 {
		return fmt.Errorf("domain timeout must be positive")
	}

	if cfg.Domain.MaxConcurrentChecks <= 0 {
		return fmt.Errorf("max concurrent checks must be positive")
	}

	return nil
}
