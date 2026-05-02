package scanner

// Finding is a single rule match discovered in a scanned file.
type Finding struct {
	RuleID         string   `json:"ruleId"`
	RuleName       string   `json:"ruleName"`
	Severity       string   `json:"severity"`
	Category       string   `json:"category,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	File           string   `json:"file"`
	Line           int      `json:"line"`
	Column         int      `json:"column"`
	FileType       string   `json:"fileType"`
	Match          string   `json:"match"`
	Description    string   `json:"description"`
	Recommendation string   `json:"recommendation"`
	Fingerprint    string   `json:"fingerprint,omitempty"`
}

// Result contains aggregate scan output.
type Result struct {
	Root                    string    `json:"root"`
	FilesScanned            int       `json:"filesScanned"`
	RulesLoaded             int       `json:"rulesLoaded"`
	Findings                []Finding `json:"findings"`
	SuppressedCount         int       `json:"suppressedCount,omitempty"`
	BaselineSuppressedCount int       `json:"baselineSuppressedCount,omitempty"`
}
