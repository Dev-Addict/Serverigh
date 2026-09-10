package config

import (
	"os"
	"path/filepath"
	"testing"
)

func unsetConfigEnv(t *testing.T) {
	t.Helper()

	names := []string{
		"SERVERIGH_ROOT",
		"SERVERIGH_HOST",
		"SERVERIGH_PORT",
		"SERVERIGH_THEME",
		"SERVERIGH_WRITE",
		"SERVERIGH_SHOW_HIDDEN",
		"SERVERIGH_MAX_PREVIEW_BYTES",
		"SERVERIGH_COLUMN_SIZE",
		"SERVERIGH_COLUMN_MODIFIED",
		"SERVERIGH_COLUMN_CREATED",
		"SERVERIGH_COLUMN_MODE",
		"XDG_CONFIG_HOME",
	}
	for _, name := range names {
		unsetEnv(t, name)
	}
}

func unsetEnv(t *testing.T, name string) {
	t.Helper()

	value, ok := os.LookupEnv(name)
	if err := os.Unsetenv(name); err != nil {
		t.Fatalf("unset %s: %v", name, err)
	}
	t.Cleanup(func() {
		if ok {
			_ = os.Setenv(name, value)
		} else {
			_ = os.Unsetenv(name)
		}
	})
}

func writeConfig(t *testing.T, root string, content string) {
	t.Helper()

	filename := filepath.Join(root, ".serverigh", "config.toml")
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatalf("create config directory: %v", err)
	}
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}
