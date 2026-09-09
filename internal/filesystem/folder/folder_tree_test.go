package folder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBranchShowsAncestorSiblings(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "docs", "drafts"))
	mkdir(t, filepath.Join(root, "docs", "published"))
	mkdir(t, filepath.Join(root, "media"))
	writeFile(t, filepath.Join(root, "note.txt"), "hello")
	service := NewService(root, false)

	tree, err := service.Branch("/docs/drafts")
	if err != nil {
		t.Fatalf("folder branch: %v", err)
	}

	assertTreePath(t, tree.Entries, "/")
	assertTreePath(t, tree.Entries, "/docs")
	assertTreePath(t, tree.Entries, "/docs/drafts")
	assertTreePath(t, tree.Entries, "/docs/published")
	assertTreePath(t, tree.Entries, "/media")
	assertSelected(t, tree.Entries, "/docs/drafts")
	if len(tree.Entries) != 5 {
		t.Fatalf("expected branch and siblings only, got %#v", tree.Entries)
	}
}

func TestChildrenListsOneLevel(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "docs", "drafts", "deep"))
	mkdir(t, filepath.Join(root, "docs", "published"))
	service := NewService(root, false)

	children, _, err := service.Children("/docs")
	if err != nil {
		t.Fatalf("folder children: %v", err)
	}

	assertTreePath(t, children, "/docs/drafts")
	assertTreePath(t, children, "/docs/published")
	if len(children) != 2 {
		t.Fatalf("expected one child level, got %#v", children)
	}
}

func assertTreePath(t *testing.T, entries []Entry, expected string) {
	t.Helper()

	for _, entry := range entries {
		if entry.Path == expected {
			return
		}
	}

	t.Fatalf("expected folder tree to include %q: %#v", expected, entries)
}

func assertSelected(t *testing.T, entries []Entry, expected string) {
	t.Helper()

	for _, entry := range entries {
		if entry.Path == expected && entry.Selected {
			return
		}
	}

	t.Fatalf("expected folder tree to select %q: %#v", expected, entries)
}

func mkdir(t *testing.T, name string) {
	t.Helper()

	if err := os.MkdirAll(name, 0o755); err != nil {
		t.Fatalf("create directory %q: %v", name, err)
	}
}

func writeFile(t *testing.T, name string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatalf("create parent directory %q: %v", name, err)
	}
	if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %q: %v", name, err)
	}
}
