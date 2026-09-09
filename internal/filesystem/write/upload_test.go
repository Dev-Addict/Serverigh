package write

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadsFile(t *testing.T) {
	root := t.TempDir()
	service := NewService(root)

	err := service.Upload("/", Upload{
		Name:   "upload.txt",
		Source: strings.NewReader("uploaded"),
	})
	if err != nil {
		t.Fatalf("upload file: %v", err)
	}

	assertFileContent(t, filepath.Join(root, "upload.txt"), "uploaded")
}

func TestUploadsNestedFolderFile(t *testing.T) {
	root := t.TempDir()
	service := NewService(root)

	err := service.Upload("/", Upload{
		Name:   "folder/nested/upload.txt",
		Source: strings.NewReader("uploaded"),
	})
	if err != nil {
		t.Fatalf("upload nested file: %v", err)
	}

	assertFileContent(
		t,
		filepath.Join(root, "folder", "nested", "upload.txt"),
		"uploaded",
	)
}

func TestRejectsUploadTraversal(t *testing.T) {
	service := NewService(t.TempDir())

	err := service.Upload("/", Upload{
		Name:   "../outside.txt",
		Source: strings.NewReader("uploaded"),
	})
	if err == nil {
		t.Fatalf("expected invalid upload path")
	}
}
