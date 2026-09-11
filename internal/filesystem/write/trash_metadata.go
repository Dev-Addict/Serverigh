package write

import (
	"encoding/json"
	"os"
	"time"
)

type trashMetadata struct {
	OriginalPath string `json:"original_path"`
	TrashedAt    string `json:"trashed_at"`
}

func (s Service) writeTrashMetadata(name string, originalPath string) error {
	metadata := trashMetadata{
		OriginalPath: originalPath,
		TrashedAt:    time.Now().UTC().Format(time.RFC3339Nano),
	}
	content, err := json.Marshal(metadata)
	if err != nil {
		return wrapWriteError("encode trash metadata", err)
	}
	if err := os.WriteFile(trashMetaPath(s.root, name), content, 0o600); err != nil {
		return wrapWriteError("write trash metadata", err)
	}

	return nil
}

func (s Service) readTrashMetadata(name string) (trashMetadata, error) {
	content, err := os.ReadFile(trashMetaPath(s.root, name))
	if err != nil {
		return trashMetadata{}, wrapWriteError("read trash metadata", err)
	}

	var metadata trashMetadata
	if err := json.Unmarshal(content, &metadata); err != nil {
		return trashMetadata{}, ErrInvalidPath
	}

	return metadata, nil
}

func (s Service) removeTrashMetadata(name string) error {
	err := os.Remove(trashMetaPath(s.root, name))
	if err == nil || os.IsNotExist(err) {
		return nil
	}

	return wrapWriteError("remove trash metadata", err)
}
