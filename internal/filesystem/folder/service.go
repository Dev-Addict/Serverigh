package folder

type Service struct {
	root       string
	showHidden bool
}

func NewService(root string, showHidden bool) Service {
	return Service{
		root:       root,
		showHidden: showHidden,
	}
}
