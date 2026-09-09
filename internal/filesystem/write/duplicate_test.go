package write

import (
	"path/filepath"
	"testing"
)

func TestDuplicatesFileWithIncrementingName(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "hello")
	writeFile(t, filepath.Join(root, "note 1.txt"), "first")
	service := NewService(root)

	if err := service.Duplicate("/note.txt"); err != nil {
		t.Fatalf("duplicate path: %v", err)
	}

	assertFileContent(t, filepath.Join(root, "note 2.txt"), "hello")
	assertFileContent(t, filepath.Join(root, "note.txt"), "hello")
}
