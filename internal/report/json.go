package report

import (
	"encoding/json"
	"io"

	"github.com/hermesscan/hermesscan/internal/scanner"
)

// WriteJSON writes the scan result as indented JSON.
func WriteJSON(writer io.Writer, result scanner.Result) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}
