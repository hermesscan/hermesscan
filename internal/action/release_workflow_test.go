package action

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseWorkflowPublishesSBOMWithChecksums(t *testing.T) {
	path := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read release workflow: %v", err)
	}

	content := string(data)
	required := []string{
		"uses: anchore/sbom-action@v0",
		"format: spdx-json",
		"output-file: dist/hermesscan.spdx.json",
		"upload-artifact: false",
		"upload-release-assets: false",
		"sha256sum hermesscan-* hermesscan.spdx.json hermesscan.sarif hermesscan.md > checksums.txt",
		"files: dist/*",
		"actions: read",
	}
	for _, value := range required {
		if !strings.Contains(content, value) {
			t.Fatalf("release.yml is missing %q", value)
		}
	}
}

func TestReleaseWorkflowPublishesHermesScanReports(t *testing.T) {
	path := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read release workflow: %v", err)
	}

	content := string(data)
	required := []string{
		"security-events: write",
		"Generate HermesScan SARIF report",
		"--format sarif",
		"--output dist/hermesscan.sarif",
		"github/codeql-action/upload-sarif@v4",
		"sarif_file: dist/hermesscan.sarif",
		"Generate HermesScan Markdown report",
		"--format markdown",
		"--output dist/hermesscan.md",
		"name: hermesscan-release-evidence",
		"Gate release with HermesScan",
		"--fail-on high",
	}
	for _, value := range required {
		if !strings.Contains(content, value) {
			t.Fatalf("release.yml is missing %q", value)
		}
	}
}

func TestReleaseSmokeWorkflowVerifiesReleaseEvidenceAssets(t *testing.T) {
	path := filepath.Join("..", "..", ".github", "workflows", "release-smoke.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read release smoke workflow: %v", err)
	}

	content := string(data)
	expectedAssets := []string{
		"'hermesscan.spdx.json'",
		"'hermesscan.sarif'",
		"'hermesscan.md'",
	}
	for _, asset := range expectedAssets {
		if occurrences := strings.Count(content, asset); occurrences < 2 {
			t.Fatalf("release-smoke.yml should download and checksum %s, found %d references", asset, occurrences)
		}
	}
}
