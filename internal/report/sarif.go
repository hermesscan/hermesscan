package report

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/hermesscan/hermesscan/internal/scanner"
)

type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri,omitempty"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	ShortDescription sarifMessage      `json:"shortDescription"`
	FullDescription  sarifMessage      `json:"fullDescription"`
	Help             sarifMessage      `json:"help"`
	Properties       sarifRuleProperty `json:"properties"`
}

type sarifRuleProperty struct {
	ProblemSeverity string   `json:"problem.severity"`
	Category        string   `json:"category,omitempty"`
	Tags            []string `json:"tags,omitempty"`
}

type sarifResult struct {
	RuleID              string            `json:"ruleId"`
	Level               string            `json:"level"`
	Message             sarifMessage      `json:"message"`
	Locations           []sarifLocation   `json:"locations"`
	PartialFingerprints map[string]string `json:"partialFingerprints,omitempty"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
}

// WriteSARIF writes a SARIF 2.1.0 report compatible with GitHub code scanning upload.
func WriteSARIF(writer io.Writer, result scanner.Result) error {
	ruleMap := make(map[string]scanner.Finding)
	orderedRuleIDs := make([]string, 0)

	for _, finding := range result.Findings {
		if _, exists := ruleMap[finding.RuleID]; !exists {
			ruleMap[finding.RuleID] = finding
			orderedRuleIDs = append(orderedRuleIDs, finding.RuleID)
		}
	}

	sarifRules := make([]sarifRule, 0, len(orderedRuleIDs))
	for _, ruleID := range orderedRuleIDs {
		finding := ruleMap[ruleID]
		sarifRules = append(sarifRules, sarifRule{
			ID:               finding.RuleID,
			Name:             finding.RuleName,
			ShortDescription: sarifMessage{Text: finding.RuleName},
			FullDescription:  sarifMessage{Text: finding.Description},
			Help:             sarifMessage{Text: finding.Recommendation},
			Properties: sarifRuleProperty{
				ProblemSeverity: strings.ToLower(finding.Severity),
				Category:        finding.Category,
				Tags:            append([]string{}, finding.Tags...),
			},
		})
	}

	results := make([]sarifResult, 0, len(result.Findings))
	for _, finding := range result.Findings {
		results = append(results, sarifResult{
			RuleID: finding.RuleID,
			Level:  sarifLevel(finding.Severity),
			Message: sarifMessage{
				Text: finding.RuleName + ": " + finding.Description + " Recommendation: " + finding.Recommendation,
			},
			PartialFingerprints: map[string]string{"hermesscan/v1": finding.Fingerprint},
			Locations: []sarifLocation{
				{
					PhysicalLocation: sarifPhysicalLocation{
						ArtifactLocation: sarifArtifactLocation{URI: normalizePathForSARIF(finding.File)},
						Region: sarifRegion{
							StartLine:   finding.Line,
							StartColumn: finding.Column,
						},
					},
				},
			},
		})
	}

	log := sarifLog{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:  "HermesScan",
						Rules: sarifRules,
					},
				},
				Results: results,
			},
		},
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(log)
}

func sarifLevel(severity string) string {
	switch scanner.SeverityRank(severity) {
	case 5, 4:
		return "error"
	case 3:
		return "warning"
	case 2:
		return "note"
	default:
		return "none"
	}
}

func normalizePathForSARIF(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
