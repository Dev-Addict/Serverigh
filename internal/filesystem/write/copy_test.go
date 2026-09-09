package write

import (
	"path/filepath"
	"testing"
)

func TestCopiesFile(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "copies"))
	writeFile(t, filepath.Join(root, "note.txt"), "hello")
	service := NewService(root)

	if err := service.Copy("/note.txt", "/copies"); err != nil {
		t.Fatalf("copy path: %v", err)
	}

	assertFileContent(t, filepath.Join(root, "copies", "note.txt"), "hello")
	assertFileContent(t, filepath.Join(root, "note.txt"), "hello")
}

func TestCopiesFileWithUniqueName(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "hello")
	service := NewService(root)

	if err := service.Copy("/note.txt", "/"); err != nil {
		t.Fatalf("copy path: %v", err)
	}

	assertFileContent(t, filepath.Join(root, "note copy.txt"), "hello")
	assertFileContent(t, filepath.Join(root, "note.txt"), "hello")
}

func TestRejectsCopyingDirectoryIntoItself(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "docs", "drafts"))
	service := NewService(root)

	if err := service.Copy("/docs", "/docs/drafts"); err == nil {
		t.Fatalf("expected copy into child directory to fail")
	}
}
