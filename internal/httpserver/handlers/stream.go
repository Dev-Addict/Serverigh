package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

func sendOpenedFile(c *fiber.Ctx, file filesystem.OpenedFile) error {
	size := streamSize(file.Size)
	if size < 0 {
		if err := c.SendStream(file.Handle); err != nil {
			file.Close()

			return err
		}

		return nil
	}

	if err := c.SendStream(file.Handle, size); err != nil {
		file.Close()

		return err
	}

	return nil
}

func streamSize(size int64) int {
	maxInt := int64(int(^uint(0) >> 1))
	if size < 0 || size > maxInt {
		return -1
	}

	return int(size)
}
