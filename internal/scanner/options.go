package scanner

import "strings"

// Options controls scanning behavior after rule loading.
type Options struct {
	Exclude             []string
	Include             []string
	DisabledRules       map[string]bool
	SeverityOverrides   map[string]string
	SuppressionsEnabled bool
	Categories          []string
	Tags                []string
	ChangedOnly         bool
	ChangedBase         string
}

// DefaultOptions returns scan defaults.
func DefaultOptions() Options {
	return Options{SuppressionsEnabled: true}
}

func normalizeRuleSet(values []string) map[string]bool {
	result := make(map[string]bool)
	for _, value := range values {
		value = strings.ToUpper(strings.TrimSpace(value))
		if value != "" {
			result[value] = true
		}
	}
	return result
}

func normalizeStringSet(values []string) map[string]bool {
	result := make(map[string]bool)
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value != "" {
			result[value] = true
		}
	}
	return result
}
