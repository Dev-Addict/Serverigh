package searcher

type Service struct {
	root       string
	showHidden bool
	Index      *Index
}

func NewService(root string, showHidden bool, index *Index) Service {
	if index == nil {
		index = NewIndex()
	}

	return Service{
		root:       root,
		showHidden: showHidden,
		Index:      index,
	}
}
