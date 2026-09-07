package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/httpserver/view"
)

func (h Handlers) Files(c *fiber.Ctx) error {
	options := listOptionsFromRequest(c)
	listing, err := h.files.ListWithOptions(c.Query("path", "/"), options)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	files := view.Files(listing)

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
