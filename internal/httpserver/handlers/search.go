package handlers

import (
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

type SearchBoxView struct {
	Path      string
	Sort      string
	Direction string
	OOB       bool
}

type SearchResultsView struct {
	Query     string
	Results   []SearchResultView
	Truncated bool
}

type SearchResultView struct {
	Name         string
	ParentPath   string
	Kind         string
	SizeLabel    string
	ModTimeLabel string
	CreatedLabel string
	Mode         string
	IsDir        bool
	BrowseURL    string
	FilesURL     string
	PreviewURL   string
}

func (h Handlers) Search(c *fiber.Ctx) error {
	activePath := c.Query("path", "/")
	options := listOptionsFromRequest(c)
	if strings.TrimSpace(c.Query("q")) == "" {
		listing, err := h.files.ListWithOptions(activePath, options)
		if err != nil {
			return h.writeOperationalError(c, err)
		}

		return h.render(c, "files.html", filesView(listing))
	}

	results, err := h.files.Search(filesystem.SearchOptions{
		Query:      c.Query("q"),
		ActivePath: activePath,
	})
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	return h.render(c, "search_panel.html", searchResultsView(results))
}

func searchBoxView(
	activePath string,
	options filesystem.ListOptions,
	oob bool,
) SearchBoxView {
	options = filesystem.NormalizeListOptions(options)

	return SearchBoxView{
		Path:      activePath,
		Sort:      string(options.Sort),
		Direction: string(options.Direction),
		OOB:       oob,
	}
}

func searchResultsView(results filesystem.SearchResults) SearchResultsView {
	view := SearchResultsView{
		Query:     results.Query,
		Truncated: results.Truncated,
	}

	for _, result := range results.Results {
		view.Results = append(view.Results, SearchResultView{
			Name:         result.Name,
			ParentPath:   result.ParentPath,
			Kind:         result.Kind,
			SizeLabel:    result.SizeLabel,
			ModTimeLabel: result.ModTimeLabel,
			CreatedLabel: result.CreatedLabel,
			Mode:         result.Mode,
			IsDir:        result.IsDir,
			BrowseURL:    searchBrowseURL(result),
			FilesURL:     searchFilesURL(result),
			PreviewURL:   previewURL(result.Path),
		})
	}

	return view
}

func searchBrowseURL(result filesystem.SearchResult) string {
	if result.IsDir {
		return browseURL(result.Path)
	}

	return browseURLWithFile(result.ParentPath, result.Path)
}

func searchFilesURL(result filesystem.SearchResult) string {
	if !result.IsDir {
		return ""
	}

	return filesURL(result.Path, filesystem.DefaultListOptions())
}

func previewURL(filePath string) string {
	return "/preview?path=" + url.QueryEscape(filePath)
}
