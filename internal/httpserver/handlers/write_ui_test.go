package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/config"
)

func TestFilesRendersWriteMenusInWriteMode(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root, func(cfg *config.Config) {
		cfg.Write = true
	})

	app := fiber.New()
	app.Get("/partials/files", h.Files)

	resp, body := testRequest(t, app, http.MethodGet, "/partials/files?path=/")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	for _, part := range writeMenuParts() {
		if !strings.Contains(body, part) {
			t.Fatalf("expected %q in write-mode response, got %q", part, body)
		}
	}
}

func writeMenuParts() []string {
	return []string{
		`data-write-menu`,
		`class="write-menu-plus" aria-hidden="true">+</span>`,
		`data-write-context-menu`,
		`data-write-create="folder"`,
		`data-write-create="file"`,
		`data-write-upload`,
		`data-context-write="rename"`,
		`data-context-write="duplicate"`,
		`data-context-write="delete"`,
		`data-context-copy="relative"`,
		`webkitdirectory`,
	}
}
