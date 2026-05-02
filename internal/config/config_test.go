package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteDefaultAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".hermesscan.json")

	if err := WriteDefault(path, false); err != nil {
		t.Fatalf("WriteDefault failed: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Rules != "" {
		t.Fatalf("expected embedded rules by default")
	}
	if len(cfg.Exclude) == 0 {
		t.Fatalf("expected default excludes")
	}
	if !cfg.SuppressionsEnabledValue() {
		t.Fatalf("expected suppressions enabled by default")
	}
}

func TestFindDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".hermesscan.json")
	if err := os.WriteFile(path, []byte(`{}`), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if got := FindDefault(dir); got != path {
		t.Fatalf("expected %q, got %q", path, got)
	}
}
