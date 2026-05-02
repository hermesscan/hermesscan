package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hermesscan/hermesscan/internal/rules"
)

func TestDefaultRulePrecisionPostgresPortRequiresExposureContext(t *testing.T) {
	noRisk := scanDefaultRuleFixture(t, "workflow.yml", "env:\n  POSTGRES_DOC: \"PostgreSQL normally listens on 5432\"\n", "HMS0002")
	if len(noRisk.Findings) != 0 {
		t.Fatalf("expected no bare-port findings, got %#v", noRisk.Findings)
	}

	risk := scanDefaultRuleFixture(t, "workflow.yml", "services:\n  postgres:\n    ports:\n      - \"5432:5432\"\n", "HMS0002")
	if len(risk.Findings) != 1 {
		t.Fatalf("expected one exposed-port finding, got %#v", risk.Findings)
	}
}

func TestDefaultRulePrecisionMutableActionAllowsVersionTags(t *testing.T) {
	noRisk := scanDefaultRuleFixture(t, "workflow.yml", "steps:\n  - uses: actions/checkout@v4\n", "HMS0009")
	if len(noRisk.Findings) != 0 {
		t.Fatalf("expected no version-tag findings, got %#v", noRisk.Findings)
	}

	risk := scanDefaultRuleFixture(t, "workflow.yml", "steps:\n  - uses: example/action@main\n", "HMS0009")
	if len(risk.Findings) != 1 {
		t.Fatalf("expected one mutable-action finding, got %#v", risk.Findings)
	}
}

func TestDefaultRulePrecisionBroadCacheKeySpecificity(t *testing.T) {
	bareRunnerOS := scanDefaultRuleFixture(t, "workflow.yml", "key: ${{ runner.os }}\n", "HMS0016")
	if len(bareRunnerOS.Findings) != 1 {
		t.Fatalf("expected one bare runner OS cache-key finding, got %#v", bareRunnerOS.Findings)
	}

	staticPrefix := scanDefaultRuleFixture(t, "workflow.yml", "key: npm-${{ runner.os }}\n", "HMS0016")
	if len(staticPrefix.Findings) != 1 {
		t.Fatalf("expected one static-prefix cache-key finding, got %#v", staticPrefix.Findings)
	}

	lockfileHash := scanDefaultRuleFixture(t, "workflow.yml", "key: ${{ runner.os }}-${{ hashFiles('**/go.sum') }}\n", "HMS0016")
	if len(lockfileHash.Findings) != 0 {
		t.Fatalf("expected no lockfile-hash cache-key findings, got %#v", lockfileHash.Findings)
	}

	matrixAndLockfileHash := scanDefaultRuleFixture(t, "workflow.yml", "key: npm-${{ runner.os }}-${{ matrix.node-version }}-${{ hashFiles('**/package-lock.json') }}\n", "HMS0016")
	if len(matrixAndLockfileHash.Findings) != 0 {
		t.Fatalf("expected no matrix/lockfile cache-key findings, got %#v", matrixAndLockfileHash.Findings)
	}

	refSpecific := scanDefaultRuleFixture(t, "workflow.yml", "key: deps-${{ runner.os }}-${{ github.ref }}\n", "HMS0016")
	if len(refSpecific.Findings) != 0 {
		t.Fatalf("expected no ref-specific cache-key findings, got %#v", refSpecific.Findings)
	}

	lockfileName := scanDefaultRuleFixture(t, "workflow.yml", "key: pip-${{ runner.os }}-requirements.txt\n", "HMS0016")
	if len(lockfileName.Findings) != 0 {
		t.Fatalf("expected no lockfile-name cache-key findings, got %#v", lockfileName.Findings)
	}
}

func TestDefaultRulePrecisionPackageInstallCacheIsolation(t *testing.T) {
	npmRisk := scanDefaultRuleFixture(t, "ci.sh", "npm ci\n", "HMS0010")
	if len(npmRisk.Findings) != 1 {
		t.Fatalf("expected one npm install cache finding, got %#v", npmRisk.Findings)
	}

	pipRisk := scanDefaultRuleFixture(t, "ci.sh", "pip install -r requirements.txt\n", "HMS0010")
	if len(pipRisk.Findings) != 1 {
		t.Fatalf("expected one pip install cache finding, got %#v", pipRisk.Findings)
	}

	npmCache := scanDefaultRuleFixture(t, "ci.sh", "npm ci --cache \"$RUNNER_TEMP/npm-cache\"\n", "HMS0010")
	if len(npmCache.Findings) != 0 {
		t.Fatalf("expected no npm cache-isolated findings, got %#v", npmCache.Findings)
	}

	pipInlineCache := scanDefaultRuleFixture(t, "ci.sh", "PIP_CACHE_DIR=\"$RUNNER_TEMP/pip-cache\" pip install -r requirements.txt\n", "HMS0010")
	if len(pipInlineCache.Findings) != 0 {
		t.Fatalf("expected no inline pip cache-isolated findings, got %#v", pipInlineCache.Findings)
	}

	yarnInlineCache := scanDefaultRuleFixture(t, "ci.sh", "YARN_CACHE_FOLDER=\"$RUNNER_TEMP/yarn-cache\" yarn install\n", "HMS0010")
	if len(yarnInlineCache.Findings) != 0 {
		t.Fatalf("expected no inline yarn cache-isolated findings, got %#v", yarnInlineCache.Findings)
	}

	pipNearbyCache := scanDefaultRuleFixture(t, "ci.sh", "export PIP_CACHE_DIR=\"$RUNNER_TEMP/pip-cache\"\npip install -r requirements.txt\n", "HMS0010")
	if len(pipNearbyCache.Findings) != 0 {
		t.Fatalf("expected no nearby pip cache-isolated findings, got %#v", pipNearbyCache.Findings)
	}

	yamlCache := scanDefaultRuleFixture(t, "workflow.yml", "env:\n  PIP_CACHE_DIR: ${{ runner.temp }}/pip-cache\nsteps:\n  - run: pip install -r requirements.txt\n", "HMS0010")
	if len(yamlCache.Findings) != 0 {
		t.Fatalf("expected no YAML env cache-isolated findings, got %#v", yamlCache.Findings)
	}
}

func TestDefaultRulePrecisionDockerComposeProjectName(t *testing.T) {
	risk := scanDefaultRuleFixture(t, "ci.sh", "docker compose up -d\n", "HMS0007")
	if len(risk.Findings) != 1 {
		t.Fatalf("expected one compose project-name finding, got %#v", risk.Findings)
	}

	projectFlag := scanDefaultRuleFixture(t, "ci.sh", "docker compose -p \"$CI_RUN_ID\" up -d\n", "HMS0007")
	if len(projectFlag.Findings) != 0 {
		t.Fatalf("expected no findings with compose project flag, got %#v", projectFlag.Findings)
	}

	inlineProjectName := scanDefaultRuleFixture(t, "ci.sh", "COMPOSE_PROJECT_NAME=\"$CI_RUN_ID\" docker compose up -d\n", "HMS0007")
	if len(inlineProjectName.Findings) != 0 {
		t.Fatalf("expected no findings with inline COMPOSE_PROJECT_NAME, got %#v", inlineProjectName.Findings)
	}

	nearbyProjectName := scanDefaultRuleFixture(t, "ci.sh", "export COMPOSE_PROJECT_NAME=\"$CI_RUN_ID\"\ndocker compose up -d\n", "HMS0007")
	if len(nearbyProjectName.Findings) != 0 {
		t.Fatalf("expected no findings with nearby COMPOSE_PROJECT_NAME, got %#v", nearbyProjectName.Findings)
	}
}

func scanDefaultRuleFixture(t *testing.T, filename string, content string, ruleID string) Result {
	t.Helper()

	root := t.TempDir()
	path := filepath.Join(root, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	loaded, err := rules.LoadDefault()
	if err != nil {
		t.Fatalf("LoadDefault returned error: %v", err)
	}

	options := NewOptionsFromConfigValues(nil, nil, []string{ruleID}, nil, nil, true)
	result, err := ScanWithOptions(root, loaded, options)
	if err != nil {
		t.Fatalf("ScanWithOptions returned error: %v", err)
	}
	return result
}
