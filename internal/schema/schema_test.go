package schema

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/hermesscan/hermesscan/internal/config"
)

func TestConfigProfilesMatchSchema(t *testing.T) {
	schema := compileSchema(t, filepath.Join(repoRoot(t), "schemas", "hermesscan-config.schema.json"))

	for _, profile := range config.SupportedProfiles() {
		t.Run(profile, func(t *testing.T) {
			data, err := json.Marshal(config.Profile(profile))
			if err != nil {
				t.Fatalf("marshal config: %v", err)
			}
			var value any
			if err := json.Unmarshal(data, &value); err != nil {
				t.Fatalf("unmarshal config: %v", err)
			}
			if err := schema.Validate(value); err != nil {
				t.Fatalf("validate %s profile: %v", profile, err)
			}
		})
	}
}

func TestRuleCatalogsMatchSchema(t *testing.T) {
	root := repoRoot(t)
	schema := compileSchema(t, filepath.Join(root, "schemas", "hermesscan-rule-catalog.schema.json"))

	for _, path := range []string{
		filepath.Join(root, "rules", "hermes.rules.json"),
		filepath.Join(root, "internal", "rules", "defaults", "hermes.rules.json"),
	} {
		t.Run(filepath.Base(filepath.Dir(path)), func(t *testing.T) {
			value := readJSON(t, path)
			if err := schema.Validate(value); err != nil {
				t.Fatalf("validate %s: %v", path, err)
			}
		})
	}
}

func compileSchema(t *testing.T, path string) *jsonschema.Schema {
	t.Helper()
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(path)
	if err != nil {
		t.Fatalf("compile schema %s: %v", path, err)
	}
	return schema
}

func readJSON(t *testing.T, path string) any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return value
}

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	return root
}
