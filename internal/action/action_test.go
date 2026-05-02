package action

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestActionSupportsCommaDelimitedRuleInput(t *testing.T) {
	path := filepath.Join("..", "..", "action.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read action.yml: %v", err)
	}

	content := string(data)
	required := []string{
		"comma-delimited list",
		"HERMESSCAN_RULE_INPUT: ${{ inputs.rule }}",
		`IFS=',' read -r -a rule_values <<< "${rule_input}"`,
		`args+=(--rule "${rule_value}")`,
	}
	for _, value := range required {
		if !strings.Contains(content, value) {
			t.Fatalf("action.yml is missing %q", value)
		}
	}
}
