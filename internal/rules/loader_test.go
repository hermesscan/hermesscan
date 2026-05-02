package rules

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadRejectsDuplicateRuleIDs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.json")
	data := `[
{"id":"HMS0001","name":"one","severity":"Low","fileTypes":["bash"],"pattern":"x","description":"d","recommendation":"r"},
{"id":"hms0001","name":"two","severity":"Low","fileTypes":["bash"],"pattern":"y","description":"d","recommendation":"r"}
]`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatalf("write rules: %v", err)
	}
	if _, err := Load(path); err == nil {
		t.Fatalf("expected duplicate id error")
	}
}

func TestLoadRejectsUnsupportedSeverity(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.json")
	data := `[{"id":"HMS0001","name":"one","severity":"Bad","fileTypes":["bash"],"pattern":"x","description":"d","recommendation":"r"}]`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatalf("write rules: %v", err)
	}
	if _, err := Load(path); err == nil {
		t.Fatalf("expected severity error")
	}
}

func TestValidateCatalogAcceptsDefaultRules(t *testing.T) {
	loaded, err := LoadDefault()
	if err != nil {
		t.Fatalf("LoadDefault returned error: %v", err)
	}
	if err := ValidateCatalog(loaded); err != nil {
		t.Fatalf("ValidateCatalog returned error: %v", err)
	}
}

func TestValidateCatalogRejectsIncompleteRule(t *testing.T) {
	loaded := []Rule{
		{
			ID:             "HMS0001",
			Name:           "Incomplete",
			Severity:       "Low",
			FileTypes:      []string{"bash"},
			Tags:           []string{"test"},
			Pattern:        "sleep",
			Description:    "description",
			Recommendation: "recommendation",
		},
	}
	err := ValidateCatalog(loaded)
	if err == nil || !strings.Contains(err.Error(), "missing category") {
		t.Fatalf("expected missing category error, got %v", err)
	}
}

func TestValidateCatalogRejectsInvalidPattern(t *testing.T) {
	loaded := []Rule{
		{
			ID:             "HMS0001",
			Name:           "Invalid regex",
			Severity:       "Low",
			Category:       "test",
			FileTypes:      []string{"bash"},
			Tags:           []string{"test"},
			Pattern:        "(",
			Description:    "description",
			Recommendation: "recommendation",
		},
	}
	err := ValidateCatalog(loaded)
	if err == nil || !strings.Contains(err.Error(), "invalid pattern") {
		t.Fatalf("expected invalid pattern error, got %v", err)
	}
}

func TestDefaultCatalogMatchesRepositoryCatalog(t *testing.T) {
	embedded, err := LoadDefault()
	if err != nil {
		t.Fatalf("LoadDefault returned error: %v", err)
	}

	repositoryPath := filepath.Join("..", "..", "rules", "hermes.rules.json")
	repository, err := Load(repositoryPath)
	if err != nil {
		t.Fatalf("Load(%q) returned error: %v", repositoryPath, err)
	}

	if !reflect.DeepEqual(repository, embedded) {
		t.Fatalf("repository rule catalog and embedded default catalog differ")
	}
}
