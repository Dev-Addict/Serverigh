package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

type BrowseView struct {
	Root        string
	Path        string
	Mode        string
	Listing     FilesView
	PathSummary PathSummaryView
	Breadcrumbs BreadcrumbsView
	Search      SearchBoxView
	EmptyState  EmptyPreviewView
	Status      StatusView
	Preview     *filesystem.Preview
}

func (h Handlers) Browse(c *fiber.Ctx) error {
	activePath := c.Query("path", "/")
	options := listOptionsFromRequest(c)
	listing, err := h.files.ListWithOptions(activePath, options)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	files := filesView(listing)

	var preview *filesystem.Preview
	selectedFile := c.Query("file")
	if selectedFile != "" {
		selectedPreview, err := h.files.Preview(selectedFile)
		if err != nil {
			return h.writeOperationalError(c, err)
		}

		preview = &selectedPreview
	}

	view := BrowseView{
		Root:    h.config.Root,
		Path:    listing.Path,
		Mode:    modeLabel(h.config.Write),
		Listing: files,
		PathSummary: PathSummaryView{
			Root: h.config.Root,
			Path: listing.Path,
		},
		Breadcrumbs: breadcrumbsView(listing.Path, false),
		Search:      searchBoxView(listing.Path, listing.Options, false),
		EmptyState: EmptyPreviewView{
			Path: listing.Path,
		},
		Preview: preview,
		Status:  statusView(listing, h.config.Write, false),
	}

	return h.render(c, "browse.html", view)
}

func statusLabel(listing filesystem.DirectoryListing) string {
	if listing.Truncated {
		return "Showing first entries"
	}

	if len(listing.Entries) == 0 {
		return "Folder empty"
	}

	return "Ready"
}

func statusView(
	listing filesystem.DirectoryListing,
	writeEnabled bool,
	oob bool,
) StatusView {
	return StatusView{
		Label:      statusLabel(listing),
		ItemCount:  len(listing.Entries),
		Mode:       modeLabel(writeEnabled),
		Truncated:  listing.Truncated,
		EntryLimit: listing.EntryLimit,
		OOB:        oob,
	}
}
