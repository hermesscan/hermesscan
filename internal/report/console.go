package report

import (
	"fmt"
	"io"

	"github.com/hermesscan/hermesscan/internal/scanner"
)

// WriteConsole writes a human-readable report.
func WriteConsole(writer io.Writer, result scanner.Result) error {
	_, err := fmt.Fprintf(writer, "HermesScan Report\n")
	if err != nil {
		return err
	}
	fmt.Fprintf(writer, "=================\n\n")
	fmt.Fprintf(writer, "Root: %s\n", result.Root)
	fmt.Fprintf(writer, "Files scanned: %d\n", result.FilesScanned)
	fmt.Fprintf(writer, "Rules loaded: %d\n", result.RulesLoaded)
	fmt.Fprintf(writer, "Findings: %d\n", len(result.Findings))
	if result.SuppressedCount > 0 {
		fmt.Fprintf(writer, "Suppressed: %d\n", result.SuppressedCount)
	}
	if result.BaselineSuppressedCount > 0 {
		fmt.Fprintf(writer, "Baseline suppressed: %d\n", result.BaselineSuppressedCount)
	}
	fmt.Fprintf(writer, "\n")

	if len(result.Findings) == 0 {
		fmt.Fprintf(writer, "No findings.\n")
		return nil
	}

	for _, finding := range result.Findings {
		fmt.Fprintf(writer, "[%s] %s - %s\n", finding.Severity, finding.RuleID, finding.RuleName)
		fmt.Fprintf(writer, "  File: %s:%d:%d\n", finding.File, finding.Line, finding.Column)
		if finding.Category != "" {
			fmt.Fprintf(writer, "  Category: %s\n", finding.Category)
		}
		fmt.Fprintf(writer, "  Match: %s\n", finding.Match)
		fmt.Fprintf(writer, "  Why: %s\n", finding.Description)
		fmt.Fprintf(writer, "  Fix: %s\n\n", finding.Recommendation)
	}

	return nil
}
