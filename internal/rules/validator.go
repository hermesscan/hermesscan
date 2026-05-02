package rules

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateCatalog checks that a rule catalog is complete enough to publish.
func ValidateCatalog(loaded []Rule) error {
	if len(loaded) == 0 {
		return fmt.Errorf("rule catalog is empty")
	}

	seenIDs := make(map[string]bool)
	for index, rule := range loaded {
		label := fmt.Sprintf("rule at index %d", index)
		if strings.TrimSpace(rule.ID) != "" {
			label = fmt.Sprintf("rule %q", rule.ID)
		}

		if strings.TrimSpace(rule.ID) == "" {
			return fmt.Errorf("%s is missing id", label)
		}
		if strings.TrimSpace(rule.ID) != rule.ID {
			return fmt.Errorf("%s id must not have leading or trailing whitespace", label)
		}
		upperID := strings.ToUpper(rule.ID)
		if seenIDs[upperID] {
			return fmt.Errorf("duplicate rule id %q", rule.ID)
		}
		seenIDs[upperID] = true

		if strings.TrimSpace(rule.Name) == "" {
			return fmt.Errorf("%s is missing name", label)
		}
		if strings.TrimSpace(rule.Severity) == "" {
			return fmt.Errorf("%s is missing severity", label)
		}
		if !validSeverities[strings.ToLower(rule.Severity)] {
			return fmt.Errorf("%s has unsupported severity %q", label, rule.Severity)
		}
		if strings.TrimSpace(rule.Category) == "" {
			return fmt.Errorf("%s is missing category", label)
		}
		if len(rule.Tags) == 0 {
			return fmt.Errorf("%s is missing tags", label)
		}
		if len(rule.FileTypes) == 0 {
			return fmt.Errorf("%s is missing fileTypes", label)
		}
		if strings.TrimSpace(rule.Pattern) == "" {
			return fmt.Errorf("%s is missing pattern", label)
		}
		if _, err := regexp.Compile(rule.Pattern); err != nil {
			return fmt.Errorf("%s has invalid pattern: %w", label, err)
		}
		if strings.TrimSpace(rule.Description) == "" {
			return fmt.Errorf("%s is missing description", label)
		}
		if strings.TrimSpace(rule.Recommendation) == "" {
			return fmt.Errorf("%s is missing recommendation", label)
		}

		seenTags := make(map[string]bool)
		for _, tag := range rule.Tags {
			cleanTag := strings.TrimSpace(tag)
			if cleanTag == "" {
				return fmt.Errorf("%s has an empty tag", label)
			}
			key := strings.ToLower(cleanTag)
			if seenTags[key] {
				return fmt.Errorf("%s has duplicate tag %q", label, tag)
			}
			seenTags[key] = true
		}
		for _, fileType := range rule.FileTypes {
			if strings.TrimSpace(fileType) == "" {
				return fmt.Errorf("%s has an empty file type", label)
			}
		}
	}

	return nil
}
