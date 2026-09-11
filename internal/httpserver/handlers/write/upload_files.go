package write

import (
	"mime/multipart"

	"serverigh/internal/filesystem"
)

func (h Handlers) uploadFile(
	parentPath string,
	uploadPath string,
	fileHeader *multipart.FileHeader,
) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	err = h.Files.Upload(parentPath, filesystem.WriteUpload{
		Name:   uploadPath,
		Source: file,
	})
	closeErr := file.Close()
	if err != nil {
		return err
	}

	return closeErr
}
