package scanner

import (
	"strings"
	"unicode"
)

type suppressionDirective struct {
	Kind    string
	RuleIDs []string
}

func parseSuppression(line string) []suppressionDirective {
	lower := strings.ToLower(line)
	markers := []string{
		"hermesscan:disable-next-line",
		"hermesscan:disable-line",
		"hermesscan:disable-file",
	}

	var result []suppressionDirective
	for _, marker := range markers {
		idx := strings.Index(lower, marker)
		if idx < 0 {
			continue
		}
		remainder := line[idx+len(marker):]
		result = append(result, suppressionDirective{
			Kind:    strings.TrimPrefix(marker, "hermesscan:"),
			RuleIDs: parseSuppressionRuleIDs(remainder),
		})
	}
	return result
}

func parseSuppressionRuleIDs(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{"all"}
	}

	fields := strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || r == ',' || r == ';' || r == ':'
	})

	var ids []string
	for _, field := range fields {
		field = strings.Trim(field, "#/-*[](){}")
		upper := strings.ToUpper(field)
		if upper == "ALL" {
			ids = append(ids, "all")
			continue
		}
		if strings.HasPrefix(upper, "HMS") {
			ids = append(ids, upper)
		}
	}
	if len(ids) == 0 {
		return []string{"all"}
	}
	return ids
}

func suppressionMatches(ids []string, ruleID string) bool {
	for _, id := range ids {
		if strings.EqualFold(id, "all") || strings.EqualFold(id, ruleID) {
			return true
		}
	}
	return false
}
