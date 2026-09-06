package filesystem

import (
	"bytes"
	"io"
)

type Preview struct {
	File
	Content   string
	BytesRead int64
	Truncated bool
	IsBinary  bool
}

func (s Service) Preview(requestPath string) (Preview, error) {
	file, err := s.Open(requestPath)
	if err != nil {
		return Preview{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file.Handle, s.maxPreviewBytes+1))
	if err != nil {
		return Preview{}, wrapFileError("read preview file", err)
	}

	truncated := int64(len(data)) > s.maxPreviewBytes
	if truncated {
		data = data[:s.maxPreviewBytes]
	}

	isBinary := bytes.IndexByte(data, 0) >= 0
	content := ""
	if !isBinary {
		content = string(bytes.ToValidUTF8(data, []byte("?")))
	}

	return Preview{
		File:      file.File,
		Content:   content,
		BytesRead: int64(len(data)),
		Truncated: truncated,
		IsBinary:  isBinary,
	}, nil
}
