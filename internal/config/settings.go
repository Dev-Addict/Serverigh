package config

type Settings struct {
	Theme           string
	MaxPreviewBytes int64
	Columns         Columns
}

func SettingsFromConfig(cfg Config) Settings {
	return Settings{
		Theme:           cfg.Theme,
		MaxPreviewBytes: cfg.MaxPreviewBytes,
		Columns:         cfg.Columns,
	}
}

func (s Settings) Validate() error {
	if s.Theme != "light" && s.Theme != "dark" {
		return wrapConfigOperation("validate theme", errInvalidTheme)
	}
	if s.MaxPreviewBytes < 1 {
		return wrapConfigOperation(
			"validate max preview bytes",
			errInvalidMaxPreviewBytes,
		)
	}

	return nil
}
