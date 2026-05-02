package rules

// Rule describes a single regex-based scan rule.
type Rule struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Severity       string   `json:"severity"`
	Category       string   `json:"category,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	FileTypes      []string `json:"fileTypes"`
	Pattern        string   `json:"pattern"`
	Description    string   `json:"description"`
	Recommendation string   `json:"recommendation"`
}
