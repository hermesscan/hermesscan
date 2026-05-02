package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config describes optional repository-level HermesScan settings.
type Config struct {
	Rules               string            `json:"rules,omitempty"`
	Exclude             []string          `json:"exclude,omitempty"`
	Include             []string          `json:"include,omitempty"`
	DisabledRules       []string          `json:"disabledRules,omitempty"`
	SeverityOverrides   map[string]string `json:"severityOverrides,omitempty"`
	FailOn              string            `json:"failOn,omitempty"`
	MinSeverity         string            `json:"minSeverity,omitempty"`
	SuppressionsEnabled *bool             `json:"suppressionsEnabled,omitempty"`
}

// Default returns a conservative starter configuration.
func Default() Config {
	enabled := true
	return Config{
		Rules: "",
		Exclude: []string{
			"dist/**",
			"build/**",
			"node_modules/**",
			"vendor/**",
			"reports/**",
			"coverage/**",
			"tmp/**",
			".git/**",
		},
		Include:             []string{},
		DisabledRules:       []string{},
		SeverityOverrides:   map[string]string{},
		FailOn:              "high",
		MinSeverity:         "",
		SuppressionsEnabled: &enabled,
	}
}

// Load reads a JSON configuration file.
func Load(path string) (Config, error) {
	if path == "" {
		return Config{}, fmt.Errorf("config path is required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}
	return cfg, nil
}

// FindDefault returns .hermesscan.json under root when present.
func FindDefault(root string) string {
	candidate := filepath.Join(root, ".hermesscan.json")
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return ""
}

// WriteDefault writes a starter configuration to path.
func WriteDefault(path string, overwrite bool) error {
	if path == "" {
		return fmt.Errorf("config path is required")
	}
	if !overwrite {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("config %q already exists", path)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	data, err := json.MarshalIndent(Default(), "", "  ")
	if err != nil {
		return fmt.Errorf("serialize default config: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config %q: %w", path, err)
	}
	return nil
}

// SuppressionsEnabledValue returns true unless explicitly disabled.
func (c Config) SuppressionsEnabledValue() bool {
	if c.SuppressionsEnabled == nil {
		return true
	}
	return *c.SuppressionsEnabled
}
