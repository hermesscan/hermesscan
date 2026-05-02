package scanner

import "testing"

func TestParseSuppressionDisableNextLine(t *testing.T) {
	directives := parseSuppression("# hermesscan:disable-next-line HMS0001,HMS0002 reason")
	if len(directives) != 1 {
		t.Fatalf("expected one directive")
	}
	if directives[0].Kind != "disable-next-line" {
		t.Fatalf("unexpected kind %q", directives[0].Kind)
	}
	if !suppressionMatches(directives[0].RuleIDs, "HMS0002") {
		t.Fatalf("expected HMS0002 suppression")
	}
}

func TestParseSuppressionDefaultsToAll(t *testing.T) {
	directives := parseSuppression("# hermesscan:disable-line")
	if len(directives) != 1 {
		t.Fatalf("expected one directive")
	}
	if !suppressionMatches(directives[0].RuleIDs, "HMS9999") {
		t.Fatalf("expected all suppression")
	}
}
