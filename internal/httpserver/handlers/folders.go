package handlers

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

type folderTreeView struct {
	Entries   []filesystem.FolderTreeEntry
	Truncated bool
}

func (h Handlers) Folders(c *fiber.Ctx) error {
	if disabled, err := h.writeModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	selectedPath := c.Query("selected")
	if selectedPath != "" {
		tree, err := h.files.FolderBranch(selectedPath)
		if err != nil {
			return h.writeOperationalError(c, err)
		}
		logHandlerOperation(
			c,
			slog.LevelInfo,
			"folder branch listed",
			start,
			"selected_path",
			selectedPath,
			"entry_count",
			len(tree.Entries),
			"truncated",
			tree.Truncated,
		)

		return h.render(c, "write_folder_tree.html", tree)
	}

	parentPath := c.Query("path", "/")
	entries, truncated, err := h.files.FolderChildren(parentPath)
	if err != nil {
		return h.writeOperationalError(c, err)
	}
	logHandlerOperation(
		c,
		slog.LevelInfo,
		"folder children listed",
		start,
		"path",
		parentPath,
		"entry_count",
		len(entries),
		"truncated",
		truncated,
	)

	return h.render(c, "write_folder_tree.html", folderTreeView{
		Entries:   entries,
		Truncated: truncated,
	})
}
