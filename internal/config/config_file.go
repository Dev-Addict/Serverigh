package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func loadConfigFile(filename string) (Overrides, error) {
	file, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		return Overrides{}, nil
	}
	if err != nil {
		return Overrides{}, wrapConfigOperation("open config file", err)
	}
	defer file.Close()

	return parseConfigFile(file)
}

func parseConfigFile(file *os.File) (Overrides, error) {
	var overrides Overrides
	section := ""
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(stripComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(line[1 : len(line)-1])

			continue
		}
		if err := applyConfigLine(&overrides, section, line); err != nil {
			return Overrides{}, err
		}
	}
	if err := scanner.Err(); err != nil {
		return Overrides{}, wrapConfigOperation("read config file", err)
	}

	return overrides, nil
}

func stripComment(line string) string {
	if index := strings.Index(line, "#"); index >= 0 {
		return line[:index]
	}

	return line
}
