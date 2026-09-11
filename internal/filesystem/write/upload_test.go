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

func TestUploadsFileWithIncrementingName(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "upload.txt"), "existing")
	writeFile(t, filepath.Join(root, "upload 1.txt"), "existing numbered")
	service := NewService(root)

	err := service.Upload("/", Upload{
		Name:   "upload.txt",
		Source: strings.NewReader("uploaded"),
	})
	if err != nil {
		t.Fatalf("upload file: %v", err)
	}

	assertFileContent(t, filepath.Join(root, "upload 2.txt"), "uploaded")
	assertFileContent(t, filepath.Join(root, "upload.txt"), "existing")
}

func TestUploadsNumberedFileWithIncrementingName(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "upload 1.txt"), "existing")
	service := NewService(root)

	err := service.Upload("/", Upload{
		Name:   "upload 1.txt",
		Source: strings.NewReader("uploaded"),
	})
	if err != nil {
		t.Fatalf("upload file: %v", err)
	}

	assertFileContent(t, filepath.Join(root, "upload 2.txt"), "uploaded")
	assertFileContent(t, filepath.Join(root, "upload 1.txt"), "existing")
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

func TestUniqueUploadPathsRenamesExistingFolderRoot(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "folder"))
	service := NewService(root)

	paths, err := service.UniqueUploadPaths("/", []string{
		"folder/one.txt",
		"folder/two.txt",
	})
	if err != nil {
		t.Fatalf("unique upload paths: %v", err)
	}

	assertEqualStrings(t, paths, []string{
		"folder 1/one.txt",
		"folder 1/two.txt",
	})
}

func TestUniqueUploadPathsReservesDuplicateRootFiles(t *testing.T) {
	root := t.TempDir()
	service := NewService(root)

	paths, err := service.UniqueUploadPaths("/", []string{
		"upload.txt",
		"upload.txt",
	})
	if err != nil {
		t.Fatalf("unique upload paths: %v", err)
	}

	assertEqualStrings(t, paths, []string{
		"upload.txt",
		"upload 1.txt",
	})
}

func TestUniqueUploadPathsReservesRenamedFolderRoots(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "folder"))
	service := NewService(root)

	paths, err := service.UniqueUploadPaths("/", []string{
		"folder/one.txt",
		"folder 1/two.txt",
	})
	if err != nil {
		t.Fatalf("unique upload paths: %v", err)
	}

	assertEqualStrings(t, paths, []string{
		"folder 1/one.txt",
		"folder 2/two.txt",
	})
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

func assertEqualStrings(t *testing.T, actual []string, expected []string) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("expected %d paths, got %d: %v", len(expected), len(actual), actual)
	}
	for index := range expected {
		if actual[index] != expected[index] {
			t.Fatalf("expected paths %v, got %v", expected, actual)
		}
	}
}
