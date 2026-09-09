package view

import (
	"path/filepath"
	"sort"
	"strings"

	"serverigh/internal/filesystem"
)

type FileIconView struct {
	Name  string
	Path  string
	Label string
}

var compoundIconExtensions = sortedIconExtensions(compoundExtensionIcons)

func (v FilesView) Icon(entry filesystem.Entry) FileIconView {
	return FileIcon(entry.Name, entry.Kind, entry.IsDir)
}

func FileIcon(name string, kind string, isDir bool) FileIconView {
	icon := iconName(name, kind, isDir)

	return FileIconView{
		Name:  icon,
		Path:  "/static/icons/files/" + icon + ".svg",
		Label: iconLabel(icon),
	}
}

func iconName(name string, kind string, isDir bool) string {
	if isDir || kind == "folder" {
		return "folder"
	}
	if kind == "symlink" {
		return "symlink"
	}

	lowerName := strings.ToLower(name)
	if icon, ok := exactFileIcons[lowerName]; ok {
		return icon
	}

	if icon, ok := compoundIcon(lowerName); ok {
		return icon
	}

	if icon, ok := extensionIcons[filepath.Ext(lowerName)]; ok {
		return icon
	}

	return "file"
}

func iconLabel(icon string) string {
	return strings.ReplaceAll(icon, "-", " ") + " file"
}

func compoundIcon(name string) (string, bool) {
	for _, extension := range compoundIconExtensions {
		if strings.HasSuffix(name, extension) {
			return compoundExtensionIcons[extension], true
		}
	}

	return "", false
}

func sortedIconExtensions(icons map[string]string) []string {
	extensions := make([]string, 0, len(icons))
	for extension := range icons {
		extensions = append(extensions, extension)
	}
	sort.Slice(extensions, func(left int, right int) bool {
		return len(extensions[left]) > len(extensions[right])
	})

	return extensions
}
