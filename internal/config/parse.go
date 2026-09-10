package config

import (
	"errors"
	"strconv"
	"strings"
)

func applyConfigLine(overrides *Overrides, section string, line string) error {
	key, value, ok := strings.Cut(line, "=")
	if !ok {
		return wrapConfigOperation(
			"parse config",
			errors.New("missing '='"),
		)
	}

	if section == "column" {
		return applyColumnConfigValue(overrides, key, value)
	}
	if section != "" {
		return nil
	}

	return applyTopLevelConfigValue(overrides, key, value)
}

func applyColumnConfigValue(
	overrides *Overrides,
	key string,
	value string,
) error {
	parsed, err := parseConfigBool(value)
	if err != nil {
		return wrapConfigOperation("parse column "+strings.TrimSpace(key), err)
	}

	switch strings.TrimSpace(key) {
	case "size":
		overrides.Columns.Size = &parsed
	case "modified":
		overrides.Columns.Modified = &parsed
	case "created":
		overrides.Columns.Created = &parsed
	case "mode":
		overrides.Columns.Mode = &parsed
	}

	return nil
}

func applyTopLevelConfigValue(
	overrides *Overrides,
	key string,
	value string,
) error {
	switch strings.TrimSpace(key) {
	case "root":
		parsed := parseConfigString(value)
		overrides.Root = &parsed
	case "host":
		parsed := parseConfigString(value)
		overrides.Host = &parsed
	case "port":
		parsed, err := parseConfigInt(value)
		if err != nil {
			return wrapConfigOperation("parse port", err)
		}
		overrides.Port = &parsed
	case "theme":
		parsed := parseConfigString(value)
		overrides.Theme = &parsed
	case "write":
		parsed, err := parseConfigBool(value)
		if err != nil {
			return wrapConfigOperation("parse write", err)
		}
		overrides.Write = &parsed
	case "show_hidden":
		parsed, err := parseConfigBool(value)
		if err != nil {
			return wrapConfigOperation("parse show_hidden", err)
		}
		overrides.ShowHidden = &parsed
	case "max_preview_bytes":
		parsed, err := parseConfigInt64(value)
		if err != nil {
			return wrapConfigOperation("parse max_preview_bytes", err)
		}
		overrides.MaxPreviewBytes = &parsed
	}

	return nil
}

func parseConfigString(value string) string {
	return strings.Trim(strings.TrimSpace(value), `"`)
}

func parseConfigInt(value string) (int, error) {
	parsed, err := strconv.Atoi(parseConfigString(value))
	if err != nil {
		return 0, err
	}

	return parsed, nil
}

func parseConfigInt64(value string) (int64, error) {
	parsed, err := strconv.ParseInt(parseConfigString(value), 10, 64)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}

func parseConfigBool(value string) (bool, error) {
	trimmed := parseConfigString(value)
	switch strings.ToLower(trimmed) {
	case "yes", "on":
		return true, nil
	case "no", "off":
		return false, nil
	default:
		return strconv.ParseBool(trimmed)
	}
}
