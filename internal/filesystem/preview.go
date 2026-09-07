package filesystem

import "serverigh/internal/filesystem/previewer"

type PreviewKind = previewer.PreviewKind
type Preview = previewer.Preview

const (
	PreviewKindText     = previewer.PreviewKindText
	PreviewKindCode     = previewer.PreviewKindCode
	PreviewKindMarkdown = previewer.PreviewKindMarkdown
	PreviewKindCSV      = previewer.PreviewKindCSV
	PreviewKindJSON     = previewer.PreviewKindJSON
	PreviewKindImage    = previewer.PreviewKindImage
	PreviewKindPDF      = previewer.PreviewKindPDF
	PreviewKindAudio    = previewer.PreviewKindAudio
	PreviewKindVideo    = previewer.PreviewKindVideo
	PreviewKindBinary   = previewer.PreviewKindBinary
)

func (s Service) Preview(requestPath string) (Preview, error) {
	file, err := s.Open(requestPath)
	if err != nil {
		return Preview{}, err
	}
	defer file.Close()

	preview, err := previewer.Build(
		previewer.File(file.File),
		file.Handle,
		s.maxPreviewBytes,
	)
	if err != nil {
		return Preview{}, wrapFileError("read preview file", err)
	}

	return preview, nil
}
