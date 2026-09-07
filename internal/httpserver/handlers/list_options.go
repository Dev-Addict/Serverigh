package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

func listOptionsFromRequest(c *fiber.Ctx) filesystem.ListOptions {
	return filesystem.NormalizeListOptions(filesystem.ListOptions{
		Sort:      filesystem.ListSort(c.Query("sort")),
		Direction: filesystem.ListDirection(c.Query("dir")),
	})
}
