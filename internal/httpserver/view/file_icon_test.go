package view

import "testing"

func TestFileIconUsesFolderIconForDirectories(t *testing.T) {
	icon := FileIcon("docs", "folder", true)

	if icon.Name != "folder" {
		t.Fatalf("expected folder icon, got %q", icon.Name)
	}
}

func TestFileIconUsesSymlinkIcon(t *testing.T) {
	icon := FileIcon("latest", "symlink", false)

	if icon.Name != "symlink" {
		t.Fatalf("expected symlink icon, got %q", icon.Name)
	}
}

func TestFileIconUsesExactFilenameBeforeExtension(t *testing.T) {
	icon := FileIcon("package.json", "file", false)

	if icon.Name != "javascript" {
		t.Fatalf("expected package.json icon, got %q", icon.Name)
	}
}

func TestFileIconUsesCompoundExtension(t *testing.T) {
	icon := FileIcon("server.test.ts", "file", false)

	if icon.Name != "typescript" {
		t.Fatalf("expected compound extension icon, got %q", icon.Name)
	}
}

func TestFileIconUsesLongestCompoundExtension(t *testing.T) {
	icon := FileIcon("release.tar.gz", "file", false)

	if icon.Name != "archive" {
		t.Fatalf("expected archive icon, got %q", icon.Name)
	}
}

func TestFileIconUsesExtension(t *testing.T) {
	icon := FileIcon("photo.webp", "file", false)

	if icon.Name != "image" {
		t.Fatalf("expected image icon, got %q", icon.Name)
	}
}

func TestFileIconUsesExpandedLanguageIcon(t *testing.T) {
	icon := FileIcon("script.py", "file", false)

	if icon.Name != "python" {
		t.Fatalf("expected python icon, got %q", icon.Name)
	}
}

func TestFileIconFallsBackToGenericFile(t *testing.T) {
	icon := FileIcon("unknown.custom", "file", false)

	if icon.Name != "file" {
		t.Fatalf("expected generic icon, got %q", icon.Name)
	}
}
