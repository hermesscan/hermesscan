package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/hermesscan/hermesscan/internal/scanner"
)

// WriteMarkdown writes the scan result as Markdown.
func WriteMarkdown(writer io.Writer, result scanner.Result) error {
	fmt.Fprintf(writer, "# HermesScan Report\n\n")
	fmt.Fprintf(writer, "| Metric | Value |\n")
	fmt.Fprintf(writer, "|---|---:|\n")
	fmt.Fprintf(writer, "| Files scanned | %d |\n", result.FilesScanned)
	fmt.Fprintf(writer, "| Rules loaded | %d |\n", result.RulesLoaded)
	fmt.Fprintf(writer, "| Findings | %d |\n", len(result.Findings))
	if result.SuppressedCount > 0 {
		fmt.Fprintf(writer, "| Suppressed | %d |\n", result.SuppressedCount)
	}
	if result.BaselineSuppressedCount > 0 {
		fmt.Fprintf(writer, "| Baseline suppressed | %d |\n", result.BaselineSuppressedCount)
	}
	fmt.Fprintf(writer, "\n")

	if len(result.Findings) == 0 {
		fmt.Fprintf(writer, "No findings.\n")
		return nil
	}

	fmt.Fprintf(writer, "## Findings\n\n")
	for _, finding := range result.Findings {
		fmt.Fprintf(writer, "### %s - %s\n\n", escapeMarkdown(finding.RuleID), escapeMarkdown(finding.RuleName))
		fmt.Fprintf(writer, "| Field | Value |\n")
		fmt.Fprintf(writer, "|---|---|\n")
		fmt.Fprintf(writer, "| Severity | %s |\n", escapeMarkdown(finding.Severity))
		fmt.Fprintf(writer, "| File | `%s:%d:%d` |\n", escapeBackticks(finding.File), finding.Line, finding.Column)
		fmt.Fprintf(writer, "| File type | %s |\n", escapeMarkdown(finding.FileType))
		if finding.Category != "" {
			fmt.Fprintf(writer, "| Category | %s |\n", escapeMarkdown(finding.Category))
		}
		if len(finding.Tags) > 0 {
			fmt.Fprintf(writer, "| Tags | %s |\n", escapeMarkdown(strings.Join(finding.Tags, ", ")))
		}
		fmt.Fprintf(writer, "| Match | `%s` |\n", escapeBackticks(finding.Match))
		fmt.Fprintf(writer, "| Description | %s |\n", escapeMarkdown(finding.Description))
		fmt.Fprintf(writer, "| Recommendation | %s |\n\n", escapeMarkdown(finding.Recommendation))
	}

	return nil
}

func escapeMarkdown(value string) string {
	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.ReplaceAll(value, "\n", " ")
	return value
}

func escapeBackticks(value string) string {
	return strings.ReplaceAll(value, "`", "'")
}
