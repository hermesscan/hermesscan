package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRejectsDuplicateRuleIDs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.json")
	data := `[
{"id":"HMS0001","name":"one","severity":"Low","fileTypes":["bash"],"pattern":"x","description":"d","recommendation":"r"},
{"id":"hms0001","name":"two","severity":"Low","fileTypes":["bash"],"pattern":"y","description":"d","recommendation":"r"}
]`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatalf("write rules: %v", err)
	}
	if _, err := Load(path); err == nil {
		t.Fatalf("expected duplicate id error")
	}
}

func TestLoadRejectsUnsupportedSeverity(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.json")
	data := `[{"id":"HMS0001","name":"one","severity":"Bad","fileTypes":["bash"],"pattern":"x","description":"d","recommendation":"r"}]`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatalf("write rules: %v", err)
	}
	if _, err := Load(path); err == nil {
		t.Fatalf("expected severity error")
	}
}
