package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
	"serverigh/internal/httpserver/view"
)

func (h Handlers) Browse(c *fiber.Ctx) error {
	activePath := c.Query("path", "/")
	options := listOptionsFromRequest(c)
	listing, err := h.files.ListWithOptions(activePath, options)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	files := view.Files(listing)

	var preview *filesystem.Preview
	selectedFile := c.Query("file")
	if selectedFile != "" {
		selectedPreview, err := h.files.Preview(selectedFile)
		if err != nil {
			return h.writeOperationalError(c, err)
		}

		preview = &selectedPreview
	}

	page := view.BrowsePage{
		Root:    h.config.Root,
		Path:    listing.Path,
		Mode:    view.ModeLabel(h.config.Write),
		Listing: files,
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
	}

	return h.render(c, "browse.html", page)
}
