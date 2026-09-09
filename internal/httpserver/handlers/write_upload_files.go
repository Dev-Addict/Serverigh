package handlers

import (
	"mime/multipart"

	"serverigh/internal/filesystem"
)

func (h Handlers) uploadFile(
	parentPath string,
	relativePaths []string,
	index int,
	fileHeader *multipart.FileHeader,
) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	err = h.files.Upload(parentPath, filesystem.WriteUpload{
		Name:   uploadName(fileHeader.Filename, relativePaths, index),
		Source: file,
	})
	closeErr := file.Close()
	if err != nil {
		return err
	}

	return closeErr
}
