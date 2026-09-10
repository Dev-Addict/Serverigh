package config

import "strings"

func preservedConfigLines(content string) []string {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	var lines []string
	section := ""
	for _, rawLine := range strings.Split(content, "\n") {
		line := strings.TrimSpace(stripComment(rawLine))
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(line[1 : len(line)-1])
			if section == "column" {
				continue
			}
			lines = append(lines, rawLine)

			continue
		}
		if section == "column" {
			continue
		}
		key, _, ok := strings.Cut(line, "=")
		if ok && managedSettingLine(section, strings.TrimSpace(key)) {
			continue
		}
		lines = append(lines, rawLine)
	}

	return trimTrailingBlankLines(lines)
}

func trimTrailingBlankLines(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}
