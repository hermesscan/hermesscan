package files

import "testing"

func TestClassify(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"build.ps1", "powershell"},
		{"ci.sh", "bash"},
		{".github/workflows/ci.yml", "yaml"},
		{"Dockerfile", "docker"},
		{"docker-compose.yml", "docker"},
		{"Makefile", "makefile"},
		{"README.md", ""},
	}

	for _, tt := range tests {
		got := Classify(tt.path)
		if got != tt.want {
			t.Fatalf("Classify(%q) = %q; want %q", tt.path, got, tt.want)
		}
	}
}
