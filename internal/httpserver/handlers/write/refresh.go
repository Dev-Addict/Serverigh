package write

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/apperror"
	"serverigh/internal/filesystem"
	"serverigh/internal/httpserver/view"
)

func (h Handlers) Refresh(c *fiber.Ctx, path string) error {
	options := listOptionsFromRequest(c)
	listing, err := h.Files.ListWithOptions(path, options)
	if err != nil {
		return h.WriteOperationalError(c, err)
	}

	files := view.Files(listing, h.Config.Enabled)
	c.Set("HX-Push-Url", files.Controls.BrowseURL)

	return h.Render(c, "files_update.html", view.FilesUpdate{
		Listing: files,
		PathSummary: view.PathSummaryView{
			Root: h.Config.Root,
			Path: listing.Path,
			OOB:  true,
		},
		Breadcrumbs: view.Breadcrumbs(listing.Path, true),
		Search:      view.SearchBox(listing.Path, listing.Options, true),
		EmptyState: view.EmptyPreviewView{
			Path: listing.Path,
		},
		Status: view.Status(listing, h.Config.Enabled, true),
	})
}

func (h Handlers) ModeDisabled(c *fiber.Ctx) (bool, error) {
	if h.Config.Enabled {
		return false, nil
	}

	err := h.WriteError(
		c,
		fiber.StatusForbidden,
		apperror.CodeWriteDisabled,
		"write mode is disabled",
	)

	return true, err
}

func listOptionsFromRequest(c *fiber.Ctx) filesystem.ListOptions {
	return filesystem.NormalizeListOptions(filesystem.ListOptions{
		Sort:      filesystem.ListSort(c.Query("sort")),
		Direction: filesystem.ListDirection(c.Query("dir")),
	})
}
