package handlers

import (
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

type BreadcrumbsView struct {
	BackPath string
	Items    []Breadcrumb
	OOB      bool
}

type Breadcrumb struct {
	Name    string
	Path    string
	Current bool
}

func (h Handlers) Breadcrumbs(c *fiber.Ctx) error {
	resolved, err := h.files.Resolve(c.Query("path", "/"))
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	if !resolved.Info.IsDir() {
		return h.writeOperationalError(c, filesystem.ErrNotDirectory)
	}

	return h.render(c, "breadcrumbs.html", BreadcrumbsView{
		BackPath: parentPath(resolved.Path),
		Items:    breadcrumbsForPath(resolved.Path),
	})
}

func breadcrumbsView(activePath string, oob bool) BreadcrumbsView {
	return BreadcrumbsView{
		BackPath: parentPath(activePath),
		Items:    breadcrumbsForPath(activePath),
		OOB:      oob,
	}
}

func breadcrumbsForPath(activePath string) []Breadcrumb {
	cleanPath := path.Clean("/" + strings.TrimPrefix(activePath, "/"))
	if cleanPath == "." {
		cleanPath = "/"
	}

	items := []Breadcrumb{{
		Name:    "Root",
		Path:    "/",
		Current: cleanPath == "/",
	}}

	if cleanPath == "/" {
		return items
	}

	parts := strings.Split(strings.TrimPrefix(cleanPath, "/"), "/")
	for i, part := range parts {
		itemPath := "/" + strings.Join(parts[:i+1], "/")
		items = append(items, Breadcrumb{
			Name:    part,
			Path:    itemPath,
			Current: i == len(parts)-1,
		})
	}

	return items
}
