package report

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/hermesscan/hermesscan/internal/scanner"
)

// WriteGitHubAnnotations writes GitHub Actions workflow command annotations.
func WriteGitHubAnnotations(writer io.Writer, result scanner.Result) error {
	for _, finding := range result.Findings {
		level := githubLevel(finding.Severity)
		message := fmt.Sprintf("%s %s: %s Recommendation: %s", finding.RuleID, finding.RuleName, finding.Description, finding.Recommendation)
		if _, err := fmt.Fprintf(writer, "::%s file=%s,line=%d,col=%d,title=%s::%s\n",
			level,
			escapeGitHubAnnotationProperty(githubAnnotationPath(finding.File)),
			finding.Line,
			finding.Column,
			escapeGitHubAnnotationProperty(finding.RuleID+" "+finding.RuleName),
			escapeGitHubAnnotationMessage(message)); err != nil {
			return err
		}
	}
	return nil
}

func githubAnnotationPath(path string) string {
	normalized := filepath.ToSlash(strings.ReplaceAll(path, "\\", "/"))
	for strings.Contains(normalized, "//") {
		normalized = strings.ReplaceAll(normalized, "//", "/")
	}
	return normalized
}

func githubLevel(severity string) string {
	switch strings.ToLower(severity) {
	case "critical", "high":
		return "error"
	case "medium":
		return "warning"
	default:
		return "notice"
	}
}

func escapeGitHubAnnotationProperty(value string) string {
	value = strings.ReplaceAll(value, "%", "%25")
	value = strings.ReplaceAll(value, "\r", "%0D")
	value = strings.ReplaceAll(value, "\n", "%0A")
	value = strings.ReplaceAll(value, ":", "%3A")
	value = strings.ReplaceAll(value, ",", "%2C")
	return value
}

func escapeGitHubAnnotationMessage(value string) string {
	value = strings.ReplaceAll(value, "%", "%25")
	value = strings.ReplaceAll(value, "\r", "%0D")
	value = strings.ReplaceAll(value, "\n", "%0A")
	return value
}
