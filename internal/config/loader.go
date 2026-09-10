package config

import (
	"os"
	"path/filepath"
)

const configFilename = "config.toml"

func Load(base Config, flagOverrides Overrides) (Config, error) {
	cfg := base

	globalOverrides, err := loadGlobalOverrides()
	if err != nil {
		return Config{}, err
	}
	cfg.apply(globalOverrides)

	envOverrides, err := loadEnvOverrides()
	if err != nil {
		return Config{}, err
	}
	cfg.apply(envOverrides)
	cfg.applyRoot(flagOverrides)

	localOverrides, err := loadConfigFile(localConfigPath(cfg.Root))
	if err != nil {
		return Config{}, err
	}
	cfg.apply(localOverrides)
	cfg.apply(flagOverrides)

	return cfg, nil
}

func loadGlobalOverrides() (Overrides, error) {
	configRoot, err := os.UserConfigDir()
	if err != nil {
		return Overrides{}, nil
	}

	return loadConfigFile(
		filepath.Join(configRoot, ".serverigh", configFilename),
	)
}

func localConfigPath(root string) string {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return filepath.Join(root, ".serverigh", configFilename)
	}

	return filepath.Join(absRoot, ".serverigh", configFilename)
}
