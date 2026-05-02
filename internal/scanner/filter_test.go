package scanner

import "testing"

func TestFilterByMinSeverity(t *testing.T) {
	result := Result{
		Findings: []Finding{
			{RuleID: "LOW", Severity: "Low"},
			{RuleID: "MED", Severity: "Medium"},
			{RuleID: "HIGH", Severity: "High"},
		},
	}

	filtered := FilterByMinSeverity(result, "medium")
	if len(filtered.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(filtered.Findings))
	}
	if filtered.Findings[0].RuleID != "MED" || filtered.Findings[1].RuleID != "HIGH" {
		t.Fatalf("unexpected findings after filter: %#v", filtered.Findings)
	}
}

func TestSeverityCounts(t *testing.T) {
	result := Result{
		Findings: []Finding{
			{Severity: "Critical"},
			{Severity: "High"},
			{Severity: "High"},
			{Severity: "Medium"},
			{Severity: "Low"},
			{Severity: "Info"},
		},
	}

	counts := SeverityCounts(result)
	if counts["Critical"] != 1 || counts["High"] != 2 || counts["Medium"] != 1 || counts["Low"] != 1 || counts["Info"] != 1 {
		t.Fatalf("unexpected counts: %#v", counts)
	}
}
