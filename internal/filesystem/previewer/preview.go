package previewer

import (
	"bytes"
	"html/template"
	"io"
	"time"
)

type PreviewKind string

const (
	PreviewKindText     PreviewKind = "text"
	PreviewKindCode     PreviewKind = "code"
	PreviewKindMarkdown PreviewKind = "markdown"
	PreviewKindCSV      PreviewKind = "csv"
	PreviewKindJSON     PreviewKind = "json"
	PreviewKindImage    PreviewKind = "image"
	PreviewKindPDF      PreviewKind = "pdf"
	PreviewKindAudio    PreviewKind = "audio"
	PreviewKindVideo    PreviewKind = "video"
	PreviewKindBinary   PreviewKind = "binary"
)

type Preview struct {
	File
	Kind            PreviewKind
	Language        string
	Content         string
	HTMLContent     template.HTML
	DarkHTMLContent template.HTML
	CSVRows         [][]string
	CSVRowLimit     int
	CSVTruncated    bool
	ParseError      string
	BytesRead       int64
	Truncated       bool
	ShowTruncated   bool
	IsBinary        bool
	IsMedia         bool
	IsUnsupported   bool
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

func Build(file File, reader io.Reader, maxPreviewBytes int64) (Preview, error) {
	if mediaKind := mediaPreviewKind(file.Name, file.MIMEType); mediaKind != "" {
		return Preview{
			File:    file,
			Kind:    mediaKind,
			IsMedia: true,
		}, nil
	}

	data, err := io.ReadAll(io.LimitReader(reader, maxPreviewBytes+1))
	if err != nil {
		return Preview{}, err
	}

	truncated := int64(len(data)) > maxPreviewBytes
	if truncated {
		data = data[:maxPreviewBytes]
	}

	isBinary := bytes.IndexByte(data, 0) >= 0
	content := ""
	if !isBinary {
		content = string(bytes.ToValidUTF8(data, []byte("?")))
	}

	preview := Preview{
		File:          file,
		Content:       content,
		BytesRead:     int64(len(data)),
		Truncated:     truncated,
		ShowTruncated: truncated,
		IsBinary:      isBinary,
	}

	preview.classify()

	return preview, nil
}

func (p *Preview) classify() {
	if p.IsBinary {
		p.Kind = PreviewKindBinary
		p.IsUnsupported = true

		return
	}

	switch previewKindFromName(p.Name) {
	case PreviewKindMarkdown:
		p.Kind = PreviewKindMarkdown
		p.HTMLContent = renderMarkdown([]byte(p.Content))
	case PreviewKindCSV:
		p.Kind = PreviewKindCSV
		p.CSVRows, p.CSVTruncated, p.ParseError = parseCSVPreview(
			p.Name,
			p.Content,
		)
		p.CSVRowLimit = maxCSVPreviewRows
	case PreviewKindJSON:
		p.Kind = PreviewKindJSON
		p.Content, p.ParseError = prettyJSON(p.Content)
		p.setHighlightedContent()
	default:
		if p.setHighlightedContent() {
			p.Kind = PreviewKindCode

			return
		}

		p.Kind = PreviewKindText
	}
}

func (p *Preview) setHighlightedContent() bool {
	html, darkHTML, language, ok := highlightCode(p.Name, p.MIMEType, p.Content)
	if !ok {
		return false
	}

	p.HTMLContent = html
	p.DarkHTMLContent = darkHTML
	p.Language = language

	return true
}
