package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"serverigh/internal/apperror"
)

func TestListReturnsDirectoryEntries(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "folder"))
	writeFile(t, filepath.Join(root, "alpha.txt"), "alpha")
	writeFile(t, filepath.Join(root, ".hidden"), "hidden")

	files := newTestService(t, root, false, 1024)
	listing, err := files.List("/")
	if err != nil {
		t.Fatalf("list root: %v", err)
	}

	if len(listing.Entries) != 2 {
		t.Fatalf("expected 2 visible entries, got %d", len(listing.Entries))
	}

	if listing.Entries[0].Name != "folder" || !listing.Entries[0].IsDir {
		t.Fatalf("expected folder first, got %#v", listing.Entries[0])
	}

	if listing.Entries[1].Name != "alpha.txt" || listing.Entries[1].IsDir {
		t.Fatalf("expected file second, got %#v", listing.Entries[1])
	}
}

func TestListCanIncludeHiddenFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".hidden"), "hidden")

	files := newTestService(t, root, true, 1024)
	listing, err := files.List("/")
	if err != nil {
		t.Fatalf("list root: %v", err)
	}

	if len(listing.Entries) != 1 {
		t.Fatalf("expected hidden file, got %d entries", len(listing.Entries))
	}

	if !listing.Entries[0].IsHidden {
		t.Fatalf("expected hidden entry flag, got %#v", listing.Entries[0])
	}
}

func TestListMarksTruncatedDirectory(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.txt"), "a")
	writeFile(t, filepath.Join(root, "b.txt"), "b")
	writeFile(t, filepath.Join(root, "c.txt"), "c")

	files := newTestService(t, root, false, 1024)
	listing, err := files.list("/", 2)
	if err != nil {
		t.Fatalf("list root: %v", err)
	}

	if !listing.Truncated {
		t.Fatalf("expected listing to be truncated")
	}

	if listing.EntryLimit != 2 {
		t.Fatalf("expected entry limit 2, got %d", listing.EntryLimit)
	}

	if len(listing.Entries) != 2 {
		t.Fatalf("expected 2 listed entries, got %d", len(listing.Entries))
	}
}

func TestResolveRejectsTraversal(t *testing.T) {
	files := newTestService(t, t.TempDir(), false, 1024)

	_, err := files.Resolve("../outside")
	if !errors.Is(err, ErrOutsideRoot) {
		t.Fatalf("expected outside root error, got %v", err)
	}

	if !apperror.HasCode(err, apperror.CodeOutsideRoot) {
		t.Fatalf("expected outside root operational code, got %v", err)
	}
}

func TestResolveRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	writeFile(t, filepath.Join(outside, "secret.txt"), "secret")

	err := os.Symlink(
		filepath.Join(outside, "secret.txt"),
		filepath.Join(root, "secret-link"),
	)
	if err != nil {
		t.Skipf("symlinks are not available: %v", err)
	}

	files := newTestService(t, root, false, 1024)
	_, err = files.File("/secret-link")
	if !errors.Is(err, ErrOutsideRoot) {
		t.Fatalf("expected outside root error, got %v", err)
	}

	if !apperror.HasCode(err, apperror.CodeOutsideRoot) {
		t.Fatalf("expected outside root operational code, got %v", err)
	}
}

func TestFileMissingReturnsOperationalNotFound(t *testing.T) {
	files := newTestService(t, t.TempDir(), false, 1024)

	_, err := files.File("/missing.txt")
	if !apperror.HasCode(err, apperror.CodeNotFound) {
		t.Fatalf("expected not found operational code, got %v", err)
	}
}

func TestPreviewReadsBoundedContent(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "abcdef")

	files := newTestService(t, root, false, 3)
	preview, err := files.Preview("/note.txt")
	if err != nil {
		t.Fatalf("preview file: %v", err)
	}

	if preview.Content != "abc" {
		t.Fatalf("expected truncated content, got %q", preview.Content)
	}

	if !preview.Truncated {
		t.Fatalf("expected preview to be truncated")
	}
}

func TestPreviewRendersMarkdown(t *testing.T) {
	root := t.TempDir()
	writeFile(
		t,
		filepath.Join(root, "README.md"),
		"# Title\n\n| Name | Value |\n| --- | --- |\n| One | Two |\n\n"+
			"[bad](javascript:alert(1))\n\n<script>x</script>",
	)

	files := newTestService(t, root, false, 1024)
	preview, err := files.Preview("/README.md")
	if err != nil {
		t.Fatalf("preview markdown: %v", err)
	}

	if preview.Kind != PreviewKindMarkdown {
		t.Fatalf("expected markdown preview, got %q", preview.Kind)
	}

	html := string(preview.HTMLContent)
	if !strings.Contains(html, "<h1>Title</h1>") {
		t.Fatalf("expected rendered heading, got %q", html)
	}

	if !strings.Contains(html, "<table>") {
		t.Fatalf("expected rendered table, got %q", html)
	}

	if strings.Contains(html, "<script>") {
		t.Fatalf("expected raw html to be skipped, got %q", html)
	}

	if strings.Contains(html, "javascript:") {
		t.Fatalf("expected unsafe links to be skipped, got %q", html)
	}
}

func TestPreviewParsesCSV(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "data.csv"), "name,value\none,\"two, too\"\n")

	files := newTestService(t, root, false, 1024)
	preview, err := files.Preview("/data.csv")
	if err != nil {
		t.Fatalf("preview csv: %v", err)
	}

	if preview.Kind != PreviewKindCSV {
		t.Fatalf("expected csv preview, got %q", preview.Kind)
	}

	if preview.ParseError != "" {
		t.Fatalf("expected valid csv, got parse error %q", preview.ParseError)
	}

	if got := preview.CSVRows[1][1]; got != "two, too" {
		t.Fatalf("expected quoted csv field, got %q", got)
	}
}

func TestPreviewPrettyPrintsJSON(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "data.json"), `{"name":"one","count":2}`)

	files := newTestService(t, root, false, 1024)
	preview, err := files.Preview("/data.json")
	if err != nil {
		t.Fatalf("preview json: %v", err)
	}

	if preview.Kind != PreviewKindJSON {
		t.Fatalf("expected json preview, got %q", preview.Kind)
	}

	if !strings.Contains(preview.Content, "\n  \"name\": \"one\"") {
		t.Fatalf("expected pretty json, got %q", preview.Content)
	}
}

func TestPreviewFallsBackForInvalidJSON(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "data.json"), `{"name":`)

	files := newTestService(t, root, false, 1024)
	preview, err := files.Preview("/data.json")
	if err != nil {
		t.Fatalf("preview json: %v", err)
	}

	if preview.Kind != PreviewKindJSON {
		t.Fatalf("expected json preview, got %q", preview.Kind)
	}

	if preview.ParseError == "" {
		t.Fatalf("expected parse error")
	}

	if preview.Content != `{"name":` {
		t.Fatalf("expected original content fallback, got %q", preview.Content)
	}
}

func TestPreviewClassifiesMediaBeforeBinary(t *testing.T) {
	root := t.TempDir()
	writeBytes(
		t,
		filepath.Join(root, "image.png"),
		[]byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR"),
	)

	files := newTestService(t, root, false, 1024)
	preview, err := files.Preview("/image.png")
	if err != nil {
		t.Fatalf("preview image: %v", err)
	}

	if preview.Kind != PreviewKindImage {
		t.Fatalf("expected image preview, got %q", preview.Kind)
	}

	if !preview.IsMedia {
		t.Fatalf("expected media preview")
	}

	if preview.BytesRead != 0 {
		t.Fatalf("expected media preview to skip content read, got %d", preview.BytesRead)
	}
}

func TestPreviewUsesMediaExtensionFallbacks(t *testing.T) {
	tests := []struct {
		name string
		kind PreviewKind
	}{
		{name: "document.pdf", kind: PreviewKindPDF},
		{name: "sound.mp3", kind: PreviewKindAudio},
		{name: "movie.mp4", kind: PreviewKindVideo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeBytes(t, filepath.Join(root, tt.name), []byte{0x00, 0x01, 0x02})

			files := newTestService(t, root, false, 1024)
			preview, err := files.Preview("/" + tt.name)
			if err != nil {
				t.Fatalf("preview media: %v", err)
			}

			if preview.Kind != tt.kind {
				t.Fatalf("expected %q preview, got %q", tt.kind, preview.Kind)
			}

			if !preview.IsMedia {
				t.Fatalf("expected media preview")
			}

			if preview.Truncated {
				t.Fatalf("expected media preview to avoid bounded truncation")
			}
		})
	}
}

func TestPreviewDoesNotRenderSVGAsMedia(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "image.svg"), "<svg></svg>")

	files := newTestService(t, root, false, 1024)
	preview, err := files.Preview("/image.svg")
	if err != nil {
		t.Fatalf("preview svg: %v", err)
	}

	if preview.Kind != PreviewKindText {
		t.Fatalf("expected active svg to use text preview, got %q", preview.Kind)
	}
}

func TestPreviewClassifiesUnknownBinary(t *testing.T) {
	root := t.TempDir()
	writeBytes(t, filepath.Join(root, "archive.bin"), []byte{0x00, 0x01, 0x02})

	files := newTestService(t, root, false, 1024)
	preview, err := files.Preview("/archive.bin")
	if err != nil {
		t.Fatalf("preview binary: %v", err)
	}

	if preview.Kind != PreviewKindBinary {
		t.Fatalf("expected binary preview, got %q", preview.Kind)
	}

	if !preview.IsUnsupported {
		t.Fatalf("expected unsupported binary preview")
	}
}

func TestFileReturnsWholeFileMetadata(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "hello")

	files := newTestService(t, root, false, 1024)
	file, err := files.File("/note.txt")
	if err != nil {
		t.Fatalf("read file metadata: %v", err)
	}

	if file.Name != "note.txt" {
		t.Fatalf("expected file name, got %q", file.Name)
	}

	if file.Size != 5 {
		t.Fatalf("expected size 5, got %d", file.Size)
	}

	if file.Absolute == "" {
		t.Fatalf("expected absolute path")
	}

	if file.ModifiedTime.IsZero() {
		t.Fatalf("expected modification time")
	}

	if file.ModifiedTimeLabel == "" {
		t.Fatalf("expected modification time label")
	}

	if file.CreatedTimeLabel == "" {
		t.Fatalf("expected creation time label")
	}
}

func TestOpenReturnsReadableValidatedHandle(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "note.txt"), "hello")

	files := newTestService(t, root, false, 1024)
	file, err := files.Open("/note.txt")
	if err != nil {
		t.Fatalf("open file: %v", err)
	}
	defer file.Close()

	body := make([]byte, 5)
	if _, err := file.Handle.Read(body); err != nil {
		t.Fatalf("read opened handle: %v", err)
	}

	if string(body) != "hello" {
		t.Fatalf("expected file body, got %q", string(body))
	}
}

func newTestService(
	t *testing.T,
	root string,
	showHidden bool,
	maxPreviewBytes int64,
) Service {
	t.Helper()

	files, err := New(root, showHidden, maxPreviewBytes)
	if err != nil {
		t.Fatalf("create filesystem service: %v", err)
	}

	return files
}

func mkdir(t *testing.T, name string) {
	t.Helper()

	if err := os.MkdirAll(name, 0o755); err != nil {
		t.Fatalf("create directory %q: %v", name, err)
	}
}

func writeFile(t *testing.T, name string, content string) {
	t.Helper()

	writeBytes(t, name, []byte(content))
}

func writeBytes(t *testing.T, name string, content []byte) {
	t.Helper()

	if err := os.WriteFile(name, content, 0o644); err != nil {
		t.Fatalf("write file %q: %v", name, err)
	}
}
