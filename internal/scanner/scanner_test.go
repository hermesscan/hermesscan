package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hermesscan/hermesscan/internal/rules"
)

func TestScanFindsSleep(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ci.sh")
	if err := os.WriteFile(path, []byte("#!/usr/bin/env bash\nsleep 30\n"), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	loaded := []rules.Rule{
		{
			ID:             "HMS0001",
			Name:           "Sleep-based synchronization",
			Severity:       "Medium",
			FileTypes:      []string{"bash"},
			Pattern:        "(?i)\\bsleep\\s+[0-9]+",
			Description:    "test description",
			Recommendation: "test recommendation",
		},
	}

	result, err := Scan(root, loaded)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings = %d; want 1", len(result.Findings))
	}
	if result.Findings[0].Line != 2 {
		t.Fatalf("line = %d; want 2", result.Findings[0].Line)
	}
}

func TestScanSuppressesNextLine(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ci.sh")
	content := "#!/usr/bin/env bash\n# hermesscan:disable-next-line HMS0001\nsleep 30\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	loaded := []rules.Rule{
		{
			ID:             "HMS0001",
			Name:           "Sleep-based synchronization",
			Severity:       "Medium",
			FileTypes:      []string{"bash"},
			Pattern:        "(?i)\\bsleep\\s+[0-9]+",
			Description:    "test description",
			Recommendation: "test recommendation",
		},
	}

	result, err := ScanWithOptions(root, loaded, DefaultOptions())
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(result.Findings) != 0 {
		t.Fatalf("findings = %d; want 0", len(result.Findings))
	}
	if result.SuppressedCount != 1 {
		t.Fatalf("suppressed = %d; want 1", result.SuppressedCount)
	}
}

func TestScanCanDisableRule(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ci.sh")
	if err := os.WriteFile(path, []byte("sleep 30\n"), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	loaded := []rules.Rule{
		{
			ID:             "HMS0001",
			Name:           "Sleep-based synchronization",
			Severity:       "Medium",
			FileTypes:      []string{"bash"},
			Pattern:        "(?i)\\bsleep\\s+[0-9]+",
			Description:    "test description",
			Recommendation: "test recommendation",
		},
	}

	options := NewOptionsFromConfigValues(nil, nil, nil, []string{"HMS0001"}, nil, true)
	result, err := ScanWithOptions(root, loaded, options)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(result.Findings) != 0 {
		t.Fatalf("findings = %d; want 0", len(result.Findings))
	}
	if result.RulesLoaded != 0 {
		t.Fatalf("rules loaded = %d; want 0", result.RulesLoaded)
	}
}

func TestScanCanOverrideSeverity(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ci.sh")
	if err := os.WriteFile(path, []byte("sleep 30\n"), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	loaded := []rules.Rule{
		{
			ID:             "HMS0001",
			Name:           "Sleep-based synchronization",
			Severity:       "Medium",
			FileTypes:      []string{"bash"},
			Pattern:        "(?i)\\bsleep\\s+[0-9]+",
			Description:    "test description",
			Recommendation: "test recommendation",
		},
	}

	options := NewOptionsFromConfigValues(nil, nil, nil, nil, map[string]string{"HMS0001": "Low"}, true)
	result, err := ScanWithOptions(root, loaded, options)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings = %d; want 1", len(result.Findings))
	}
	if result.Findings[0].Severity != "Low" {
		t.Fatalf("severity = %q; want Low", result.Findings[0].Severity)
	}
}

func TestScanExcludesCandidates(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "generated"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(root, "generated", "ci.sh")
	if err := os.WriteFile(path, []byte("sleep 30\n"), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	loaded := []rules.Rule{
		{
			ID:             "HMS0001",
			Name:           "Sleep-based synchronization",
			Severity:       "Medium",
			FileTypes:      []string{"bash"},
			Pattern:        "(?i)\\bsleep\\s+[0-9]+",
			Description:    "test description",
			Recommendation: "test recommendation",
		},
	}

	options := NewOptionsFromConfigValues([]string{filepath.ToSlash(filepath.Join(root, "generated")) + "/**"}, nil, nil, nil, nil, true)
	result, err := ScanWithOptions(root, loaded, options)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if result.FilesScanned != 0 {
		t.Fatalf("files scanned = %d; want 0", result.FilesScanned)
	}
}

func TestScanCanEnableSpecificRules(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ci.sh")
	if err := os.WriteFile(path, []byte("sleep 30\nnpm ci\n"), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	loaded := []rules.Rule{
		{ID: "HMS0001", Name: "Sleep", Severity: "Medium", FileTypes: []string{"bash"}, Pattern: `sleep\s+[0-9]+`},
		{ID: "HMS0010", Name: "NPM", Severity: "Low", FileTypes: []string{"bash"}, Pattern: `npm\s+ci`},
	}

	options := NewOptionsFromConfigValues(nil, nil, []string{"HMS0010"}, nil, nil, true)
	result, err := ScanWithOptions(root, loaded, options)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(result.Findings) != 1 || result.Findings[0].RuleID != "HMS0010" {
		t.Fatalf("expected only HMS0010, got %#v", result.Findings)
	}
	if result.RulesLoaded != 1 {
		t.Fatalf("rules loaded = %d; want 1", result.RulesLoaded)
	}
}

func TestScanRequiresTriggerFilePatternWhenConfigured(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ci.sh")
	if err := os.WriteFile(path, []byte("sleep 30\n"), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	loaded := []rules.Rule{
		{
			ID:                 "HMS0001",
			Name:               "Sleep",
			Severity:           "Medium",
			FileTypes:          []string{"bash"},
			Pattern:            `sleep\s+[0-9]+`,
			TriggerFilePattern: `deploy`,
			Description:        "test description",
			Recommendation:     "test recommendation",
		},
	}

	result, err := ScanWithOptions(root, loaded, DefaultOptions())
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(result.Findings) != 0 {
		t.Fatalf("findings = %d; want 0", len(result.Findings))
	}

	if err := os.WriteFile(path, []byte("# deploy\nsleep 30\n"), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	result, err = ScanWithOptions(root, loaded, DefaultOptions())
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings = %d; want 1", len(result.Findings))
	}
}
