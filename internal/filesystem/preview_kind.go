package filesystem

import (
	"path/filepath"
	"strings"
)

func previewKindFromName(name string) PreviewKind {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".md", ".markdown", ".mdown":
		return PreviewKindMarkdown
	case ".csv", ".tsv":
		return PreviewKindCSV
	case ".json":
		return PreviewKindJSON
	default:
		return PreviewKindText
	}
}
