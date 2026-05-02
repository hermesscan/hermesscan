package scanner

// FilterByMinSeverity returns a copy of result containing only findings whose severity
// is greater than or equal to minSeverity. Empty, "none", and "info" retain all findings.
func FilterByMinSeverity(result Result, minSeverity string) Result {
	if minSeverity == "" || minSeverity == "none" || SeverityRank(minSeverity) <= SeverityRank("info") {
		return result
	}

	filtered := result
	filtered.Findings = make([]Finding, 0, len(result.Findings))
	for _, finding := range result.Findings {
		if MeetsThreshold(finding.Severity, minSeverity) {
			filtered.Findings = append(filtered.Findings, finding)
		}
	}
	return filtered
}

// SeverityCounts returns finding counts by canonical severity name.
func SeverityCounts(result Result) map[string]int {
	counts := map[string]int{
		"Critical": 0,
		"High":     0,
		"Medium":   0,
		"Low":      0,
		"Info":     0,
	}

	for _, finding := range result.Findings {
		switch SeverityRank(finding.Severity) {
		case 5:
			counts["Critical"]++
		case 4:
			counts["High"]++
		case 3:
			counts["Medium"]++
		case 2:
			counts["Low"]++
		default:
			counts["Info"]++
		}
	}

	return counts
}
