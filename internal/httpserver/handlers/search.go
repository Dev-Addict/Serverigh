package handlers

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
	"serverigh/internal/httpserver/view"
)

const searchRequestTimeout = 5 * time.Second

func (h Handlers) Search(c *fiber.Ctx) error {
	start := time.Now()
	activePath := c.Query("path", "/")
	options := listOptionsFromRequest(c)
	if strings.TrimSpace(c.Query("q")) == "" {
		listing, err := h.files.ListWithOptions(activePath, options)
		if err != nil {
			return h.writeOperationalError(c, err)
		}
		logHandlerOperation(
			c,
			slog.LevelInfo,
			"search cleared",
			start,
			"path",
			listing.Path,
			"entry_count",
			len(listing.Entries),
		)

		return h.render(c, "files.html", view.Files(listing, h.config.Write))
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
	logHandlerOperation(
		c,
		slog.LevelInfo,
		"search completed",
		start,
		"path",
		activePath,
		"query_length",
		len(c.Query("q")),
		"result_count",
		len(results.Results),
	)

	return h.render(c, "search_panel.html", view.SearchResults(results))
}
