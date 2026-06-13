package main

import (
	"testing"

	"serverigh/internal/config"
)

func TestCLIParsesServerConfig(t *testing.T) {
	root := t.TempDir()
	cfg := config.Config{
		Root:            root,
		Host:            "127.0.0.1",
		Port:            4173,
		MaxPreviewBytes: 1024,
	}

	var captured config.Config
	app := newCLI(&cfg, func(cfg config.Config) error {
		captured = cfg

		return nil
	})

	err := app.Run([]string{
		"serverigh",
		"--root",
		root,
		"--host",
		"0.0.0.0",
		"--port",
		"9999",
		"--write",
		"--show-hidden",
		"--max-preview-bytes",
		"2048",
	})
	if err != nil {
		t.Fatalf("run cli: %v", err)
	}

	if captured.Root != root {
		t.Fatalf("expected root %q, got %q", root, captured.Root)
	}

	if captured.Host != "0.0.0.0" {
		t.Fatalf("expected host 0.0.0.0, got %q", captured.Host)
	}

	if captured.Port != 9999 {
		t.Fatalf("expected port 9999, got %d", captured.Port)
	}

	if !captured.Write {
		t.Fatalf("expected write mode to be enabled")
	}

	if !captured.ShowHidden {
		t.Fatalf("expected hidden files to be shown")
	}

	if captured.MaxPreviewBytes != 2048 {
		t.Fatalf("expected max preview bytes 2048, got %d", captured.MaxPreviewBytes)
	}
}
