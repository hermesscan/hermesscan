package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ProfileMinimal     = "minimal"
	ProfileCI          = "ci"
	ProfileSupplyChain = "supply-chain"
)

// Config describes optional repository-level HermesScan settings.
type Config struct {
	Rules               string            `json:"rules,omitempty"`
	Exclude             []string          `json:"exclude,omitempty"`
	Include             []string          `json:"include,omitempty"`
	EnabledRules        []string          `json:"enabledRules,omitempty"`
	DisabledRules       []string          `json:"disabledRules,omitempty"`
	Categories          []string          `json:"categories,omitempty"`
	Tags                []string          `json:"tags,omitempty"`
	SeverityOverrides   map[string]string `json:"severityOverrides,omitempty"`
	FailOn              string            `json:"failOn,omitempty"`
	MinSeverity         string            `json:"minSeverity,omitempty"`
	SuppressionsEnabled *bool             `json:"suppressionsEnabled,omitempty"`
}

// Default returns a conservative starter configuration.
func Default() Config {
	return Profile(ProfileCI)
}

// Profile returns a starter configuration for a named adoption profile.
func Profile(name string) Config {
	enabled := true
	cfg := Config{
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
		EnabledRules:        []string{},
		DisabledRules:       []string{},
		Categories:          []string{},
		Tags:                []string{},
		SeverityOverrides:   map[string]string{},
		FailOn:              "high",
		MinSeverity:         "",
		SuppressionsEnabled: &enabled,
	}

	switch normalizeProfileName(name) {
	case ProfileMinimal:
		cfg.FailOn = ""
	case ProfileSupplyChain:
		cfg.Categories = []string{"supply-chain"}
	}

	return cfg
}

// SupportedProfiles returns the accepted init profile names.
func SupportedProfiles() []string {
	return []string{ProfileMinimal, ProfileCI, ProfileSupplyChain}
}

// NormalizeProfile returns the canonical profile name or an error.
func NormalizeProfile(name string) (string, error) {
	normalized := normalizeProfileName(name)
	for _, profile := range SupportedProfiles() {
		if normalized == profile {
			return normalized, nil
		}
	}
	return "", fmt.Errorf("unsupported profile %q; expected one of: %s", name, strings.Join(SupportedProfiles(), ", "))
}

func normalizeProfileName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return ProfileCI
	}
	return name
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
	return WriteProfile(path, overwrite, ProfileCI)
}

// WriteProfile writes a starter configuration profile to path.
func WriteProfile(path string, overwrite bool, profile string) error {
	if path == "" {
		return fmt.Errorf("config path is required")
	}
	profile, err := NormalizeProfile(profile)
	if err != nil {
		return err
	}
	if !overwrite {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("config %q already exists", path)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	data, err := json.MarshalIndent(Profile(profile), "", "  ")
	if err != nil {
		return fmt.Errorf("serialize %s config: %w", profile, err)
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
