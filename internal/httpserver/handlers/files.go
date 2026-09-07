package handlers

import "github.com/gofiber/fiber/v2"

type FilesUpdateView struct {
	Listing     FilesView
	PathSummary PathSummaryView
	Breadcrumbs BreadcrumbsView
	Search      SearchBoxView
	EmptyState  EmptyPreviewView
	Status      StatusView
}

func (h Handlers) Files(c *fiber.Ctx) error {
	options := listOptionsFromRequest(c)
	listing, err := h.files.ListWithOptions(c.Query("path", "/"), options)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	files := filesView(listing)

	if c.Get("HX-Request") == "true" {
		c.Set("HX-Push-Url", files.Controls.BrowseURL)

		return h.render(c, "files_update.html", FilesUpdateView{
			Listing: files,
			PathSummary: PathSummaryView{
				Root: h.config.Root,
				Path: listing.Path,
				OOB:  true,
			},
			Breadcrumbs: breadcrumbsView(listing.Path, true),
			Search:      searchBoxView(listing.Path, listing.Options, true),
			EmptyState: EmptyPreviewView{
				Path: listing.Path,
			},
			Status: statusView(listing, h.config.Write, true),
		})
	}

	return h.render(c, "files.html", files)
}
