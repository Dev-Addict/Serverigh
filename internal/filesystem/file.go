package filesystem

import (
	"io"
	"net/http"
	"os"
	"time"

	"serverigh/internal/filesystem/fileinfo"
)

type OpenedFile struct {
	File
	Handle *os.File
}

type File struct {
	Name              string
	Path              string
	Absolute          string
	Size              int64
	SizeLabel         string
	Mode              string
	CreatedTime       time.Time
	CreatedTimeLabel  string
	CreatedTimeValue  string
	ModifiedTime      time.Time
	ModifiedTimeLabel string
	ModifiedTimeValue string
	CreationTimeKnown bool
	MIMEType          string
}

func (s Service) File(requestPath string) (File, error) {
	opened, err := s.Open(requestPath)
	if err != nil {
		return File{}, err
	}
	defer opened.Handle.Close()

	return opened.File, nil
}

func (s Service) Open(requestPath string) (OpenedFile, error) {
	resolved, err := s.Resolve(requestPath)
	if err != nil {
		return OpenedFile{}, err
	}

	if resolved.Info.IsDir() {
		return OpenedFile{}, ErrIsDirectory
	}

	handle, err := os.Open(resolved.Absolute)
	if err != nil {
		return OpenedFile{}, wrapFileError("open file", err)
	}

	info, err := handle.Stat()
	if err != nil {
		handle.Close()

		return OpenedFile{}, wrapFileError("inspect opened file", err)
	}

	if info.IsDir() {
		handle.Close()

		return OpenedFile{}, ErrIsDirectory
	}

	if !os.SameFile(resolved.Info, info) {
		handle.Close()

		return OpenedFile{}, ErrPathChanged
	}

	mimeType, err := detectMIMEType(handle)
	if err != nil {
		handle.Close()

		return OpenedFile{}, err
	}

	return OpenedFile{
		File:   fileFromInfo(resolved, info, mimeType),
		Handle: handle,
	}, nil
}

func (f OpenedFile) Close() error {
	return f.Handle.Close()
}

func detectMIMEType(file *os.File) (string, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", wrapFileError("read file for mime detection", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", wrapFileError("rewind file after mime detection", err)
	}

	return http.DetectContentType(buffer[:n]), nil
}

func fileFromInfo(
	resolved ResolvedPath,
	info os.FileInfo,
	mimeType string,
) File {
	createdTime, creationTimeKnown := fileinfo.CreationTime(resolved.Absolute)
	modifiedTime := info.ModTime()

	return File{
		Name:              info.Name(),
		Path:              resolved.Path,
		Absolute:          resolved.Absolute,
		Size:              info.Size(),
		SizeLabel:         fileinfo.FormatSize(info.Size()),
		Mode:              info.Mode().String(),
		CreatedTime:       createdTime,
		CreatedTimeLabel:  fileinfo.FormatOptionalTime(createdTime, creationTimeKnown),
		CreatedTimeValue:  fileinfo.FormatOptionalTimeValue(createdTime, creationTimeKnown),
		ModifiedTime:      modifiedTime,
		ModifiedTimeLabel: fileinfo.FormatTime(modifiedTime),
		ModifiedTimeValue: fileinfo.FormatTimeValue(modifiedTime),
		CreationTimeKnown: creationTimeKnown,
		MIMEType:          mimeType,
	}
}
