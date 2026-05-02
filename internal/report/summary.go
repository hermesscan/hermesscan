package report

import (
	"fmt"
	"io"

	"github.com/hermesscan/hermesscan/internal/scanner"
)

// WriteSummary writes a compact human-readable report.
func WriteSummary(writer io.Writer, result scanner.Result) error {
	counts := scanner.SeverityCounts(result)

	_, err := fmt.Fprintf(writer, "HermesScan: %d findings across %d files\n", len(result.Findings), result.FilesScanned)
	if err != nil {
		return err
	}

	fmt.Fprintf(writer, "Critical: %d\n", counts["Critical"])
	fmt.Fprintf(writer, "High: %d\n", counts["High"])
	fmt.Fprintf(writer, "Medium: %d\n", counts["Medium"])
	fmt.Fprintf(writer, "Low: %d\n", counts["Low"])
	fmt.Fprintf(writer, "Info: %d\n", counts["Info"])
	if result.SuppressedCount > 0 {
		fmt.Fprintf(writer, "Suppressed: %d\n", result.SuppressedCount)
	}
	if result.BaselineSuppressedCount > 0 {
		fmt.Fprintf(writer, "Baseline suppressed: %d\n", result.BaselineSuppressedCount)
	}

	return nil
}
