package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveLocalSettingsWritesSettingsConfig(t *testing.T) {
	root := t.TempDir()

	err := SaveLocalSettings(root, Settings{
		Theme:           "dark",
		MaxPreviewBytes: 2048,
		Columns: Columns{
			Size:     true,
			Modified: false,
			Created:  true,
			Mode:     false,
		},
	})
	if err != nil {
		t.Fatalf("save local settings: %v", err)
	}

	content := readConfig(t, root)
	expected := []string{
		`theme = "dark"`,
		`max_preview_bytes = 2048`,
		`[column]`,
		`size = true`,
		`modified = false`,
		`created = true`,
		`mode = false`,
	}
	for _, value := range expected {
		if !strings.Contains(content, value) {
			t.Fatalf("expected %q in config:\n%s", value, content)
		}
	}
}

func TestSaveLocalSettingsPreservesStartupConfig(t *testing.T) {
	root := t.TempDir()
	writeConfig(
		t,
		root,
		"root = \"/workspace\"\nwrite = true\n\n[column]\nsize = false\n",
	)

	err := SaveLocalSettings(root, Settings{
		Theme:           "light",
		MaxPreviewBytes: 4096,
		Columns:         DefaultColumns(),
	})
	if err != nil {
		t.Fatalf("save local settings: %v", err)
	}

	content := readConfig(t, root)
	if !strings.Contains(content, `root = "/workspace"`) {
		t.Fatalf("expected root to be preserved:\n%s", content)
	}
	if !strings.Contains(content, `write = true`) {
		t.Fatalf("expected write to be preserved:\n%s", content)
	}
	if strings.Count(content, "[column]") != 1 {
		t.Fatalf("expected one column section:\n%s", content)
	}
}

func readConfig(t *testing.T, root string) string {
	t.Helper()

	content, err := os.ReadFile(
		filepath.Join(root, ".serverigh", "config.toml"),
	)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	return string(content)
}
