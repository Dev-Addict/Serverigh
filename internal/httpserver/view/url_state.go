package view

import (
	"net/url"
	"path"
	"strings"

	"serverigh/internal/filesystem"
)

func BrowseURL(activePath string) string {
	return "/browse?path=" + url.QueryEscape(activePath)
}

func BrowseURLWithFile(activePath string, filePath string) string {
	return BrowseURL(activePath) + "&file=" + url.QueryEscape(filePath)
}

func BrowseURLWithOptions(
	activePath string,
	options filesystem.ListOptions,
) string {
	return "/browse?" + listURLValues(activePath, options).Encode()
}

func BrowseURLWithFileAndOptions(
	activePath string,
	filePath string,
	options filesystem.ListOptions,
) string {
	values := listURLValues(activePath, options)
	values.Set("file", filePath)

	return "/browse?" + values.Encode()
}

func FilesURL(activePath string, options filesystem.ListOptions) string {
	return "/partials/files?" + listURLValues(activePath, options).Encode()
}

func FolderForFile(filePath string) string {
	return parentOfPath(filePath)
}

func ParentPath(activePath string) string {
	cleanPath := cleanURLPath(activePath)
	if cleanPath == "/" {
		return "/"
	}

	return parentOfPath(cleanPath)
}

func parentOfPath(activePath string) string {
	folder := path.Dir(cleanURLPath(activePath))
	if folder == "." {
		return "/"
	}

	return folder
}

func cleanURLPath(activePath string) string {
	cleanPath := path.Clean("/" + strings.TrimPrefix(activePath, "/"))
	if cleanPath == "." {
		return "/"
	}

	return cleanPath
}

func listURLValues(activePath string, options filesystem.ListOptions) url.Values {
	options = filesystem.NormalizeListOptions(options)
	values := url.Values{}
	values.Set("path", activePath)
	values.Set("sort", string(options.Sort))
	values.Set("dir", string(options.Direction))

	return values
}
