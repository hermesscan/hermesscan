package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hermesscan/hermesscan/internal/rules"
)

func TestDefaultRulePrecisionPostgresPortRequiresExposureContext(t *testing.T) {
	noRisk := scanDefaultRuleFixture(t, "workflow.yml", "env:\n  POSTGRES_DOC: \"PostgreSQL normally listens on 5432\"\n", "HMS0002")
	if len(noRisk.Findings) != 0 {
		t.Fatalf("expected no bare-port findings, got %#v", noRisk.Findings)
	}

	risk := scanDefaultRuleFixture(t, "workflow.yml", "services:\n  postgres:\n    ports:\n      - \"5432:5432\"\n", "HMS0002")
	if len(risk.Findings) != 1 {
		t.Fatalf("expected one exposed-port finding, got %#v", risk.Findings)
	}
}

func TestDefaultRulePrecisionMutableActionAllowsVersionTags(t *testing.T) {
	noRisk := scanDefaultRuleFixture(t, "workflow.yml", "steps:\n  - uses: actions/checkout@v4\n", "HMS0009")
	if len(noRisk.Findings) != 0 {
		t.Fatalf("expected no version-tag findings, got %#v", noRisk.Findings)
	}

	risk := scanDefaultRuleFixture(t, "workflow.yml", "steps:\n  - uses: example/action@main\n", "HMS0009")
	if len(risk.Findings) != 1 {
		t.Fatalf("expected one mutable-action finding, got %#v", risk.Findings)
	}
}

func TestDefaultRulePrecisionBroadCacheAllowsLockfileHash(t *testing.T) {
	noRisk := scanDefaultRuleFixture(t, "workflow.yml", "key: ${{ runner.os }}-${{ hashFiles('**/go.sum') }}\n", "HMS0016")
	if len(noRisk.Findings) != 0 {
		t.Fatalf("expected no specific-cache-key findings, got %#v", noRisk.Findings)
	}

	risk := scanDefaultRuleFixture(t, "workflow.yml", "key: ${{ runner.os }}\n", "HMS0016")
	if len(risk.Findings) != 1 {
		t.Fatalf("expected one broad-cache-key finding, got %#v", risk.Findings)
	}
}

func scanDefaultRuleFixture(t *testing.T, filename string, content string, ruleID string) Result {
	t.Helper()

	root := t.TempDir()
	path := filepath.Join(root, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	loaded, err := rules.LoadDefault()
	if err != nil {
		t.Fatalf("LoadDefault returned error: %v", err)
	}

	options := NewOptionsFromConfigValues(nil, nil, []string{ruleID}, nil, nil, true)
	result, err := ScanWithOptions(root, loaded, options)
	if err != nil {
		t.Fatalf("ScanWithOptions returned error: %v", err)
	}
	return result
}
