package write

type Service struct {
	root string
}

func NewService(root string) Service {
	return Service{
		root: root,
	}
}
