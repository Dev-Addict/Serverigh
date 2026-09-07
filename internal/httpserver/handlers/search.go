package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
	"serverigh/internal/httpserver/view"
)

const searchRequestTimeout = 5 * time.Second

func (h Handlers) Search(c *fiber.Ctx) error {
	activePath := c.Query("path", "/")
	options := listOptionsFromRequest(c)
	if strings.TrimSpace(c.Query("q")) == "" {
		listing, err := h.files.ListWithOptions(activePath, options)
		if err != nil {
			return h.writeOperationalError(c, err)
		}

		return h.render(c, "files.html", view.Files(listing))
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), searchRequestTimeout)
	defer cancel()

	results, err := h.files.Search(filesystem.SearchOptions{
		Context:    ctx,
		Query:      c.Query("q"),
		ActivePath: activePath,
	})
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	return h.render(c, "search_panel.html", view.SearchResults(results))
}
