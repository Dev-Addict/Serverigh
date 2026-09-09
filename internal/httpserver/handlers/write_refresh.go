package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/apperror"
	"serverigh/internal/httpserver/view"
)

func (h Handlers) writeRefresh(c *fiber.Ctx, path string) error {
	options := listOptionsFromRequest(c)
	listing, err := h.files.ListWithOptions(path, options)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	files := view.Files(listing, h.config.Write)
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

func (h Handlers) writeModeDisabled(c *fiber.Ctx) (bool, error) {
	if h.config.Write {
		return false, nil
	}

	err := h.writeError(
		c,
		fiber.StatusForbidden,
		apperror.CodeWriteDisabled,
		"write mode is disabled",
	)

	return true, err
}
