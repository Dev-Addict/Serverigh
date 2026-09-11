package bulk

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/apperror"
	"serverigh/internal/filesystem"
)

func bulkTargets(c *fiber.Ctx) []string {
	values := c.Request().PostArgs().PeekMulti("target")
	targets := make([]string, 0, len(values))
	for _, value := range values {
		if len(value) > 0 {
			targets = append(targets, string(value))
		}
	}

	return targets
}

func (h Handlers) bulkResolvedTargets(
	c *fiber.Ctx,
) ([]filesystem.ResolvedPath, error) {
	targets := bulkTargets(c)
	if len(targets) == 0 {
		return nil, apperror.New(
			apperror.CodeInvalidPath,
			"select at least one entry",
		)
	}

	resolved := make([]filesystem.ResolvedPath, 0, len(targets))
	seen := map[string]bool{}
	for _, target := range targets {
		entry, err := h.Files.Resolve(target)
		if err != nil {
			return nil, err
		}
		if seen[entry.Path] {
			return nil, apperror.New(
				apperror.CodeInvalidPath,
				"entry selected more than once",
			)
		}
		seen[entry.Path] = true
		resolved = append(resolved, entry)
	}

	return resolved, nil
}

func (h Handlers) bulkDestination(path string) (filesystem.ResolvedPath, error) {
	destination, err := h.Files.Resolve(path)
	if err != nil {
		return filesystem.ResolvedPath{}, err
	}
	if !destination.Info.IsDir() {
		return filesystem.ResolvedPath{}, apperror.New(
			apperror.CodeNotDirectory,
			"path is not a directory",
		)
	}

	return destination, nil
}
