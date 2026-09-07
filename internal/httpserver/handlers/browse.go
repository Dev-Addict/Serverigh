package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

type BrowseView struct {
	Root        string
	Path        string
	Mode        string
	ItemCount   int
	StatusLabel string
	Listing     filesystem.DirectoryListing
}

func (h Handlers) Browse(c *fiber.Ctx) error {
	activePath := c.Query("path", "/")
	listing, err := h.files.List(activePath)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	mode := "read-only"
	if h.config.Write {
		mode = "write-enabled"
	}

	view := BrowseView{
		Root:        h.config.Root,
		Path:        listing.Path,
		Mode:        mode,
		ItemCount:   len(listing.Entries),
		StatusLabel: statusLabel(listing),
		Listing:     listing,
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
