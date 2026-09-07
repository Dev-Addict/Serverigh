package previewer

import (
	"path/filepath"
	"strings"
)

func mediaPreviewKind(name string, mimeType string) PreviewKind {
	if isActiveMediaExtension(name) {
		return ""
	}

	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	mimeType, _, _ = strings.Cut(mimeType, ";")

	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return PreviewKindImage
	case mimeType == "application/pdf":
		return PreviewKindPDF
	case strings.HasPrefix(mimeType, "audio/"):
		return PreviewKindAudio
	case strings.HasPrefix(mimeType, "video/"):
		return PreviewKindVideo
	}

	switch strings.ToLower(filepath.Ext(name)) {
	case ".avif", ".gif", ".jpeg", ".jpg", ".png", ".webp":
		return PreviewKindImage
	case ".pdf":
		return PreviewKindPDF
	case ".aac", ".flac", ".m4a", ".mp3", ".oga", ".ogg", ".opus", ".wav":
		return PreviewKindAudio
	case ".m4v", ".mov", ".mp4", ".ogv", ".webm":
		return PreviewKindVideo
	default:
		return ""
	}
}

func isActiveMediaExtension(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".svg", ".xml", ".xhtml", ".html", ".htm":
		return true
	default:
		return false
	}
}
