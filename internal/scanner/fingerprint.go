package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
)

// Fingerprint returns a stable-enough identifier for a finding.
func Fingerprint(f Finding) string {
	path := filepath.ToSlash(strings.ToLower(f.File))
	match := strings.TrimSpace(f.Match)
	payload := fmt.Sprintf("%s|%s|%d|%s", strings.ToUpper(f.RuleID), path, f.Line, match)
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
