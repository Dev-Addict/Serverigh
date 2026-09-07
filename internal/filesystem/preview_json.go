package filesystem

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

func prettyJSON(content string) (string, string) {
	var value any
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return content, err.Error()
	}

	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return content, "multiple JSON values"
		}

		return content, err.Error()
	}

	var formatted bytes.Buffer
	encoder := json.NewEncoder(&formatted)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return content, err.Error()
	}

	return strings.TrimSuffix(formatted.String(), "\n"), ""
}
