package write

import (
	"io"
	"os"
)

type Upload struct {
	Name   string
	Source io.Reader
}

func (s Service) Upload(parentPath string, upload Upload) error {
	target, err := s.uploadTarget(parentPath, upload.Name)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return wrapWriteError("create uploaded file", err)
	}

	if _, err := io.Copy(file, upload.Source); err != nil {
		_ = file.Close()

		return wrapWriteError("write uploaded file", err)
	}
	if err := file.Close(); err != nil {
		return wrapWriteError("close uploaded file", err)
	}

	return nil
}
