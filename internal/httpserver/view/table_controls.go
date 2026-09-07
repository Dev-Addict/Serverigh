package view

import "serverigh/internal/filesystem"

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

func Files(listing filesystem.DirectoryListing) FilesView {
	controls := TableControls(listing.Path, listing.Options)

	return FilesView{
		DirectoryListing: listing,
		Controls:         controls,
	}
}

func (v FilesView) BrowseURL(path string) string {
	return BrowseURLWithOptions(path, v.Options)
}

func (v FilesView) FilesURL(path string) string {
	return FilesURL(path, v.Options)
}

func (v FilesView) FileBrowseURL(filePath string) string {
	return BrowseURLWithFileAndOptions(v.Path, filePath, v.Options)
}

func TableControls(
	activePath string,
	options filesystem.ListOptions,
) TableControlsView {
	options = filesystem.NormalizeListOptions(options)
	controls := TableControlsView{
		Path:      activePath,
		BrowseURL: BrowseURLWithOptions(activePath, options),
		FilesURL:  FilesURL(activePath, options),
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
		BrowseURL: BrowseURLWithOptions(activePath, nextOptions),
		FilesURL:  FilesURL(activePath, nextOptions),
		Active:    active,
	}
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
