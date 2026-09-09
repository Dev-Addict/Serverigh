package write

import (
	"path/filepath"
	"testing"
)

func TestMovesPath(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "docs"))
	writeFile(t, filepath.Join(root, "note.txt"), "hello")
	service := NewService(root)

	if err := service.Move("/note.txt", "/docs"); err != nil {
		t.Fatalf("move path: %v", err)
	}

	assertFileContent(t, filepath.Join(root, "docs", "note.txt"), "hello")
	assertMissing(t, filepath.Join(root, "note.txt"))
}

func TestRejectsMovingDirectoryIntoItself(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "docs", "drafts"))
	service := NewService(root)

	if err := service.Move("/docs", "/docs/drafts"); err == nil {
		t.Fatalf("expected move into child directory to fail")
	}
}
