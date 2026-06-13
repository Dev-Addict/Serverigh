package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultHost            = "127.0.0.1"
	defaultPort            = 4173
	defaultWrite           = false
	defaultShowHidden      = false
	defaultMaxPreviewBytes = int64(1_048_576)
)

type Config struct {
	Root            string
	Host            string
	Port            int
	Write           bool
	ShowHidden      bool
	MaxPreviewBytes int64
}

func Default() (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("get working directory: %w", err)
	}

	return Config{
		Root:            wd,
		Host:            defaultHost,
		Port:            defaultPort,
		Write:           defaultWrite,
		ShowHidden:      defaultShowHidden,
		MaxPreviewBytes: defaultMaxPreviewBytes,
	}, nil
}

func (c *Config) Normalize() error {
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if c.MaxPreviewBytes < 1 {
		return fmt.Errorf("max preview bytes must be greater than zero")
	}

	root, err := filepath.Abs(c.Root)
	if err != nil {
		return fmt.Errorf("resolve root path: %w", err)
	}

	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("inspect root path %q: %w", root, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("root path %q is not a directory", root)
	}

	c.Root = root

	return nil
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
