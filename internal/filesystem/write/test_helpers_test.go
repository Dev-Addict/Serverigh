package write

import (
	"os"
	"path/filepath"
	"testing"
)

func assertDirExists(t *testing.T, filename string) {
	t.Helper()

	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("stat directory: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected directory at %q", filename)
	}
}

func assertFileContent(t *testing.T, filename string, expected string) {
	t.Helper()

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("read file %q: %v", filename, err)
	}
	if string(content) != expected {
		t.Fatalf("expected %q content %q, got %q", filename, expected, content)
	}
}

func assertMissing(t *testing.T, filename string) {
	t.Helper()

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		t.Fatalf("expected %q to be missing, got %v", filename, err)
	}
}

func mkdir(t *testing.T, name string) {
	t.Helper()

	if err := os.MkdirAll(name, 0o755); err != nil {
		t.Fatalf("create directory %q: %v", name, err)
	}
}

func writeFile(t *testing.T, name string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatalf("create parent directory %q: %v", name, err)
	}
	if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %q: %v", name, err)
	}
}
