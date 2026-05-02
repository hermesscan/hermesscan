package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hermesscan/hermesscan/internal/scanner"
)

func TestWriteSARIF(t *testing.T) {
	result := scanner.Result{
		Findings: []scanner.Finding{
			{
				RuleID:         "HMS0001",
				RuleName:       "Sleep-based synchronization",
				Severity:       "Medium",
				File:           "examples\\ci.ps1",
				Line:           3,
				Column:         1,
				Description:    "Sleep was used.",
				Recommendation: "Use readiness checks.",
			},
		},
	}

	var buffer bytes.Buffer
	if err := WriteSARIF(&buffer, result); err != nil {
		t.Fatalf("WriteSARIF returned error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(buffer.Bytes(), &decoded); err != nil {
		t.Fatalf("SARIF output is not valid JSON: %v", err)
	}

	if decoded["version"] != "2.1.0" {
		t.Fatalf("unexpected SARIF version: %#v", decoded["version"])
	}
}
