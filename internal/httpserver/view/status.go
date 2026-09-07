package view

import "serverigh/internal/filesystem"

type StatusView struct {
	Label      string
	ItemCount  int
	Mode       string
	Truncated  bool
	EntryLimit int
	OOB        bool
}

func ModeLabel(writeEnabled bool) string {
	if writeEnabled {
		return "write-enabled"
	}

	return "read-only"
}

func Status(
	listing filesystem.DirectoryListing,
	writeEnabled bool,
	oob bool,
) StatusView {
	return StatusView{
		Label:      StatusLabel(listing),
		ItemCount:  len(listing.Entries),
		Mode:       ModeLabel(writeEnabled),
		Truncated:  listing.Truncated,
		EntryLimit: listing.EntryLimit,
		OOB:        oob,
	}
}

func StatusLabel(listing filesystem.DirectoryListing) string {
	if listing.Truncated {
		return "Showing first entries"
	}

	if len(listing.Entries) == 0 {
		return "Folder empty"
	}

	return "Ready"
}
