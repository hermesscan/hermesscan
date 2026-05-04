package packaging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

type scoopManifest struct {
	Version      string                    `json:"version"`
	Description  string                    `json:"description"`
	Homepage     string                    `json:"homepage"`
	License      string                    `json:"license"`
	Architecture map[string]scoopArchEntry `json:"architecture"`
	Bin          string                    `json:"bin"`
	Checkver     struct {
		GitHub string `json:"github"`
	} `json:"checkver"`
	Autoupdate struct {
		Architecture map[string]scoopArchEntry `json:"architecture"`
	} `json:"autoupdate"`
}

type scoopArchEntry struct {
	URL  string `json:"url"`
	Hash string `json:"hash,omitempty"`
}

func TestScoopManifestContract(t *testing.T) {
	path := filepath.Join("..", "..", "packaging", "scoop", "hermesscan.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Scoop manifest: %v", err)
	}

	var manifest scoopManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse Scoop manifest: %v", err)
	}

	if !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(manifest.Version) {
		t.Fatalf("version should be release semver without v prefix, got %q", manifest.Version)
	}
	if manifest.Description != "Static analyzer for build scripts, CI scripts, and pipeline definitions." {
		t.Fatalf("unexpected description: %q", manifest.Description)
	}
	if manifest.Homepage != "https://github.com/hermesscan/hermesscan" {
		t.Fatalf("unexpected homepage: %q", manifest.Homepage)
	}
	if manifest.License != "MIT" {
		t.Fatalf("unexpected license: %q", manifest.License)
	}
	if manifest.Bin != "hermesscan.exe" {
		t.Fatalf("unexpected bin: %q", manifest.Bin)
	}
	if manifest.Checkver.GitHub != "https://github.com/hermesscan/hermesscan" {
		t.Fatalf("unexpected checkver github URL: %q", manifest.Checkver.GitHub)
	}

	expectedAssets := map[string]string{
		"64bit": "hermesscan-windows-amd64.exe",
		"arm64": "hermesscan-windows-arm64.exe",
	}
	for arch, asset := range expectedAssets {
		entry, ok := manifest.Architecture[arch]
		if !ok {
			t.Fatalf("missing architecture entry %q", arch)
		}
		expectedURL := fmt.Sprintf("https://github.com/hermesscan/hermesscan/releases/download/v%s/%s#/hermesscan.exe", manifest.Version, asset)
		if entry.URL != expectedURL {
			t.Fatalf("unexpected %s URL: got %q, want %q", arch, entry.URL, expectedURL)
		}
		if !regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString(entry.Hash) {
			t.Fatalf("unexpected %s hash: %q", arch, entry.Hash)
		}

		updateEntry, ok := manifest.Autoupdate.Architecture[arch]
		if !ok {
			t.Fatalf("missing autoupdate architecture entry %q", arch)
		}
		expectedUpdateURL := fmt.Sprintf("https://github.com/hermesscan/hermesscan/releases/download/v$version/%s#/hermesscan.exe", asset)
		if updateEntry.URL != expectedUpdateURL {
			t.Fatalf("unexpected %s autoupdate URL: got %q, want %q", arch, updateEntry.URL, expectedUpdateURL)
		}
	}
}
