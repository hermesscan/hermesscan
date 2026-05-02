package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hermesscan/hermesscan/internal/rules"
)

func TestIncludeFiltersCandidates(t *testing.T) {
	root := t.TempDir()
	keep := filepath.Join(root, "keep")
	skip := filepath.Join(root, "skip")
	if err := os.MkdirAll(keep, 0755); err != nil {
		t.Fatalf("mkdir keep: %v", err)
	}
	if err := os.MkdirAll(skip, 0755); err != nil {
		t.Fatalf("mkdir skip: %v", err)
	}
	if err := os.WriteFile(filepath.Join(keep, "ci.sh"), []byte("sleep 30\n"), 0644); err != nil {
		t.Fatalf("write keep: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skip, "ci.sh"), []byte("sleep 30\n"), 0644); err != nil {
		t.Fatalf("write skip: %v", err)
	}

	rule := rules.Rule{ID: "HMS0001", Name: "sleep", Severity: "Medium", FileTypes: []string{"bash"}, Pattern: `sleep\s+\d+`}
	includePattern := filepath.ToSlash(filepath.Join(root, "keep")) + "/**"
	options := NewOptionsFromConfigValues(nil, []string{includePattern}, nil, nil, nil, true)
	result, err := ScanWithOptions(root, []rules.Rule{rule}, options)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if result.FilesScanned != 1 {
		t.Fatalf("expected 1 scanned file, got %d", result.FilesScanned)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result.Findings))
	}
}
