package scanner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hermesscan/hermesscan/internal/files"
	"github.com/hermesscan/hermesscan/internal/rules"
)

type compiledRule struct {
	rule               rules.Rule
	regex              *regexp.Regexp
	excludeRegex       *regexp.Regexp
	contextBeforeRegex *regexp.Regexp
}

// Scan applies rules to candidate files under root using default options.
func Scan(root string, loadedRules []rules.Rule) (Result, error) {
	return ScanWithOptions(root, loadedRules, DefaultOptions())
}

// ScanWithOptions applies rules to candidate files under root using explicit options.
func ScanWithOptions(root string, loadedRules []rules.Rule, options Options) (Result, error) {
	candidates, err := files.Discover(root)
	if err != nil {
		return Result{}, fmt.Errorf("discover files: %w", err)
	}
	candidates = filterIncludedCandidates(candidates, options.Include)
	if options.ChangedOnly {
		changed, err := changedFiles(root, options.ChangedBase)
		if err != nil {
			return Result{}, fmt.Errorf("discover changed files: %w", err)
		}
		candidates = filterChangedCandidates(candidates, changed)
	}
	candidates = filterExcludedCandidates(candidates, options.Exclude)

	compiled, err := compileRules(loadedRules, options)
	if err != nil {
		return Result{}, err
	}

	result := Result{
		Root:            root,
		FilesScanned:    len(candidates),
		RulesLoaded:     len(compiled),
		Findings:        []Finding{},
		SuppressedCount: 0,
	}

	for _, candidate := range candidates {
		findings, suppressed, err := scanFile(candidate, compiled, options)
		if err != nil {
			return Result{}, err
		}
		result.Findings = append(result.Findings, findings...)
		result.SuppressedCount += suppressed
	}

	return result, nil
}

func compileRules(loadedRules []rules.Rule, options Options) ([]compiledRule, error) {
	compiled := make([]compiledRule, 0, len(loadedRules))
	for _, rule := range loadedRules {
		ruleID := strings.ToUpper(strings.TrimSpace(rule.ID))
		if len(options.EnabledRules) > 0 && !options.EnabledRules[ruleID] {
			continue
		}
		if options.DisabledRules != nil && options.DisabledRules[ruleID] {
			continue
		}
		if !ruleMatchesCategoryAndTag(rule, options) {
			continue
		}
		if options.SeverityOverrides != nil {
			if severity, ok := options.SeverityOverrides[strings.ToUpper(rule.ID)]; ok && severity != "" {
				rule.Severity = severity
			}
		}
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return nil, fmt.Errorf("compile rule %q: %w", rule.ID, err)
		}
		var excludeRe *regexp.Regexp
		if strings.TrimSpace(rule.ExcludePattern) != "" {
			excludeRe, err = regexp.Compile(rule.ExcludePattern)
			if err != nil {
				return nil, fmt.Errorf("compile rule %q exclude pattern: %w", rule.ID, err)
			}
		}
		var contextBeforeRe *regexp.Regexp
		if strings.TrimSpace(rule.ContextBeforePattern) != "" {
			contextBeforeRe, err = regexp.Compile(rule.ContextBeforePattern)
			if err != nil {
				return nil, fmt.Errorf("compile rule %q context-before pattern: %w", rule.ID, err)
			}
		}
		compiled = append(compiled, compiledRule{rule: rule, regex: re, excludeRegex: excludeRe, contextBeforeRegex: contextBeforeRe})
	}
	return compiled, nil
}

func scanFile(candidate files.Candidate, compiled []compiledRule, options Options) ([]Finding, int, error) {
	file, err := os.Open(candidate.Path)
	if err != nil {
		return nil, 0, fmt.Errorf("open %q: %w", candidate.Path, err)
	}
	defer file.Close()

	var findings []Finding
	suppressedCount := 0
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	fileSuppressions := []string{}
	nextLineSuppressions := make(map[int][]string)
	previousLines := []string{}

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		lineSuppressions := []string{}

		if options.SuppressionsEnabled {
			for _, directive := range parseSuppression(line) {
				switch directive.Kind {
				case "disable-file":
					fileSuppressions = append(fileSuppressions, directive.RuleIDs...)
				case "disable-next-line":
					nextLineSuppressions[lineNumber+1] = append(nextLineSuppressions[lineNumber+1], directive.RuleIDs...)
				case "disable-line":
					lineSuppressions = append(lineSuppressions, directive.RuleIDs...)
				}
			}
		}

		for _, item := range compiled {
			if !ruleAppliesToType(item.rule, candidate.Type) {
				continue
			}

			location := item.regex.FindStringIndex(line)
			if location == nil {
				continue
			}

			if hasRuleExclusionContext(item, line, previousLines) {
				continue
			}

			if options.SuppressionsEnabled && isSuppressed(item.rule.ID, fileSuppressions, nextLineSuppressions[lineNumber], lineSuppressions) {
				suppressedCount++
				continue
			}

			matchedText := line[location[0]:location[1]]
			finding := Finding{
				RuleID:         item.rule.ID,
				RuleName:       item.rule.Name,
				Severity:       item.rule.Severity,
				File:           candidate.Path,
				Line:           lineNumber,
				Column:         location[0] + 1,
				FileType:       candidate.Type,
				Match:          matchedText,
				Description:    item.rule.Description,
				Recommendation: item.rule.Recommendation,
				Category:       item.rule.Category,
				Tags:           append([]string{}, item.rule.Tags...),
			}
			finding.Fingerprint = Fingerprint(finding)
			findings = append(findings, finding)
		}

		previousLines = append(previousLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, suppressedCount, fmt.Errorf("scan %q: %w", candidate.Path, err)
	}

	return findings, suppressedCount, nil
}

func hasRuleExclusionContext(item compiledRule, line string, previousLines []string) bool {
	if item.excludeRegex != nil && item.excludeRegex.MatchString(line) {
		return true
	}
	if item.contextBeforeRegex == nil {
		return false
	}
	window := item.rule.ContextBeforeLines
	if window < 1 {
		window = 1
	}
	start := len(previousLines) - window
	if start < 0 {
		start = 0
	}
	for _, previousLine := range previousLines[start:] {
		if item.contextBeforeRegex.MatchString(previousLine) {
			return true
		}
	}
	return false
}

func isSuppressed(ruleID string, groups ...[]string) bool {
	for _, group := range groups {
		if suppressionMatches(group, ruleID) {
			return true
		}
	}
	return false
}

func ruleAppliesToType(rule rules.Rule, fileType string) bool {
	if len(rule.FileTypes) == 0 {
		return true
	}
	for _, value := range rule.FileTypes {
		if strings.EqualFold(value, fileType) || value == "*" {
			return true
		}
	}
	return false
}

func filterIncludedCandidates(candidates []files.Candidate, patterns []string) []files.Candidate {
	if len(patterns) == 0 {
		return candidates
	}
	filtered := make([]files.Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidateMatchesAny(candidate.Path, patterns) {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}

func filterExcludedCandidates(candidates []files.Candidate, patterns []string) []files.Candidate {
	if len(patterns) == 0 {
		return candidates
	}
	filtered := make([]files.Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidateMatchesAny(candidate.Path, patterns) {
			continue
		}
		filtered = append(filtered, candidate)
	}
	return filtered
}

func filterChangedCandidates(candidates []files.Candidate, changed map[string]bool) []files.Candidate {
	if len(changed) == 0 {
		return []files.Candidate{}
	}
	filtered := make([]files.Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		clean := filepath.ToSlash(candidate.Path)
		if changed[clean] || changed[filepath.Base(clean)] {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}

func changedFiles(root string, base string) (map[string]bool, error) {
	args := []string{"-C", root, "diff", "--name-only"}
	if base != "" {
		args = append(args, base)
	} else {
		args = append(args, "HEAD")
	}
	output, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool)
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(filepath.ToSlash(line))
		if line == "" {
			continue
		}
		result[line] = true
		result[filepath.ToSlash(filepath.Join(root, line))] = true
	}
	return result, nil
}

func ruleMatchesCategoryAndTag(rule rules.Rule, options Options) bool {
	categories := normalizeStringSet(options.Categories)
	if len(categories) > 0 && !categories[strings.ToLower(strings.TrimSpace(rule.Category))] {
		return false
	}
	tags := normalizeStringSet(options.Tags)
	if len(tags) > 0 {
		matched := false
		for _, tag := range rule.Tags {
			if tags[strings.ToLower(strings.TrimSpace(tag))] {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

func candidateMatchesAny(path string, patterns []string) bool {
	clean := filepath.ToSlash(path)
	base := filepath.Base(clean)
	for _, pattern := range patterns {
		pattern = strings.TrimSpace(filepath.ToSlash(pattern))
		if pattern == "" {
			continue
		}
		if strings.HasSuffix(pattern, "/**") {
			prefix := strings.TrimSuffix(pattern, "/**")
			if strings.HasPrefix(clean, prefix+"/") || clean == prefix {
				return true
			}
		}
		if ok, _ := filepath.Match(pattern, clean); ok {
			return true
		}
		if ok, _ := filepath.Match(pattern, base); ok {
			return true
		}
	}
	return false
}

// NewOptionsFromConfigValues creates Options from config-shaped values.
func NewOptionsFromConfigValues(exclude []string, include []string, enabled []string, disabled []string, overrides map[string]string, suppressionsEnabled bool) Options {
	normalizedOverrides := make(map[string]string)
	for key, value := range overrides {
		normalizedOverrides[strings.ToUpper(strings.TrimSpace(key))] = value
	}
	return Options{
		Exclude:             exclude,
		Include:             include,
		EnabledRules:        normalizeRuleSet(enabled),
		DisabledRules:       normalizeRuleSet(disabled),
		SeverityOverrides:   normalizedOverrides,
		SuppressionsEnabled: suppressionsEnabled,
	}
}
