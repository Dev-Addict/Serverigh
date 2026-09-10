package config

import (
	"strconv"
	"strings"
)

func appendSettingsLines(lines []string, settings Settings) []string {
	if len(lines) > 0 {
		lines = append(lines, "")
	}

	lines = append(lines,
		"theme = "+quoteTOMLString(settings.Theme),
		"max_preview_bytes = "+strconv.FormatInt(settings.MaxPreviewBytes, 10),
		"",
		"[column]",
		"size = "+strconv.FormatBool(settings.Columns.Size),
		"modified = "+strconv.FormatBool(settings.Columns.Modified),
		"created = "+strconv.FormatBool(settings.Columns.Created),
		"mode = "+strconv.FormatBool(settings.Columns.Mode),
	)

	return lines
}

func managedSettingLine(section string, key string) bool {
	if section == "" {
		return topLevelSettings[key]
	}
	if section == "column" {
		return columnSettings[key]
	}

	return false
}

func quoteTOMLString(value string) string {
	return strconv.Quote(value)
}

func joinConfigLines(lines []string) string {
	return strings.Join(lines, "\n") + "\n"
}
