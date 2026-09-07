package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
	"serverigh/internal/httpserver/view"
)

func (h Handlers) Breadcrumbs(c *fiber.Ctx) error {
	resolved, err := h.files.Resolve(c.Query("path", "/"))
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	if !resolved.Info.IsDir() {
		return h.writeOperationalError(c, filesystem.ErrNotDirectory)
	}

	return h.render(c, "breadcrumbs.html", view.Breadcrumbs(resolved.Path, false))
}
