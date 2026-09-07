package handlers

import (
	"net/url"
	"path"
	"strings"
)

func browseURL(activePath string) string {
	return "/browse?path=" + url.QueryEscape(activePath)
}

func browseURLWithFile(activePath string, filePath string) string {
	return browseURL(activePath) + "&file=" + url.QueryEscape(filePath)
}

func folderForFile(filePath string) string {
	return parentOfPath(filePath)
}

func parentPath(activePath string) string {
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
