package bulk

import (
	"archive/zip"
	"os"
	"testing"
)

func assertZipEntries(t *testing.T, reader *zip.Reader, names ...string) {
	t.Helper()

	entries := map[string]bool{}
	for _, file := range reader.File {
		entries[file.Name] = true
	}
	for _, name := range names {
		if !entries[name] {
			t.Fatalf("expected %q in zip, got %#v", name, entries)
		}
	}
}

func assertMissing(t *testing.T, filename string) {
	t.Helper()

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		t.Fatalf("expected %q to be missing, got %v", filename, err)
	}
}

func assertFileContent(t *testing.T, filename string, expected string) {
	t.Helper()

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("read %q: %v", filename, err)
	}
	if string(content) != expected {
		t.Fatalf("expected %q content %q, got %q", filename, expected, content)
	}
}
