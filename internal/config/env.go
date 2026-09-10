package config

import (
	"os"
)

func loadEnvOverrides() (Overrides, error) {
	var overrides Overrides

	if err := applyEnvString(&overrides.Root, "SERVERIGH_ROOT"); err != nil {
		return Overrides{}, err
	}
	if err := applyEnvString(&overrides.Host, "SERVERIGH_HOST"); err != nil {
		return Overrides{}, err
	}
	if err := applyEnvInt(&overrides.Port, "SERVERIGH_PORT"); err != nil {
		return Overrides{}, err
	}
	if err := applyEnvString(&overrides.Theme, "SERVERIGH_THEME"); err != nil {
		return Overrides{}, err
	}
	if err := applyEnvBool(&overrides.Write, "SERVERIGH_WRITE"); err != nil {
		return Overrides{}, err
	}
	if err := applyEnvBool(
		&overrides.ShowHidden,
		"SERVERIGH_SHOW_HIDDEN",
	); err != nil {
		return Overrides{}, err
	}
	if err := applyEnvInt64(
		&overrides.MaxPreviewBytes,
		"SERVERIGH_MAX_PREVIEW_BYTES",
	); err != nil {
		return Overrides{}, err
	}
	if err := applyEnvBool(
		&overrides.Columns.Size,
		"SERVERIGH_COLUMN_SIZE",
	); err != nil {
		return Overrides{}, err
	}
	if err := applyEnvBool(
		&overrides.Columns.Modified,
		"SERVERIGH_COLUMN_MODIFIED",
	); err != nil {
		return Overrides{}, err
	}
	if err := applyEnvBool(
		&overrides.Columns.Created,
		"SERVERIGH_COLUMN_CREATED",
	); err != nil {
		return Overrides{}, err
	}
	if err := applyEnvBool(
		&overrides.Columns.Mode,
		"SERVERIGH_COLUMN_MODE",
	); err != nil {
		return Overrides{}, err
	}

	return overrides, nil
}

func applyEnvString(target **string, name string) error {
	value, ok := os.LookupEnv(name)
	if !ok {
		return nil
	}

	*target = &value

	return nil
}

func applyEnvInt(target **int, name string) error {
	value, ok := os.LookupEnv(name)
	if !ok {
		return nil
	}

	parsed, err := parseConfigInt(value)
	if err != nil {
		return wrapConfigOperation("parse "+name, err)
	}

	*target = &parsed

	return nil
}

func applyEnvInt64(target **int64, name string) error {
	value, ok := os.LookupEnv(name)
	if !ok {
		return nil
	}

	parsed, err := parseConfigInt64(value)
	if err != nil {
		return wrapConfigOperation("parse "+name, err)
	}

	*target = &parsed

	return nil
}

func applyEnvBool(target **bool, name string) error {
	value, ok := os.LookupEnv(name)
	if !ok {
		return nil
	}

	parsed, err := parseConfigBool(value)
	if err != nil {
		return wrapConfigOperation("parse "+name, err)
	}

	*target = &parsed

	return nil
}
