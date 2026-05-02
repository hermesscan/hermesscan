package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hermesscan/hermesscan/internal/rules"
)

func TestScanWithCategoryFilter(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ci.sh")
	if err := os.WriteFile(path, []byte("sleep 30\nnpm ci\n"), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	loadedRules := []rules.Rule{
		{ID: "HMS0001", Name: "Sleep", Severity: "Medium", Category: "reliability", FileTypes: []string{"bash"}, Pattern: `sleep\s+[0-9]+`},
		{ID: "HMS0010", Name: "NPM", Severity: "Medium", Category: "cache", FileTypes: []string{"bash"}, Pattern: `npm\s+ci`},
	}

	options := DefaultOptions()
	options.Categories = []string{"cache"}
	result, err := ScanWithOptions(root, loadedRules, options)
	if err != nil {
		t.Fatalf("scan returned error: %v", err)
	}
	if len(result.Findings) != 1 || result.Findings[0].RuleID != "HMS0010" {
		t.Fatalf("expected only HMS0010, got %#v", result.Findings)
	}
}

func TestScanWithTagFilter(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ci.sh")
	if err := os.WriteFile(path, []byte("sleep 30\ndocker build .\n"), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	loadedRules := []rules.Rule{
		{ID: "HMS0001", Name: "Sleep", Severity: "Medium", Tags: []string{"timing"}, FileTypes: []string{"bash"}, Pattern: `sleep\s+[0-9]+`},
		{ID: "HMS0011", Name: "Docker", Severity: "Medium", Tags: []string{"docker"}, FileTypes: []string{"bash"}, Pattern: `docker\s+build`},
	}

	options := DefaultOptions()
	options.Tags = []string{"docker"}
	result, err := ScanWithOptions(root, loadedRules, options)
	if err != nil {
		t.Fatalf("scan returned error: %v", err)
	}
	if len(result.Findings) != 1 || result.Findings[0].RuleID != "HMS0011" {
		t.Fatalf("expected only HMS0011, got %#v", result.Findings)
	}
}
