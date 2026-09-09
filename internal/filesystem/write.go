package filesystem

import (
	"io"

	writefs "serverigh/internal/filesystem/write"
)

type WriteUpload struct {
	Name   string
	Source io.Reader
}

func (s Service) CreateDirectory(parentPath string, name string) error {
	return s.invalidateAfterWrite(
		s.writeService().CreateDirectory(parentPath, name),
	)
}

func (s Service) CreateFile(parentPath string, name string) error {
	return s.invalidateAfterWrite(s.writeService().CreateFile(parentPath, name))
}

func (s Service) Rename(requestPath string, name string) error {
	return s.invalidateAfterWrite(s.writeService().Rename(requestPath, name))
}

func (s Service) Move(requestPath string, targetParentPath string) error {
	return s.invalidateAfterWrite(
		s.writeService().Move(requestPath, targetParentPath),
	)
}

func (s Service) Copy(requestPath string, targetParentPath string) error {
	return s.invalidateAfterWrite(
		s.writeService().Copy(requestPath, targetParentPath),
	)
}

func (s Service) Upload(parentPath string, upload WriteUpload) error {
	err := s.writeService().Upload(parentPath, writefs.Upload{
		Name:   upload.Name,
		Source: upload.Source,
	})

	return s.invalidateAfterWrite(err)
}

func (s Service) Duplicate(requestPath string) error {
	return s.invalidateAfterWrite(s.writeService().Duplicate(requestPath))
}

func (s Service) Trash(requestPath string) error {
	return s.invalidateAfterWrite(s.writeService().Trash(requestPath))
}

func (s Service) writeService() writefs.Service {
	return writefs.NewService(s.root)
}

func (s Service) invalidateAfterWrite(err error) error {
	if err != nil {
		return err
	}

	s.invalidateSearch()

	return nil
}
