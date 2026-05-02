package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/hermesscan/hermesscan/internal/scanner"
)

// File is the persisted HermesScan baseline format.
type File struct {
	Version   int     `json:"version"`
	CreatedAt string  `json:"createdAt"`
	Findings  []Entry `json:"findings"`
}

// Entry identifies a finding accepted into a baseline.
type Entry struct {
	Fingerprint string `json:"fingerprint"`
	RuleID      string `json:"ruleId"`
	Severity    string `json:"severity"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Match       string `json:"match"`
}

// Load reads a baseline file.
func Load(path string) (File, error) {
	if path == "" {
		return File{}, fmt.Errorf("baseline path is required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return File{}, fmt.Errorf("read baseline %q: %w", path, err)
	}
	var file File
	if err := json.Unmarshal(data, &file); err != nil {
		return File{}, fmt.Errorf("parse baseline %q: %w", path, err)
	}
	return file, nil
}

// FromResult creates a baseline file from a scan result.
func FromResult(result scanner.Result) File {
	entries := make([]Entry, 0, len(result.Findings))
	seen := make(map[string]bool)
	for _, finding := range result.Findings {
		fp := finding.Fingerprint
		if fp == "" {
			fp = scanner.Fingerprint(finding)
		}
		if seen[fp] {
			continue
		}
		seen[fp] = true
		entries = append(entries, Entry{
			Fingerprint: fp,
			RuleID:      finding.RuleID,
			Severity:    finding.Severity,
			File:        filepath.ToSlash(finding.File),
			Line:        finding.Line,
			Match:       finding.Match,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].File == entries[j].File {
			if entries[i].Line == entries[j].Line {
				return entries[i].RuleID < entries[j].RuleID
			}
			return entries[i].Line < entries[j].Line
		}
		return entries[i].File < entries[j].File
	})
	return File{Version: 1, CreatedAt: time.Now().UTC().Format(time.RFC3339), Findings: entries}
}

// Save writes a baseline file.
func Save(path string, file File) error {
	if path == "" {
		return fmt.Errorf("baseline path is required")
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create baseline directory %q: %w", dir, err)
		}
	}
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize baseline: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write baseline %q: %w", path, err)
	}
	return nil
}

// Apply removes baseline-matching findings and increments Result.BaselineSuppressedCount.
func Apply(result scanner.Result, file File) scanner.Result {
	known := make(map[string]bool)
	for _, entry := range file.Findings {
		if entry.Fingerprint != "" {
			known[entry.Fingerprint] = true
		}
	}
	kept := make([]scanner.Finding, 0, len(result.Findings))
	removed := 0
	for _, finding := range result.Findings {
		fp := finding.Fingerprint
		if fp == "" {
			fp = scanner.Fingerprint(finding)
		}
		if known[fp] {
			removed++
			continue
		}
		kept = append(kept, finding)
	}
	result.Findings = kept
	result.BaselineSuppressedCount += removed
	return result
}
