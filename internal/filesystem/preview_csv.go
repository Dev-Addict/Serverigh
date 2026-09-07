package filesystem

import (
	"encoding/csv"
	"io"
	"path/filepath"
	"strings"
)

const (
	maxCSVPreviewRows = 500
)

func parseCSVPreview(name string, content string) ([][]string, bool, string) {
	reader := csv.NewReader(strings.NewReader(content))
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true
	if strings.EqualFold(filepath.Ext(name), ".tsv") {
		reader.Comma = '\t'
	}

	rows := make([][]string, 0, 64)
	for len(rows) < maxCSVPreviewRows {
		record, err := reader.Read()
		if err == io.EOF {
			return rows, false, ""
		}

		if err != nil {
			return rows, false, err.Error()
		}

		row := append([]string(nil), record...)
		rows = append(rows, row)
	}

	_, err := reader.Read()
	if err == io.EOF {
		return rows, false, ""
	}

	if err != nil {
		return rows, false, err.Error()
	}

	return rows, true, ""
}
