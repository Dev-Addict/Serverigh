package view

import (
	"path"
	"strings"
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

func Breadcrumbs(activePath string, oob bool) BreadcrumbsView {
	return BreadcrumbsView{
		BackPath: ParentPath(activePath),
		Items:    BreadcrumbsForPath(activePath),
		OOB:      oob,
	}
}

func BreadcrumbsForPath(activePath string) []Breadcrumb {
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
