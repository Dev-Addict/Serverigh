package write

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTrashesPath(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "hello")
	service := NewService(root)

	if err := service.Trash("/note.txt"); err != nil {
		t.Fatalf("trash path: %v", err)
	}

	assertMissing(t, filepath.Join(root, "note.txt"))
	entries, err := os.ReadDir(filepath.Join(root, trashFolderName))
	if err != nil {
		t.Fatalf("read trash: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one trash entry, got %d", len(entries))
	}
}
