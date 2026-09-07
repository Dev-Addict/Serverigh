package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

type BrowseView struct {
	Root        string
	Path        string
	Mode        string
	Listing     filesystem.DirectoryListing
	PathSummary PathSummaryView
	Breadcrumbs BreadcrumbsView
	EmptyState  EmptyPreviewView
	Status      StatusView
	Preview     *filesystem.Preview
}

func (h Handlers) Browse(c *fiber.Ctx) error {
	activePath := c.Query("path", "/")
	listing, err := h.files.List(activePath)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

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
		Listing: listing,
		PathSummary: PathSummaryView{
			Root: h.config.Root,
			Path: listing.Path,
		},
		Breadcrumbs: breadcrumbsView(listing.Path, false),
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
