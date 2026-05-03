package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hermesscan/hermesscan/internal/baseline"
	"github.com/hermesscan/hermesscan/internal/config"
	"github.com/hermesscan/hermesscan/internal/report"
	"github.com/hermesscan/hermesscan/internal/rules"
	"github.com/hermesscan/hermesscan/internal/scanner"
)

var version = "0.9.0"

type repeatFlag []string

func (r *repeatFlag) String() string {
	return strings.Join(*r, ",")
}

func (r *repeatFlag) Set(value string) error {
	*r = append(*r, value)
	return nil
}

type scanOptions struct {
	path              string
	rulePath          string
	configPath        string
	format            string
	outputPath        string
	failOn            string
	minSeverity       string
	summary           bool
	quiet             bool
	noColor           bool
	rulesProvided     bool
	configProvided    bool
	noFail            bool
	baselinePath      string
	createBaseline    string
	changedOnly       bool
	changedBase       string
	githubAnnotations bool
	exclude           repeatFlag
	include           repeatFlag
	category          repeatFlag
	tag               repeatFlag
	rule              repeatFlag
}

func main() {
	if len(os.Args) < 2 {
		printUsage(os.Stderr)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "scan":
		code := runScan(os.Args[2:])
		os.Exit(code)
	case "init":
		code := runInit(os.Args[2:])
		os.Exit(code)
	case "rules":
		code := runRules(os.Args[2:])
		os.Exit(code)
	case "version", "--version", "-version":
		fmt.Fprintf(os.Stdout, "HermesScan %s\n", version)
		os.Exit(0)
	case "help", "--help", "-h":
		printUsage(os.Stdout)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage(os.Stderr)
		os.Exit(2)
	}
}

func runScan(args []string) int {
	options, err := parseScanOptions(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	cfg, cfgPath, err := loadOptionalConfig(options.path, options.configPath, options.configProvided)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	options = applyConfigDefaults(options, cfg, cfgPath)

	loadedRules, err := rules.Load(options.rulePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	scanOptions := scanner.NewOptionsFromConfigValues(
		append(cfg.Exclude, options.exclude...),
		append(cfg.Include, options.include...),
		append(cfg.EnabledRules, options.rule...),
		cfg.DisabledRules,
		cfg.SeverityOverrides,
		cfg.SuppressionsEnabledValue(),
	)
	scanOptions.Categories = append(scanOptions.Categories, cfg.Categories...)
	scanOptions.Categories = append(scanOptions.Categories, options.category...)
	scanOptions.Tags = append(scanOptions.Tags, cfg.Tags...)
	scanOptions.Tags = append(scanOptions.Tags, options.tag...)
	scanOptions.ChangedOnly = options.changedOnly
	scanOptions.ChangedBase = options.changedBase

	result, err := scanner.ScanWithOptions(options.path, loadedRules, scanOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	if options.createBaseline != "" {
		if err := baseline.Save(options.createBaseline, baseline.FromResult(result)); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 2
		}
		fmt.Fprintf(os.Stderr, "Created baseline %s with %d findings.\n", options.createBaseline, len(result.Findings))
	}

	if options.baselinePath != "" {
		loadedBaseline, err := baseline.Load(options.baselinePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 2
		}
		result = baseline.Apply(result, loadedBaseline)
	}

	reportedResult := scanner.FilterByMinSeverity(result, options.minSeverity)

	if err := writeReport(options, reportedResult); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	if !options.noFail && thresholdExceeded(result, options.failOn) {
		fmt.Fprintf(os.Stderr, "HermesScan detected findings at or above severity %q.\n", options.failOn)
		return 1
	}

	return 0
}

func parseScanOptions(args []string) (scanOptions, error) {
	pathArg, flagArgs := splitPathAndFlags(args)

	set := flag.NewFlagSet("scan", flag.ContinueOnError)
	set.SetOutput(io.Discard)

	options := scanOptions{}
	set.StringVar(&options.rulePath, "rules", "", "path to JSON rule file")
	set.StringVar(&options.configPath, "config", "", "optional .hermesscan.json configuration file")
	set.StringVar(&options.format, "format", "console", "report format: console, summary, markdown, json, sarif")
	set.StringVar(&options.outputPath, "output", "", "optional output file")
	set.StringVar(&options.failOn, "fail-on", "", "fail when findings meet severity: info, low, medium, high, critical")
	set.StringVar(&options.minSeverity, "min-severity", "", "only report findings at or above severity: info, low, medium, high, critical")
	set.BoolVar(&options.summary, "summary", false, "write compact summary output")
	set.BoolVar(&options.quiet, "quiet", false, "suppress report output unless an output file is specified")
	set.BoolVar(&options.noColor, "no-color", false, "disable color output; currently accepted for compatibility")
	set.BoolVar(&options.noFail, "no-fail", false, "do not return a failing exit code even when fail-on is configured")
	set.StringVar(&options.baselinePath, "baseline", "", "path to baseline file used to ignore existing findings")
	set.StringVar(&options.createBaseline, "create-baseline", "", "write current findings to a baseline file")
	set.BoolVar(&options.changedOnly, "changed-only", false, "scan only files changed according to git")
	set.StringVar(&options.changedBase, "changed-base", "", "git base ref/commit for changed-only scans")
	set.BoolVar(&options.githubAnnotations, "github-annotations", false, "emit GitHub Actions workflow annotations")
	set.Var(&options.exclude, "exclude", "additional glob pattern to exclude; may be specified multiple times")
	set.Var(&options.include, "include", "glob pattern to include; may be specified multiple times")
	set.Var(&options.category, "category", "rule category to include; may be specified multiple times")
	set.Var(&options.tag, "tag", "rule tag to include; may be specified multiple times")
	set.Var(&options.rule, "rule", "rule ID to include; may be specified multiple times")

	if err := set.Parse(flagArgs); err != nil {
		return options, err
	}

	set.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "rules":
			options.rulesProvided = true
		case "config":
			options.configProvided = true
		}
	})

	if pathArg == "" {
		options.path = "."
	} else {
		options.path = pathArg
	}

	if len(set.Args()) > 0 {
		return options, fmt.Errorf("unexpected argument %q", set.Args()[0])
	}

	options.format = strings.ToLower(options.format)
	if options.summary {
		options.format = "summary"
	}
	if options.format != "console" && options.format != "summary" && options.format != "markdown" && options.format != "json" && options.format != "sarif" && options.format != "github" {
		return options, fmt.Errorf("unsupported format %q", options.format)
	}

	options.failOn = strings.ToLower(options.failOn)
	if options.failOn != "" && options.failOn != "none" && scanner.SeverityRank(options.failOn) == 0 {
		return options, fmt.Errorf("unsupported fail-on severity %q", options.failOn)
	}

	options.minSeverity = strings.ToLower(options.minSeverity)
	if options.minSeverity != "" && options.minSeverity != "none" && scanner.SeverityRank(options.minSeverity) == 0 {
		return options, fmt.Errorf("unsupported min-severity %q", options.minSeverity)
	}

	if options.githubAnnotations {
		options.format = "github"
	}
	if options.quiet && options.outputPath == "" && !options.githubAnnotations {
		options.format = "none"
	}

	_ = options.noColor

	return options, nil
}

func splitPathAndFlags(args []string) (string, []string) {
	var pathArg string
	var flagArgs []string
	expectsValue := false

	valueFlags := map[string]bool{
		"--rules": true, "-rules": true,
		"--config": true, "-config": true,
		"--format": true, "-format": true,
		"--output": true, "-output": true,
		"--fail-on": true, "-fail-on": true,
		"--min-severity": true, "-min-severity": true,
		"--exclude": true, "-exclude": true,
		"--include": true, "-include": true,
		"--category": true, "-category": true,
		"--tag": true, "-tag": true,
		"--rule": true, "-rule": true,
		"--baseline": true, "-baseline": true,
		"--create-baseline": true, "-create-baseline": true,
		"--changed-base": true, "-changed-base": true,
	}

	for _, arg := range args {
		if expectsValue {
			flagArgs = append(flagArgs, arg)
			expectsValue = false
			continue
		}

		if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
			flagArgs = append(flagArgs, arg)
			if valueFlags[arg] {
				expectsValue = true
			}
			continue
		}

		if pathArg == "" {
			pathArg = arg
			continue
		}

		flagArgs = append(flagArgs, arg)
	}

	return pathArg, flagArgs
}

func loadOptionalConfig(root string, path string, explicit bool) (config.Config, string, error) {
	if path == "" {
		path = config.FindDefault(root)
	}
	if path == "" {
		return config.Config{}, "", nil
	}
	cfg, err := config.Load(path)
	if err != nil {
		return config.Config{}, path, err
	}
	return cfg, path, nil
}

func applyConfigDefaults(options scanOptions, cfg config.Config, cfgPath string) scanOptions {
	baseDir := options.path
	if cfgPath != "" {
		baseDir = filepath.Dir(cfgPath)
	}

	if !options.rulesProvided {
		if cfg.Rules != "" {
			options.rulePath = resolveRelativePath(baseDir, cfg.Rules)
		} else {
			options.rulePath = defaultRulesPath()
		}
	}
	if options.rulePath == "" {
		options.rulePath = defaultRulesPath()
	}
	if options.failOn == "" && cfg.FailOn != "" {
		options.failOn = strings.ToLower(cfg.FailOn)
	}
	if options.minSeverity == "" && cfg.MinSeverity != "" {
		options.minSeverity = strings.ToLower(cfg.MinSeverity)
	}
	return options
}

func resolveRelativePath(baseDir string, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(baseDir, path)
}

func writeReport(options scanOptions, result scanner.Result) error {
	var writer io.Writer = os.Stdout
	var file *os.File

	if options.format == "none" {
		return nil
	}

	if options.outputPath != "" && options.outputPath != "-" {
		if err := ensureOutputDirectory(options.outputPath); err != nil {
			return err
		}
		created, err := os.Create(options.outputPath)
		if err != nil {
			return fmt.Errorf("create output %q: %w", options.outputPath, err)
		}
		file = created
		defer file.Close()
		writer = file
	}

	switch options.format {
	case "console":
		return report.WriteConsole(writer, result)
	case "summary":
		return report.WriteSummary(writer, result)
	case "markdown":
		return report.WriteMarkdown(writer, result)
	case "json":
		return report.WriteJSON(writer, result)
	case "sarif":
		return report.WriteSARIF(writer, result)
	case "github":
		return report.WriteGitHubAnnotations(writer, result)
	default:
		return fmt.Errorf("unsupported format %q", options.format)
	}
}

func ensureOutputDirectory(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create output directory %q: %w", dir, err)
	}
	return nil
}

func thresholdExceeded(result scanner.Result, failOn string) bool {
	for _, finding := range result.Findings {
		if scanner.MeetsThreshold(finding.Severity, failOn) {
			return true
		}
	}
	return false
}

func defaultRulesPath() string {
	candidates := []string{
		filepath.Join("rules", "hermes.rules.json"),
		filepath.Join("..", "..", "rules", "hermes.rules.json"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}

func runInit(args []string) int {
	set := flag.NewFlagSet("init", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	path := set.String("path", ".hermesscan.json", "config file path")
	overwrite := set.Bool("force", false, "overwrite an existing config")
	if err := set.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	if len(set.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "error: unexpected argument %q\n", set.Args()[0])
		return 2
	}
	if err := config.WriteDefault(*path, *overwrite); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	fmt.Fprintf(os.Stdout, "Created %s\n", *path)
	return 0
}

func runRules(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: hermesscan rules list|show|docs|validate|categories|tags [id] [--rules rules/hermes.rules.json]\n")
		return 2
	}

	switch args[0] {
	case "list":
		return runRulesList(args[1:])
	case "show":
		return runRulesShow(args[1:])
	case "docs":
		return runRulesDocs(args[1:])
	case "validate":
		return runRulesValidate(args[1:])
	case "categories":
		return runRulesCategories(args[1:])
	case "tags":
		return runRulesTags(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown rules command: %s\n", args[0])
		return 2
	}
}

func runRulesList(args []string) int {
	set := flag.NewFlagSet("rules list", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	rulePath := set.String("rules", defaultRulesPath(), "path to JSON rule file")
	if err := set.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	if len(set.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "error: unexpected argument %q\n", set.Args()[0])
		return 2
	}
	loadedRules, err := rules.Load(*rulePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	sort.Slice(loadedRules, func(i, j int) bool { return loadedRules[i].ID < loadedRules[j].ID })
	fmt.Fprintf(os.Stdout, "%-8s %-9s %-14s %s\n", "ID", "Severity", "Category", "Name")
	for _, rule := range loadedRules {
		fmt.Fprintf(os.Stdout, "%-8s %-9s %-14s %s\n", rule.ID, rule.Severity, rule.Category, rule.Name)
	}
	return 0
}

func runRulesShow(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: hermesscan rules show RULE_ID [--rules rules/hermes.rules.json]\n")
		return 2
	}
	ruleID := strings.ToUpper(args[0])
	set := flag.NewFlagSet("rules show", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	rulePath := set.String("rules", defaultRulesPath(), "path to JSON rule file")
	if err := set.Parse(args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	loadedRules, err := rules.Load(*rulePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	for _, rule := range loadedRules {
		if strings.ToUpper(rule.ID) == ruleID {
			fmt.Fprintf(os.Stdout, "ID: %s\n", rule.ID)
			fmt.Fprintf(os.Stdout, "Name: %s\n", rule.Name)
			fmt.Fprintf(os.Stdout, "Severity: %s\n", rule.Severity)
			fmt.Fprintf(os.Stdout, "Category: %s\n", rule.Category)
			if len(rule.Tags) > 0 {
				fmt.Fprintf(os.Stdout, "Tags: %s\n", strings.Join(rule.Tags, ", "))
			}
			fmt.Fprintf(os.Stdout, "File types: %s\n", strings.Join(rule.FileTypes, ", "))
			fmt.Fprintf(os.Stdout, "Pattern: %s\n", rule.Pattern)
			if rule.ExcludePattern != "" {
				fmt.Fprintf(os.Stdout, "Exclude pattern: %s\n", rule.ExcludePattern)
			}
			if rule.ContextBeforePattern != "" {
				fmt.Fprintf(os.Stdout, "Context before pattern: %s\n", rule.ContextBeforePattern)
				fmt.Fprintf(os.Stdout, "Context before lines: %d\n", rule.ContextBeforeLines)
			}
			if rule.RequiredFilePattern != "" {
				fmt.Fprintf(os.Stdout, "Required file pattern: %s\n", rule.RequiredFilePattern)
			}
			fmt.Fprintf(os.Stdout, "Description: %s\n", rule.Description)
			fmt.Fprintf(os.Stdout, "Recommendation: %s\n", rule.Recommendation)
			return 0
		}
	}
	fmt.Fprintf(os.Stderr, "error: rule %q was not found\n", ruleID)
	return 1
}

func runRulesDocs(args []string) int {
	set := flag.NewFlagSet("rules docs", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	rulePath := set.String("rules", defaultRulesPath(), "path to JSON rule file")
	outputPath := set.String("output", "", "optional Markdown output file")
	if err := set.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	if len(set.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "error: unexpected argument %q\n", set.Args()[0])
		return 2
	}
	loadedRules, err := rules.Load(*rulePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	sort.Slice(loadedRules, func(i, j int) bool { return loadedRules[i].ID < loadedRules[j].ID })
	var builder strings.Builder
	writeRulesMarkdown(&builder, loadedRules)
	if *outputPath == "" || *outputPath == "-" {
		fmt.Fprint(os.Stdout, builder.String())
		return 0
	}
	if err := ensureOutputDirectory(*outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	if err := os.WriteFile(*outputPath, []byte(builder.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: write rules docs %q: %v\n", *outputPath, err)
		return 2
	}
	fmt.Fprintf(os.Stdout, "Wrote %s\n", *outputPath)
	return 0
}

func runRulesValidate(args []string) int {
	set := flag.NewFlagSet("rules validate", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	rulePath := set.String("rules", defaultRulesPath(), "path to JSON rule file")
	if err := set.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	if len(set.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "error: unexpected argument %q\n", set.Args()[0])
		return 2
	}
	loadedRules, err := rules.Load(*rulePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	if err := rules.ValidateCatalog(loadedRules); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	source := *rulePath
	if source == "" {
		source = "embedded default rules"
	}
	fmt.Fprintf(os.Stdout, "Validated %d rules from %s\n", len(loadedRules), source)
	return 0
}

func runRulesCategories(args []string) int {
	return runRulesValueList(args, "categories")
}

func runRulesTags(args []string) int {
	return runRulesValueList(args, "tags")
}

func runRulesValueList(args []string, mode string) int {
	set := flag.NewFlagSet("rules "+mode, flag.ContinueOnError)
	set.SetOutput(io.Discard)
	rulePath := set.String("rules", defaultRulesPath(), "path to JSON rule file")
	if err := set.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	loadedRules, err := rules.Load(*rulePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	values := map[string]bool{}
	for _, rule := range loadedRules {
		if mode == "categories" {
			value := strings.TrimSpace(rule.Category)
			if value != "" {
				values[value] = true
			}
			continue
		}
		for _, tag := range rule.Tags {
			value := strings.TrimSpace(tag)
			if value != "" {
				values[value] = true
			}
		}
	}
	keys := make([]string, 0, len(values))
	for value := range values {
		keys = append(keys, value)
	}
	sort.Strings(keys)
	for _, value := range keys {
		fmt.Fprintln(os.Stdout, value)
	}
	return 0
}

func writeRulesMarkdown(writer io.Writer, loadedRules []rules.Rule) {
	fmt.Fprintln(writer, "# HermesScan Rule Reference")
	fmt.Fprintln(writer)
	fmt.Fprintf(writer, "Generated for HermesScan %s.\n\n", version)
	fmt.Fprintln(writer, "| ID | Severity | Category | Name |")
	fmt.Fprintln(writer, "|---|---|---|---|")
	for _, rule := range loadedRules {
		fmt.Fprintf(writer, "| `%s` | %s | %s | %s |\n", rule.ID, rule.Severity, rule.Category, markdownEscape(rule.Name))
	}
	fmt.Fprintln(writer)
	for _, rule := range loadedRules {
		fmt.Fprintf(writer, "## %s - %s\n\n", rule.ID, rule.Name)
		fmt.Fprintf(writer, "**Severity:** %s\n\n", rule.Severity)
		fmt.Fprintf(writer, "**Category:** %s\n\n", rule.Category)
		if len(rule.Tags) > 0 {
			fmt.Fprintf(writer, "**Tags:** `%s`\n\n", strings.Join(rule.Tags, "`, `"))
		}
		fmt.Fprintf(writer, "**File types:** `%s`\n\n", strings.Join(rule.FileTypes, "`, `"))
		fmt.Fprintf(writer, "%s\n\n", rule.Description)
		fmt.Fprintf(writer, "**Recommendation:** %s\n\n", rule.Recommendation)
		fmt.Fprintln(writer, "```text")
		fmt.Fprintln(writer, rule.Pattern)
		fmt.Fprintln(writer, "```")
		fmt.Fprintln(writer)
		if rule.ExcludePattern != "" {
			fmt.Fprintln(writer, "**Exclude pattern:**")
			fmt.Fprintln(writer)
			fmt.Fprintln(writer, "```text")
			fmt.Fprintln(writer, rule.ExcludePattern)
			fmt.Fprintln(writer, "```")
			fmt.Fprintln(writer)
		}
		if rule.ContextBeforePattern != "" {
			fmt.Fprintf(writer, "**Context before pattern:** within the previous %d line(s)\n\n", rule.ContextBeforeLines)
			fmt.Fprintln(writer, "```text")
			fmt.Fprintln(writer, rule.ContextBeforePattern)
			fmt.Fprintln(writer, "```")
			fmt.Fprintln(writer)
		}
		if rule.RequiredFilePattern != "" {
			fmt.Fprintln(writer, "**Required file pattern:**")
			fmt.Fprintln(writer)
			fmt.Fprintln(writer, "```text")
			fmt.Fprintln(writer, rule.RequiredFilePattern)
			fmt.Fprintln(writer, "```")
			fmt.Fprintln(writer)
		}
	}
}

func markdownEscape(value string) string {
	value = strings.ReplaceAll(value, "|", "\\|")
	return value
}

func printUsage(writer io.Writer) {
	fmt.Fprintf(writer, "HermesScan %s\n\n", version)
	fmt.Fprintf(writer, "Usage:\n")
	fmt.Fprintf(writer, "  hermesscan scan [path] [--config .hermesscan.json] [--rules rules/hermes.rules.json] [--format console|summary|markdown|json|sarif] [--output file] [--fail-on high] [--no-fail] [--baseline file] [--create-baseline file] [--min-severity medium] [--include pattern] [--exclude pattern] [--category name] [--tag name] [--rule HMS0001] [--changed-only] [--changed-base ref] [--github-annotations] [--summary] [--quiet] [--no-color]\n")
	fmt.Fprintf(writer, "  hermesscan init [--path .hermesscan.json] [--force]\n")
	fmt.Fprintf(writer, "  hermesscan rules list [--rules rules/hermes.rules.json]\n")
	fmt.Fprintf(writer, "  hermesscan rules show RULE_ID [--rules rules/hermes.rules.json]\n")
	fmt.Fprintf(writer, "  hermesscan rules docs [--rules rules/hermes.rules.json] [--output docs/rules.md]\n")
	fmt.Fprintf(writer, "  hermesscan rules validate [--rules rules/hermes.rules.json]\n")
	fmt.Fprintf(writer, "  hermesscan rules categories [--rules rules/hermes.rules.json]\n")
	fmt.Fprintf(writer, "  hermesscan rules tags [--rules rules/hermes.rules.json]\n")
	fmt.Fprintf(writer, "  hermesscan version\n")
}
