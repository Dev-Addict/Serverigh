package fileinfo

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreationTimeReturnsValidTimeWhenKnown(t *testing.T) {
	root := t.TempDir()
	filename := filepath.Join(root, "note.txt")
	if err := os.WriteFile(filename, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	createdTime, ok := CreationTime(filename)
	if !ok {
		return
	}

	if createdTime.IsZero() {
		t.Fatalf("expected non-zero creation time when known")
	}

	if createdTime.After(time.Now().Add(time.Minute)) {
		t.Fatalf("expected creation time not to be in the future")
	}
}

func TestCreationTimeReturnsUnknownForMissingFile(t *testing.T) {
	_, ok := CreationTime(filepath.Join(t.TempDir(), "missing.txt"))
	if ok {
		t.Fatalf("expected missing file creation time to be unknown")
	}
}
