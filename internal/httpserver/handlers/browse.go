package handlers

import "github.com/gofiber/fiber/v2"

type BrowseView struct {
	Root string
	Path string
	Mode string
}

func (h Handlers) Browse(c *fiber.Ctx) error {
	activePath := c.Query("path", "/")
	mode := "read-only"
	if h.config.Write {
		mode = "write-enabled"
	}

	view := BrowseView{
		Root: h.config.Root,
		Path: activePath,
		Mode: mode,
	}

	return h.render(c, "browse.html", view)
}
