package handlers

import (
	"net/url"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

type FilesView struct {
	filesystem.DirectoryListing
	Controls TableControlsView
}

type TableControlsView struct {
	Path         string
	BrowseURL    string
	FilesURL     string
	NameSort     TableSortLink
	SizeSort     TableSortLink
	ModifiedSort TableSortLink
	CreatedSort  TableSortLink
}

type TableSortLink struct {
	Label     string
	Indicator string
	AriaSort  string
	BrowseURL string
	FilesURL  string
	Active    bool
}

func filesView(listing filesystem.DirectoryListing) FilesView {
	controls := tableControlsView(listing.Path, listing.Options)

	return FilesView{
		DirectoryListing: listing,
		Controls:         controls,
	}
}

func (v FilesView) BrowseURL(path string) string {
	return browseURLWithOptions(path, v.Options)
}

func (v FilesView) FilesURL(path string) string {
	return filesURL(path, v.Options)
}

func (v FilesView) FileBrowseURL(filePath string) string {
	return browseURLWithFileAndOptions(v.Path, filePath, v.Options)
}

func listOptionsFromRequest(
	c *fiber.Ctx,
) filesystem.ListOptions {
	return filesystem.NormalizeListOptions(filesystem.ListOptions{
		Sort:      filesystem.ListSort(c.Query("sort")),
		Direction: filesystem.ListDirection(c.Query("dir")),
	})
}

func tableControlsView(
	activePath string,
	options filesystem.ListOptions,
) TableControlsView {
	options = filesystem.NormalizeListOptions(options)
	controls := TableControlsView{
		Path:      activePath,
		BrowseURL: browseURLWithOptions(activePath, options),
		FilesURL:  filesURL(activePath, options),
	}

	controls.NameSort = sortLink("Name", filesystem.ListSortName, activePath, options)
	controls.SizeSort = sortLink("Size", filesystem.ListSortSize, activePath, options)
	controls.ModifiedSort = sortLink(
		"Modified",
		filesystem.ListSortModified,
		activePath,
		options,
	)
	controls.CreatedSort = sortLink(
		"Created",
		filesystem.ListSortCreated,
		activePath,
		options,
	)

	return controls
}

func sortLink(
	label string,
	field filesystem.ListSort,
	activePath string,
	options filesystem.ListOptions,
) TableSortLink {
	nextOptions := options
	active := options.Sort == field
	direction := filesystem.ListDirectionAsc
	if active && options.Direction == filesystem.ListDirectionAsc {
		direction = filesystem.ListDirectionDesc
	}
	nextOptions.Sort = field
	nextOptions.Direction = direction

	return TableSortLink{
		Label:     label,
		Indicator: sortIndicator(active, options.Direction),
		AriaSort:  ariaSort(active, options.Direction),
		BrowseURL: browseURLWithOptions(activePath, nextOptions),
		FilesURL:  filesURL(activePath, nextOptions),
		Active:    active,
	}
}

func browseURLWithOptions(
	activePath string,
	options filesystem.ListOptions,
) string {
	return "/browse?" + listURLValues(activePath, options).Encode()
}

func browseURLWithFileAndOptions(
	activePath string,
	filePath string,
	options filesystem.ListOptions,
) string {
	values := listURLValues(activePath, options)
	values.Set("file", filePath)

	return "/browse?" + values.Encode()
}

func filesURL(activePath string, options filesystem.ListOptions) string {
	return "/partials/files?" + listURLValues(activePath, options).Encode()
}

func listURLValues(activePath string, options filesystem.ListOptions) url.Values {
	options = filesystem.NormalizeListOptions(options)
	values := url.Values{}
	values.Set("path", activePath)
	values.Set("sort", string(options.Sort))
	values.Set("dir", string(options.Direction))

	return values
}

func sortIndicator(active bool, direction filesystem.ListDirection) string {
	if !active {
		return ""
	}

	if direction == filesystem.ListDirectionDesc {
		return "Desc"
	}

	return "Asc"
}

func ariaSort(active bool, direction filesystem.ListDirection) string {
	if !active {
		return "none"
	}

	if direction == filesystem.ListDirectionDesc {
		return "descending"
	}

	return "ascending"
}
