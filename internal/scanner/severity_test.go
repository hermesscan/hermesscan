package scanner

import "testing"

func TestMeetsThreshold(t *testing.T) {
	if !MeetsThreshold("High", "medium") {
		t.Fatal("expected high to meet medium threshold")
	}
	if MeetsThreshold("Low", "high") {
		t.Fatal("did not expect low to meet high threshold")
	}
	if MeetsThreshold("Critical", "none") {
		t.Fatal("did not expect any severity to meet none threshold")
	}
}
