package handlers

import "github.com/gofiber/fiber/v2"

type filesVersionResponse struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

func (h Handlers) FilesVersion(c *fiber.Ctx) error {
	options := listOptionsFromRequest(c)
	version, err := h.files.DirectoryVersion(c.Query("path", "/"), options)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	return c.JSON(filesVersionResponse{
		Path:    version.Path,
		Version: version.Version,
	})
}
