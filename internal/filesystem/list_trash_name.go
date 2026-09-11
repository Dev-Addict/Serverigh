package filesystem

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
)

type trashEntryMetadata struct {
	OriginalPath string `json:"original_path"`
}

func (s Service) trashEntryName(name string) string {
	content, err := os.ReadFile(
		filepath.Join(
			s.root,
			StateDirectoryName,
			trashMetaDirectoryName,
			name+".json",
		),
	)
	if err != nil {
		return name
	}

	var metadata trashEntryMetadata
	if err := json.Unmarshal(content, &metadata); err != nil {
		return name
	}
	if metadata.OriginalPath == "" {
		return name
	}

	return path.Base(metadata.OriginalPath)
}
