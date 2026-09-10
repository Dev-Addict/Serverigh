package config

import (
	"net"
	"os"
	"path/filepath"
	"strconv"

	"serverigh/internal/apperror"
)

const (
	defaultHost            = "127.0.0.1"
	defaultPort            = 4173
	defaultTheme           = "light"
	defaultWrite           = false
	defaultShowHidden      = false
	defaultMaxPreviewBytes = int64(1_048_576)
)

type Config struct {
	Root            string
	Host            string
	Port            int
	Theme           string
	Write           bool
	ShowHidden      bool
	Columns         Columns
	MaxPreviewBytes int64
}

func Default() (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Config{}, apperror.WrapOperation(
			apperror.CodeFilesystem,
			"filesystem error",
			"get working directory",
			err,
		)
	}

	return Config{
		Root:            wd,
		Host:            defaultHost,
		Port:            defaultPort,
		Theme:           defaultTheme,
		Write:           defaultWrite,
		ShowHidden:      defaultShowHidden,
		Columns:         DefaultColumns(),
		MaxPreviewBytes: defaultMaxPreviewBytes,
	}, nil
}

func (c *Config) Normalize() error {
	if c.Host == "" {
		return apperror.New(
			apperror.CodeInvalidConfig,
			"host cannot be empty",
		)
	}

	if c.Port < 1 || c.Port > 65535 {
		return apperror.New(
			apperror.CodeInvalidConfig,
			"port must be between 1 and 65535",
		)
	}

	if c.MaxPreviewBytes < 1 {
		return apperror.New(
			apperror.CodeInvalidConfig,
			"max preview bytes must be greater than zero",
		)
	}

	if c.Theme == "" {
		c.Theme = defaultTheme
	}

	if c.Theme != "light" && c.Theme != "dark" {
		return wrapConfigOperation("validate theme", errInvalidTheme)
	}

	root, err := filepath.Abs(c.Root)
	if err != nil {
		return apperror.WrapOperation(
			apperror.CodeInvalidPath,
			"invalid root path",
			"resolve root path",
			err,
		)
	}

	info, err := os.Stat(root)
	if err != nil {
		return wrapRootOperation("inspect root path", err)
	}

	if !info.IsDir() {
		return apperror.New(
			apperror.CodeNotDirectory,
			"root path is not a directory",
		)
	}

	c.Root = root

	return nil
}

func (c Config) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}
