package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/apperror"
	"serverigh/internal/config"
	"serverigh/internal/filesystem"
)

func (h Handlers) Settings(c *fiber.Ctx) error {
	settings, err := h.settingsFromRequest(c)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	if err := config.SaveLocalSettings(h.config.Root, settings); err != nil {
		return h.writeOperationalError(c, err)
	}

	files, err := filesystem.New(
		h.config.Root,
		h.config.ShowHidden,
		settings.MaxPreviewBytes,
	)
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	*h.files = files
	h.config.Theme = settings.Theme
	h.config.MaxPreviewBytes = settings.MaxPreviewBytes
	h.config.Columns = settings.Columns

	return c.JSON(fiber.Map{"status": "ok"})
}

func (h Handlers) settingsFromRequest(c *fiber.Ctx) (config.Settings, error) {
	maxPreviewBytes, err := strconv.ParseInt(
		c.FormValue("max_preview_bytes"),
		10,
		64,
	)
	if err != nil {
		return config.Settings{}, apperror.WrapOperation(
			apperror.CodeInvalidConfig,
			"invalid config",
			"parse max preview bytes",
			err,
		)
	}

	settings := config.Settings{
		Theme:           c.FormValue("theme"),
		MaxPreviewBytes: maxPreviewBytes,
		Columns: config.Columns{
			Size:     requestBool(c, "column_size"),
			Modified: requestBool(c, "column_modified"),
			Created:  requestBool(c, "column_created"),
			Mode:     requestBool(c, "column_mode"),
		},
	}

	return settings, settings.Validate()
}

func requestBool(c *fiber.Ctx, name string) bool {
	value, err := strconv.ParseBool(c.FormValue(name))
	return err == nil && value
}
