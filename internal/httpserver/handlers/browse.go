package handlers

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
	"serverigh/internal/httpserver/view"
)

func (h Handlers) Browse(c *fiber.Ctx) error {
	start := time.Now()
	activePath := c.Query("path", "/")
	options := listOptionsFromRequest(c)
	listing, err := h.files.ListWithOptions(activePath, options)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	files := view.Files(listing, h.config.Write)

	var preview *filesystem.Preview
	selectedFile := c.Query("file")
	if selectedFile != "" {
		selectedPreview, err := h.files.Preview(selectedFile)
		if err != nil {
			return h.writeOperationalError(c, err)
		}

		preview = &selectedPreview
	}
	logHandlerOperation(
		c,
		slog.LevelInfo,
		"browse rendered",
		start,
		"path",
		listing.Path,
		"entry_count",
		len(listing.Entries),
		"truncated",
		listing.Truncated,
		"has_preview",
		preview != nil,
	)

	page := view.BrowsePage{
		Root:            h.config.Root,
		Path:            listing.Path,
		Mode:            view.ModeLabel(h.config.Write),
		Theme:           h.config.Theme,
		MaxPreviewBytes: h.config.MaxPreviewBytes,
		Listing:         files,
		PathSummary: view.PathSummaryView{
			Root: h.config.Root,
			Path: listing.Path,
		},
		Breadcrumbs: view.Breadcrumbs(listing.Path, false),
		Search:      view.SearchBox(listing.Path, listing.Options, false),
		EmptyState: view.EmptyPreviewView{
			Path: listing.Path,
		},
		Preview: preview,
		Status:  view.Status(listing, h.config.Write, false),
		Columns: h.config.Columns,
	}

	return h.render(c, "browse.html", page)
}
