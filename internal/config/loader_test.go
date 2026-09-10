package config

import "testing"

func TestLoadUsesConfiguredPrecedence(t *testing.T) {
	unsetConfigEnv(t)
	root := t.TempDir()
	globalRoot := t.TempDir()
	localRoot := t.TempDir()
	flagHost := "127.0.0.2"
	flagMaxPreviewBytes := int64(4096)

	t.Setenv("XDG_CONFIG_HOME", globalRoot)
	t.Setenv("SERVERIGH_ROOT", root)
	t.Setenv("SERVERIGH_HOST", "0.0.0.0")
	t.Setenv("SERVERIGH_THEME", "dark")
	t.Setenv("SERVERIGH_WRITE", "true")
	t.Setenv("SERVERIGH_MAX_PREVIEW_BYTES", "2048")

	writeConfig(t, globalRoot, "port = 9000\nshow_hidden = true\n")
	writeConfig(t, root, "root = \""+localRoot+"\"\ntheme = \"light\"\n")

	cfg, err := Default()
	if err != nil {
		t.Fatalf("default config: %v", err)
	}

	loaded, err := Load(cfg, Overrides{
		Host:            &flagHost,
		MaxPreviewBytes: &flagMaxPreviewBytes,
	})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	assertLoadedConfig(t, loaded, localRoot, flagHost, flagMaxPreviewBytes)
}

func assertLoadedConfig(
	t *testing.T,
	loaded Config,
	root string,
	host string,
	maxPreviewBytes int64,
) {
	t.Helper()

	if loaded.Root != root {
		t.Fatalf("expected local root %q, got %q", root, loaded.Root)
	}
	if loaded.Host != host {
		t.Fatalf("expected flag host %q, got %q", host, loaded.Host)
	}
	if loaded.Port != 9000 {
		t.Fatalf("expected global port 9000, got %d", loaded.Port)
	}
	if loaded.Theme != "light" {
		t.Fatalf("expected local theme light, got %q", loaded.Theme)
	}
	if !loaded.Write {
		t.Fatalf("expected env write mode")
	}
	if !loaded.ShowHidden {
		t.Fatalf("expected global show hidden")
	}
	if loaded.MaxPreviewBytes != maxPreviewBytes {
		t.Fatalf(
			"expected flag max preview bytes %d, got %d",
			maxPreviewBytes,
			loaded.MaxPreviewBytes,
		)
	}
}
