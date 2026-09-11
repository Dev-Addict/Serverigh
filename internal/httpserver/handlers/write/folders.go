package write

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
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	selectedPath := c.Query("selected")
	if selectedPath != "" {
		tree, err := h.Files.FolderBranch(selectedPath)
		if err != nil {
			return h.WriteOperationalError(c, err)
		}

		logOperation(
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

		return h.Render(c, "write_folder_tree.html", tree)
	}

	parentPath := c.Query("path", "/")
	entries, truncated, err := h.Files.FolderChildren(parentPath)
	if err != nil {
		return h.WriteOperationalError(c, err)
	}

	logOperation(
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

	return h.Render(c, "write_folder_tree.html", folderTreeView{
		Entries:   entries,
		Truncated: truncated,
	})
}
