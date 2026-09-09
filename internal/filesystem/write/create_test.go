package write

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCreatesDirectory(t *testing.T) {
	service := NewService(t.TempDir())

	if err := service.CreateDirectory("/", "docs"); err != nil {
		t.Fatalf("create directory: %v", err)
	}

	assertDirExists(t, filepath.Join(service.root, "docs"))
}

func TestCreatesFile(t *testing.T) {
	root := t.TempDir()
	service := NewService(root)

	if err := service.CreateFile("/", "note.txt"); err != nil {
		t.Fatalf("create file: %v", err)
	}

	assertFileContent(t, filepath.Join(root, "note.txt"), "")
}

func TestRejectsNestedEntryName(t *testing.T) {
	service := NewService(t.TempDir())

	err := service.CreateDirectory("/", "../outside")
	if err == nil {
		t.Fatalf("expected invalid path error")
	}

	if !strings.Contains(err.Error(), "invalid path") {
		t.Fatalf("expected invalid path error, got %v", err)
	}
}

func TestRenamesPath(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "hello")
	service := NewService(root)

	if err := service.Rename("/note.txt", "renamed.txt"); err != nil {
		t.Fatalf("rename path: %v", err)
	}

	assertFileContent(t, filepath.Join(root, "renamed.txt"), "hello")
	assertMissing(t, filepath.Join(root, "note.txt"))
}
