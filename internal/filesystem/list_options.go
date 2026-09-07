package filesystem

type ListSort string

const (
	ListSortName     ListSort = "name"
	ListSortSize     ListSort = "size"
	ListSortModified ListSort = "modified"
	ListSortCreated  ListSort = "created"
)

type ListDirection string

const (
	ListDirectionAsc  ListDirection = "asc"
	ListDirectionDesc ListDirection = "desc"
)

type ListOptions struct {
	Sort      ListSort
	Direction ListDirection
}

func DefaultListOptions() ListOptions {
	return NormalizeListOptions(ListOptions{})
}

func NormalizeListOptions(options ListOptions) ListOptions {
	switch options.Sort {
	case ListSortName,
		ListSortSize,
		ListSortModified,
		ListSortCreated:
	default:
		options.Sort = ListSortName
	}

	switch options.Direction {
	case ListDirectionAsc, ListDirectionDesc:
	default:
		options.Direction = ListDirectionAsc
	}

	return options
}
