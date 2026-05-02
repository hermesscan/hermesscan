package scanner

import "strings"

var severityRank = map[string]int{
	"info":     1,
	"low":      2,
	"medium":   3,
	"high":     4,
	"critical": 5,
}

// SeverityRank returns a sortable severity value.
func SeverityRank(severity string) int {
	value, ok := severityRank[strings.ToLower(severity)]
	if !ok {
		return 0
	}
	return value
}

// MeetsThreshold returns true when severity is greater than or equal to threshold.
func MeetsThreshold(severity string, threshold string) bool {
	if threshold == "" || strings.EqualFold(threshold, "none") {
		return false
	}
	return SeverityRank(severity) >= SeverityRank(threshold)
}
