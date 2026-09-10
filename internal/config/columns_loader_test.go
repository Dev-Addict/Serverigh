package config

import "testing"

func TestLoadColumnsUsesDefaultValues(t *testing.T) {
	unsetConfigEnv(t)

	columns, err := LoadColumns(t.TempDir(), ColumnOverrides{})
	if err != nil {
		t.Fatalf("load columns: %v", err)
	}

	if !columns.Size || !columns.Modified || !columns.Created || !columns.Mode {
		t.Fatalf("expected all columns visible by default, got %#v", columns)
	}
}

func TestLoadColumnsUsesConfiguredPrecedence(t *testing.T) {
	unsetConfigEnv(t)
	root := t.TempDir()
	globalRoot := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", globalRoot)
	t.Setenv("SERVERIGH_COLUMN_SIZE", "true")
	t.Setenv("SERVERIGH_COLUMN_MODIFIED", "false")

	writeConfig(t, globalRoot, "[column]\nsize = false\ncreated = false\n")
	writeConfig(t, root, "[column]\nsize = false\nmode = false\n")
	flagMode := true

	columns, err := LoadColumns(root, ColumnOverrides{
		Mode: &flagMode,
	})
	if err != nil {
		t.Fatalf("load columns: %v", err)
	}

	if columns.Size {
		t.Fatalf("expected local config to override size env")
	}
	if columns.Modified {
		t.Fatalf("expected env to override modified default")
	}
	if columns.Created {
		t.Fatalf("expected global config to override created default")
	}
	if !columns.Mode {
		t.Fatalf("expected flag to override local config")
	}
}
