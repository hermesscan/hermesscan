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
		"sha256sum hermesscan-* hermesscan.spdx.json > checksums.txt",
		"files: dist/*",
		"actions: read",
	}
	for _, value := range required {
		if !strings.Contains(content, value) {
			t.Fatalf("release.yml is missing %q", value)
		}
	}
}

func TestReleaseSmokeWorkflowVerifiesSBOMAsset(t *testing.T) {
	path := filepath.Join("..", "..", ".github", "workflows", "release-smoke.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read release smoke workflow: %v", err)
	}

	content := string(data)
	if occurrences := strings.Count(content, "'hermesscan.spdx.json'"); occurrences < 2 {
		t.Fatalf("release-smoke.yml should download and checksum the SBOM asset, found %d references", occurrences)
	}
}
