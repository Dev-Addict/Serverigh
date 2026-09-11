package handlers

import (
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
	"serverigh/internal/httpserver/view"
)

func (h Handlers) Breadcrumbs(c *fiber.Ctx) error {
	activePath := c.Query("path", "/")
	if path.Clean("/"+strings.TrimPrefix(activePath, "/")) == filesystem.TrashPath {
		return h.render(c, "breadcrumbs.html", view.Breadcrumbs(
			filesystem.TrashPath,
			false,
		))
	}

	resolved, err := h.files.Resolve(activePath)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	if !resolved.Info.IsDir() {
		return h.writeOperationalError(c, filesystem.ErrNotDirectory)
	}

	return h.render(c, "breadcrumbs.html", view.Breadcrumbs(resolved.Path, false))
}
