package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureOutputDirectoryCreatesParent(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "reports", "nested", "result.json")
	if err := ensureOutputDirectory(output); err != nil {
		t.Fatalf("ensureOutputDirectory failed: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(output)); err != nil {
		t.Fatalf("expected parent directory: %v", err)
	}
}

func TestParseScanOptionsTracksRulesFlag(t *testing.T) {
	options, err := parseScanOptions([]string{".", "--rules", "custom.json", "--config", ".hermesscan.json"})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if !options.rulesProvided {
		t.Fatalf("expected rulesProvided")
	}
	if !options.configProvided {
		t.Fatalf("expected configProvided")
	}
}

func TestParseScanOptionsNoFailIncludeExclude(t *testing.T) {
	options, err := parseScanOptions([]string{".", "--no-fail", "--include", "scripts/**", "--exclude", "dist/**"})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if !options.noFail {
		t.Fatalf("expected noFail")
	}
	if len(options.include) != 1 || options.include[0] != "scripts/**" {
		t.Fatalf("unexpected include values: %#v", options.include)
	}
	if len(options.exclude) != 1 || options.exclude[0] != "dist/**" {
		t.Fatalf("unexpected exclude values: %#v", options.exclude)
	}
}

func TestParseScanOptionsPhaseFiveFlags(t *testing.T) {
	options, err := parseScanOptions([]string{".", "--category", "cache", "--tag", "docker", "--changed-only", "--changed-base", "origin/main", "--github-annotations"})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(options.category) != 1 || options.category[0] != "cache" {
		t.Fatalf("unexpected category values: %#v", options.category)
	}
	if len(options.tag) != 1 || options.tag[0] != "docker" {
		t.Fatalf("unexpected tag values: %#v", options.tag)
	}
	if !options.changedOnly || options.changedBase != "origin/main" {
		t.Fatalf("unexpected changed options: %#v", options)
	}
	if !options.githubAnnotations || options.format != "github" {
		t.Fatalf("expected github annotation format: %#v", options)
	}
}
