package files

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// Candidate represents a file selected for scanning.
type Candidate struct {
	Path string
	Type string
}

var ignoredDirectories = map[string]bool{
	".git":         true,
	".hg":          true,
	".svn":         true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	"bin":          true,
	"obj":          true,
	".terraform":   true,
	".venv":        true,
	"venv":         true,
	"__pycache__":  true,
	".idea":        true,
	".vscode":      false,
}

// Discover walks root and returns candidate build, CI, and script files.
func Discover(root string) ([]Candidate, error) {
	var candidates []Candidate

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			name := entry.Name()
			if ignoredDirectories[name] {
				return filepath.SkipDir
			}
			return nil
		}

		fileType := Classify(path)
		if fileType == "" {
			return nil
		}

		candidates = append(candidates, Candidate{Path: path, Type: fileType})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return candidates, nil
}

// Classify returns a HermesScan file type for files likely relevant to CI/build risk scanning.
func Classify(path string) string {
	base := strings.ToLower(filepath.Base(path))
	ext := strings.ToLower(filepath.Ext(path))
	clean := filepath.ToSlash(strings.ToLower(path))

	switch base {
	case "makefile", "gnumakefile":
		return "makefile"
	case "dockerfile":
		return "docker"
	case "docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml":
		return "docker"
	}

	switch ext {
	case ".ps1", ".psm1", ".psd1":
		return "powershell"
	case ".sh", ".bash", ".zsh":
		return "bash"
	case ".yml", ".yaml":
		return "yaml"
	case ".mk":
		return "makefile"
	}

	if strings.Contains(clean, "/.github/workflows/") && (ext == ".yml" || ext == ".yaml") {
		return "yaml"
	}

	return ""
}
