package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

type FilesUpdateView struct {
	Listing     filesystem.DirectoryListing
	PathSummary PathSummaryView
	Breadcrumbs BreadcrumbsView
	EmptyState  EmptyPreviewView
	Status      StatusView
}

func (h Handlers) Files(c *fiber.Ctx) error {
	listing, err := h.files.List(c.Query("path", "/"))
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	if c.Get("HX-Request") == "true" {
		return h.render(c, "files_update.html", FilesUpdateView{
			Listing: listing,
			PathSummary: PathSummaryView{
				Root: h.config.Root,
				Path: listing.Path,
				OOB:  true,
			},
			Breadcrumbs: BreadcrumbsView{
				Items: breadcrumbsForPath(listing.Path),
				OOB:   true,
			},
			EmptyState: EmptyPreviewView{
				Path: listing.Path,
			},
			Status: statusView(listing, h.config.Write, true),
		})
	}

	return h.render(c, "files.html", listing)
}
