package config

import (
	"os"
	"path/filepath"
)

var topLevelSettings = map[string]bool{
	"theme":             true,
	"max_preview_bytes": true,
}

var columnSettings = map[string]bool{
	"size":     true,
	"modified": true,
	"created":  true,
	"mode":     true,
}

func SaveLocalSettings(root string, settings Settings) error {
	if err := settings.Validate(); err != nil {
		return err
	}

	filename := localConfigPath(root)
	content, err := mergedSettingsConfig(filename, settings)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return wrapConfigOperation("create config directory", err)
	}
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		return wrapConfigOperation("write config file", err)
	}

	return nil
}

func mergedSettingsConfig(filename string, settings Settings) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return "", wrapConfigOperation("read config file", err)
	}

	lines := preservedConfigLines(string(content))
	lines = appendSettingsLines(lines, settings)

	return joinConfigLines(lines), nil
}
