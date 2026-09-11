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

func TestRestoresTrashedPath(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs", "note.txt"), "hello")
	service := NewService(root)

	if err := service.Trash("/docs/note.txt"); err != nil {
		t.Fatalf("trash path: %v", err)
	}
	entry := onlyTrashEntry(t, root)

	if err := service.RestoreTrash("/trash/" + entry.Name()); err != nil {
		t.Fatalf("restore trash: %v", err)
	}

	assertFileContent(t, filepath.Join(root, "docs", "note.txt"), "hello")
	assertMissing(t, filepath.Join(root, stateFolderName, trashFolderName, entry.Name()))
}

func TestDeletesTrashPermanently(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "hello")
	service := NewService(root)

	if err := service.Trash("/note.txt"); err != nil {
		t.Fatalf("trash path: %v", err)
	}
	entry := onlyTrashEntry(t, root)

	if err := service.DeleteTrash("/trash/" + entry.Name()); err != nil {
		t.Fatalf("delete trash: %v", err)
	}

	assertMissing(t, filepath.Join(root, "note.txt"))
	assertMissing(t, filepath.Join(root, stateFolderName, trashFolderName, entry.Name()))
}

func TestRejectsTrashingStateDirectory(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, stateFolderName, trashFolderName))
	service := NewService(root)

	if err := service.Trash("/" + stateFolderName); err == nil {
		t.Fatalf("expected state directory trash to fail")
	}
}

func onlyTrashEntry(t *testing.T, root string) os.DirEntry {
	t.Helper()

	entries, err := os.ReadDir(
		filepath.Join(root, stateFolderName, trashFolderName),
	)
	if err != nil {
		t.Fatalf("read trash: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one trash entry, got %d", len(entries))
	}

	return entries[0]
}
