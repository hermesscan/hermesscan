package baseline

import (
	"path/filepath"
	"testing"

	"github.com/hermesscan/hermesscan/internal/scanner"
)

func TestBaselineApplyRemovesKnownFindings(t *testing.T) {
	finding := scanner.Finding{RuleID: "HMS0001", Severity: "Medium", File: "scripts/ci.sh", Line: 3, Match: "sleep 30"}
	finding.Fingerprint = scanner.Fingerprint(finding)
	result := scanner.Result{Findings: []scanner.Finding{finding}}
	file := FromResult(result)

	filtered := Apply(result, file)
	if len(filtered.Findings) != 0 {
		t.Fatalf("expected finding to be baseline-suppressed")
	}
	if filtered.BaselineSuppressedCount != 1 {
		t.Fatalf("expected baseline suppressed count 1, got %d", filtered.BaselineSuppressedCount)
	}
}

func TestBaselineSaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", ".hermesscan-baseline.json")
	finding := scanner.Finding{RuleID: "HMS0001", Severity: "Medium", File: "scripts/ci.sh", Line: 3, Match: "sleep 30"}
	finding.Fingerprint = scanner.Fingerprint(finding)
	file := FromResult(scanner.Result{Findings: []scanner.Finding{finding}})
	if err := Save(path, file); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if len(loaded.Findings) != 1 || loaded.Findings[0].Fingerprint == "" {
		t.Fatalf("unexpected loaded baseline: %#v", loaded)
	}
}
