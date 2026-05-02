package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hermesscan/hermesscan/internal/scanner"
)

func TestWriteGitHubAnnotations(t *testing.T) {
	result := scanner.Result{Findings: []scanner.Finding{{RuleID: "HMS0001", RuleName: "Sleep", Severity: "High", File: "examples/test.sh", Line: 3, Column: 2, Description: "bad", Recommendation: "fix"}}}
	var buf bytes.Buffer
	if err := WriteGitHubAnnotations(&buf, result); err != nil {
		t.Fatalf("WriteGitHubAnnotations returned error: %v", err)
	}
	text := buf.String()
	if !strings.Contains(text, "::error file=examples/test.sh,line=3,col=2") {
		t.Fatalf("annotation was not written correctly: %s", text)
	}
}

func TestWriteGitHubAnnotationsNormalizesWindowsPaths(t *testing.T) {
	result := scanner.Result{Findings: []scanner.Finding{{RuleID: "HMS0002", RuleName: "Port", Severity: "Medium", File: `examples\\ci-with-risks.ps1`, Line: 7, Column: 14, Description: "bad", Recommendation: "fix"}}}
	var buf bytes.Buffer
	if err := WriteGitHubAnnotations(&buf, result); err != nil {
		t.Fatalf("WriteGitHubAnnotations returned error: %v", err)
	}
	text := buf.String()
	if !strings.Contains(text, "file=examples/ci-with-risks.ps1") {
		t.Fatalf("annotation path was not normalized: %s", text)
	}
}
