package rules

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//go:embed defaults/hermes.rules.json
var embeddedRules embed.FS

var validSeverities = map[string]bool{
	"info": true, "low": true, "medium": true, "high": true, "critical": true,
}

// Load reads a JSON rule catalog from disk. If path is empty, it loads the
// embedded default HermesScan rule catalog.
func Load(path string) ([]Rule, error) {
	if path == "" {
		return LoadDefault()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read rules %q: %w", path, err)
	}
	return parseRuleCatalog(data, path)
}

// LoadDefault returns the embedded default HermesScan rule catalog.
func LoadDefault() ([]Rule, error) {
	data, err := embeddedRules.ReadFile("defaults/hermes.rules.json")
	if err != nil {
		return nil, fmt.Errorf("read embedded rules: %w", err)
	}
	return parseRuleCatalog(data, "embedded defaults")
}

func parseRuleCatalog(data []byte, source string) ([]Rule, error) {
	var loaded []Rule
	if err := json.Unmarshal(data, &loaded); err != nil {
		return nil, fmt.Errorf("parse rules %q: %w", source, err)
	}

	seen := make(map[string]bool)
	for i, rule := range loaded {
		if rule.ID == "" {
			return nil, fmt.Errorf("rule at index %d is missing id", i)
		}
		upperID := strings.ToUpper(rule.ID)
		if seen[upperID] {
			return nil, fmt.Errorf("duplicate rule id %q", rule.ID)
		}
		seen[upperID] = true
		if rule.Pattern == "" {
			return nil, fmt.Errorf("rule %q is missing pattern", rule.ID)
		}
		if rule.Severity == "" {
			return nil, fmt.Errorf("rule %q is missing severity", rule.ID)
		}
		if !validSeverities[strings.ToLower(rule.Severity)] {
			return nil, fmt.Errorf("rule %q has unsupported severity %q", rule.ID, rule.Severity)
		}
	}

	return loaded, nil
}
