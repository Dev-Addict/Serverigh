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
	entries, err := os.ReadDir(
		filepath.Join(root, stateFolderName, trashFolderName),
	)
	if err != nil {
		t.Fatalf("read trash: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one trash entry, got %d", len(entries))
	}
}

func TestRejectsTrashingStateDirectory(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, stateFolderName, trashFolderName))
	service := NewService(root)

	if err := service.Trash("/" + stateFolderName); err == nil {
		t.Fatalf("expected state directory trash to fail")
	}
}
