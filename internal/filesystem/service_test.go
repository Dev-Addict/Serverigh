package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"serverigh/internal/apperror"
)

func TestListReturnsDirectoryEntries(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "folder"))
	writeFile(t, filepath.Join(root, "alpha.txt"), "alpha")
	writeFile(t, filepath.Join(root, ".hidden"), "hidden")

	files := newTestService(t, root, false, 1024)
	listing, err := files.List("/")
	if err != nil {
		t.Fatalf("list root: %v", err)
	}

	if len(listing.Entries) != 2 {
		t.Fatalf("expected 2 visible entries, got %d", len(listing.Entries))
	}

	if listing.Entries[0].Name != "folder" || !listing.Entries[0].IsDir {
		t.Fatalf("expected folder first, got %#v", listing.Entries[0])
	}

	if listing.Entries[1].Name != "alpha.txt" || listing.Entries[1].IsDir {
		t.Fatalf("expected file second, got %#v", listing.Entries[1])
	}
}

func TestListCanIncludeHiddenFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".hidden"), "hidden")

	files := newTestService(t, root, true, 1024)
	listing, err := files.List("/")
	if err != nil {
		t.Fatalf("list root: %v", err)
	}

	if len(listing.Entries) != 1 {
		t.Fatalf("expected hidden file, got %d entries", len(listing.Entries))
	}

	if !listing.Entries[0].IsHidden {
		t.Fatalf("expected hidden entry flag, got %#v", listing.Entries[0])
	}
}

func TestListMarksTruncatedDirectory(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.txt"), "a")
	writeFile(t, filepath.Join(root, "b.txt"), "b")
	writeFile(t, filepath.Join(root, "c.txt"), "c")

	files := newTestService(t, root, false, 1024)
	listing, err := files.list("/", 2)
	if err != nil {
		t.Fatalf("list root: %v", err)
	}

	if !listing.Truncated {
		t.Fatalf("expected listing to be truncated")
	}

	if listing.EntryLimit != 2 {
		t.Fatalf("expected entry limit 2, got %d", listing.EntryLimit)
	}

	if len(listing.Entries) != 2 {
		t.Fatalf("expected 2 listed entries, got %d", len(listing.Entries))
	}
}

func TestResolveRejectsTraversal(t *testing.T) {
	files := newTestService(t, t.TempDir(), false, 1024)

	_, err := files.Resolve("../outside")
	if !errors.Is(err, ErrOutsideRoot) {
		t.Fatalf("expected outside root error, got %v", err)
	}

	if !apperror.HasCode(err, apperror.CodeOutsideRoot) {
		t.Fatalf("expected outside root operational code, got %v", err)
	}
}

func TestResolveRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	writeFile(t, filepath.Join(outside, "secret.txt"), "secret")

	err := os.Symlink(
		filepath.Join(outside, "secret.txt"),
		filepath.Join(root, "secret-link"),
	)
	if err != nil {
		t.Skipf("symlinks are not available: %v", err)
	}

	files := newTestService(t, root, false, 1024)
	_, err = files.File("/secret-link")
	if !errors.Is(err, ErrOutsideRoot) {
		t.Fatalf("expected outside root error, got %v", err)
	}

	if !apperror.HasCode(err, apperror.CodeOutsideRoot) {
		t.Fatalf("expected outside root operational code, got %v", err)
	}
}

func TestFileMissingReturnsOperationalNotFound(t *testing.T) {
	files := newTestService(t, t.TempDir(), false, 1024)

	_, err := files.File("/missing.txt")
	if !apperror.HasCode(err, apperror.CodeNotFound) {
		t.Fatalf("expected not found operational code, got %v", err)
	}
}

func TestPreviewReadsBoundedContent(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "abcdef")

	files := newTestService(t, root, false, 3)
	preview, err := files.Preview("/note.txt")
	if err != nil {
		t.Fatalf("preview file: %v", err)
	}

	if preview.Content != "abc" {
		t.Fatalf("expected truncated content, got %q", preview.Content)
	}

	if !preview.Truncated {
		t.Fatalf("expected preview to be truncated")
	}
}

func TestFileReturnsWholeFileMetadata(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "hello")

	files := newTestService(t, root, false, 1024)
	file, err := files.File("/note.txt")
	if err != nil {
		t.Fatalf("read file metadata: %v", err)
	}

	if file.Name != "note.txt" {
		t.Fatalf("expected file name, got %q", file.Name)
	}

	if file.Size != 5 {
		t.Fatalf("expected size 5, got %d", file.Size)
	}

	if file.Absolute == "" {
		t.Fatalf("expected absolute path")
	}

	if file.ModifiedTime.IsZero() {
		t.Fatalf("expected modification time")
	}

	if file.ModifiedTimeLabel == "" {
		t.Fatalf("expected modification time label")
	}

	if file.CreatedTimeLabel == "" {
		t.Fatalf("expected creation time label")
	}
}

func TestOpenReturnsReadableValidatedHandle(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "hello")

	files := newTestService(t, root, false, 1024)
	file, err := files.Open("/note.txt")
	if err != nil {
		t.Fatalf("open file: %v", err)
	}
	defer file.Close()

	body := make([]byte, 5)
	if _, err := file.Handle.Read(body); err != nil {
		t.Fatalf("read opened handle: %v", err)
	}

	if string(body) != "hello" {
		t.Fatalf("expected file body, got %q", string(body))
	}
}

func newTestService(
	t *testing.T,
	root string,
	showHidden bool,
	maxPreviewBytes int64,
) Service {
	t.Helper()

	files, err := New(root, showHidden, maxPreviewBytes)
	if err != nil {
		t.Fatalf("create filesystem service: %v", err)
	}

	return files
}

func mkdir(t *testing.T, name string) {
	t.Helper()

	if err := os.MkdirAll(name, 0o755); err != nil {
		t.Fatalf("create directory %q: %v", name, err)
	}
}

func writeFile(t *testing.T, name string, content string) {
	t.Helper()

	if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %q: %v", name, err)
	}
}
