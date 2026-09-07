package filesystem

import "sort"

func sortedListing(listing DirectoryListing) DirectoryListing {
	sort.SliceStable(listing.Entries, func(i int, j int) bool {
		left := listing.Entries[i]
		right := listing.Entries[j]
		if left.IsDir != right.IsDir {
			return left.IsDir
		}

		return compareEntries(left, right, listing.Options) < 0
	})

	return listing
}

func compareEntries(left Entry, right Entry, options ListOptions) int {
	result := 0
	switch options.Sort {
	case ListSortSize:
		result = compareInt64(left.Size, right.Size)
	case ListSortModified:
		result = compareTimes(left.ModTime, right.ModTime)
	case ListSortCreated:
		result = compareOptionalTimes(
			left.CreatedTime,
			left.CreatedKnown,
			right.CreatedTime,
			right.CreatedKnown,
		)
	default:
		result = compareNames(left.Name, right.Name)
	}

	if result == 0 {
		result = compareNames(left.Name, right.Name)
	}

	if options.Direction == ListDirectionDesc {
		return -result
	}

	return result
}
