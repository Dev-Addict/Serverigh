package view

import (
	"net/url"

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
	RelativePath string
	AbsolutePath string
	Kind         string
	SizeLabel    string
	ModTimeLabel string
	ModTimeValue string
	CreatedLabel string
	CreatedValue string
	CreatedKnown bool
	Mode         string
	IsDir        bool
	Icon         FileIconView
	BrowseURL    string
	FilesURL     string
	PreviewURL   string
}

func SearchBox(
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

func SearchResults(results filesystem.SearchResults) SearchResultsView {
	view := SearchResultsView{
		Query:     results.Query,
		Truncated: results.Truncated,
	}

	for _, result := range results.Results {
		view.Results = append(view.Results, SearchResultView{
			Name:         result.Name,
			ParentPath:   result.ParentPath,
			RelativePath: result.RelativePath,
			AbsolutePath: result.AbsolutePath,
			Kind:         result.Kind,
			SizeLabel:    result.SizeLabel,
			ModTimeLabel: result.ModTimeLabel,
			ModTimeValue: result.ModTimeValue,
			CreatedLabel: result.CreatedLabel,
			CreatedValue: result.CreatedValue,
			CreatedKnown: result.CreatedKnown,
			Mode:         result.Mode,
			IsDir:        result.IsDir,
			Icon:         FileIcon(result.Name, result.Kind, result.IsDir),
			BrowseURL:    searchBrowseURL(result),
			FilesURL:     searchFilesURL(result),
			PreviewURL:   previewURL(result.Path),
		})
	}

	return view
}

func searchBrowseURL(result filesystem.SearchResult) string {
	if result.IsDir {
		return BrowseURL(result.Path)
	}

	return BrowseURLWithFile(result.ParentPath, result.Path)
}

func searchFilesURL(result filesystem.SearchResult) string {
	if !result.IsDir {
		return ""
	}

	return FilesURL(result.Path, filesystem.DefaultListOptions())
}

func previewURL(filePath string) string {
	return "/preview?path=" + url.QueryEscape(filePath)
}
