package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main WAF configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Security SecurityConfig `yaml:"security"`
	Logging  LoggingConfig  `yaml:"logging"`
	Rules    RulesConfig    `yaml:"rules"`
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	ListenAddress      string `yaml:"listen_address"`
	UpstreamURL        string `yaml:"upstream_url"`
	ReadTimeoutSeconds int    `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds int   `yaml:"write_timeout_seconds"`
	IdleTimeoutSeconds int    `yaml:"idle_timeout_seconds"`
}

// SecurityConfig contains security-related settings
type SecurityConfig struct {
	AnomalyThreshold int  `yaml:"anomaly_threshold"`
	LogRequestBody   bool `yaml:"log_request_bodies"`
	RateLimit        RateLimitConfig `yaml:"rate_limit"`
	IPFilter         IPFilterConfig  `yaml:"ip_filter"`
}

// RateLimitConfig contains rate limiting settings
type RateLimitConfig struct {
	Enabled      bool `yaml:"enabled"`
	MaxRequests  int  `yaml:"max_requests"`
	WindowSeconds int `yaml:"window_seconds"`
}

// IPFilterConfig contains IP filtering settings
type IPFilterConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Whitelist []string `yaml:"whitelist"`
	Blacklist []string `yaml:"blacklist"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

// RulesConfig contains rule file paths
type RulesConfig struct {
	Files []string `yaml:"files"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults if not specified
	if cfg.Server.ListenAddress == "" {
		cfg.Server.ListenAddress = ":8080"
	}
	if cfg.Server.ReadTimeoutSeconds == 0 {
		cfg.Server.ReadTimeoutSeconds = 10
	}
	if cfg.Server.WriteTimeoutSeconds == 0 {
		cfg.Server.WriteTimeoutSeconds = 10
	}
	if cfg.Server.IdleTimeoutSeconds == 0 {
		cfg.Server.IdleTimeoutSeconds = 60
	}
	if cfg.Security.AnomalyThreshold == 0 {
		cfg.Security.AnomalyThreshold = 10
	}
	// Rate limit defaults
	if cfg.Security.RateLimit.MaxRequests == 0 {
		cfg.Security.RateLimit.MaxRequests = 100
	}
	if cfg.Security.RateLimit.WindowSeconds == 0 {
		cfg.Security.RateLimit.WindowSeconds = 60
	}
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Logging.Output == "" {
		cfg.Logging.Output = "stdout"
	}

	return &cfg, nil
}

// ReadTimeout returns the read timeout as a time.Duration
func (s *ServerConfig) ReadTimeout() time.Duration {
	return time.Duration(s.ReadTimeoutSeconds) * time.Second
}

// WriteTimeout returns the write timeout as a time.Duration
func (s *ServerConfig) WriteTimeout() time.Duration {
	return time.Duration(s.WriteTimeoutSeconds) * time.Second
}

// IdleTimeout returns the idle timeout as a time.Duration
func (s *ServerConfig) IdleTimeout() time.Duration {
	return time.Duration(s.IdleTimeoutSeconds) * time.Second
}

