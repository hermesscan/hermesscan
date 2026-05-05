package rules

// Rule describes a single regex-based scan rule.
type Rule struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	Severity             string   `json:"severity"`
	Category             string   `json:"category,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	FileTypes            []string `json:"fileTypes"`
	Pattern              string   `json:"pattern"`
	ExcludePattern       string   `json:"excludePattern,omitempty"`
	ContextBeforePattern string   `json:"contextBeforePattern,omitempty"`
	ContextBeforeLines   int      `json:"contextBeforeLines,omitempty"`
	TriggerFilePattern   string   `json:"triggerFilePattern,omitempty"`
	RequiredFilePattern  string   `json:"requiredFilePattern,omitempty"`
	Description          string   `json:"description"`
	Recommendation       string   `json:"recommendation"`
}
