package handlers

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/httpserver/view"
)

func (h Handlers) Files(c *fiber.Ctx) error {
	start := time.Now()
	options := listOptionsFromRequest(c)
	listing, err := h.files.ListWithOptions(c.Query("path", "/"), options)
	if err != nil {
		return h.writeOperationalError(c, err)
	}
	logHandlerOperation(
		c,
		slog.LevelInfo,
		"directory listed",
		start,
		"path",
		listing.Path,
		"entry_count",
		len(listing.Entries),
		"truncated",
		listing.Truncated,
	)

	files := view.Files(listing, h.config.Write)

	if c.Get("HX-Request") == "true" {
		c.Set("HX-Push-Url", files.Controls.BrowseURL)

		return h.render(c, "files_update.html", view.FilesUpdate{
			Listing: files,
			PathSummary: view.PathSummaryView{
				Root: h.config.Root,
				Path: listing.Path,
				OOB:  true,
			},
			Breadcrumbs: view.Breadcrumbs(listing.Path, true),
			Search:      view.SearchBox(listing.Path, listing.Options, true),
			EmptyState: view.EmptyPreviewView{
				Path: listing.Path,
			},
			Status: view.Status(listing, h.config.Write, true),
		})
	}

	return h.render(c, "files.html", files)
}
