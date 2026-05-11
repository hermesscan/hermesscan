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

func TestWriteProfileMinimalIsAdvisory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".hermesscan.json")

	if err := WriteProfile(path, false, ProfileMinimal); err != nil {
		t.Fatalf("WriteProfile failed: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.FailOn != "" {
		t.Fatalf("expected no fail threshold for minimal profile, got %q", cfg.FailOn)
	}
	if len(cfg.Categories) != 0 {
		t.Fatalf("expected no category filter for minimal profile, got %#v", cfg.Categories)
	}
}

func TestWriteProfileSupplyChainFiltersCategory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".hermesscan.json")

	if err := WriteProfile(path, false, ProfileSupplyChain); err != nil {
		t.Fatalf("WriteProfile failed: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.FailOn != "high" {
		t.Fatalf("expected high fail threshold, got %q", cfg.FailOn)
	}
	if len(cfg.Categories) != 1 || cfg.Categories[0] != "supply-chain" {
		t.Fatalf("unexpected category filter: %#v", cfg.Categories)
	}
}

func TestNormalizeProfileRejectsUnknownValue(t *testing.T) {
	if _, err := NormalizeProfile("strict"); err == nil {
		t.Fatalf("expected unknown profile to fail")
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
