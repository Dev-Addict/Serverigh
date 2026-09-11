package bulk

import (
	"archive/zip"
	"io"
)

func (h Handlers) writeBulkZip(
	writer io.Writer,
	entries []bulkZipEntry,
) error {
	archive := zip.NewWriter(writer)
	for _, entry := range entries {
		if err := h.addFileToZip(archive, entry.Path); err != nil {
			archive.Close()

			return err
		}
	}

	return archive.Close()
}

func (h Handlers) addFileToZip(archive *zip.Writer, requestPath string) error {
	file, err := h.Files.Open(requestPath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Handle.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = archiveName(file.Path)
	header.Method = zip.Deflate

	target, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(target, file.Handle)

	return err
}
