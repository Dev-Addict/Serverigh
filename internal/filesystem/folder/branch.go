package folder

func (s Service) Branch(selectedPath string) (Tree, error) {
	selected, err := s.resolve(selectedPath)
	if err != nil {
		return Tree{}, err
	}
	if !selected.Info.IsDir() {
		return Tree{}, ErrNotDirectory
	}

	tree := Tree{}
	treeAncestors := ancestors(selected.Path)
	tree.append(Entry{
		Name:        "/",
		Path:        "/",
		Depth:       0,
		Expanded:    selected.Path != "/",
		HasChildren: true,
		Selected:    selected.Path == "/",
	})

	for index := 0; index < len(treeAncestors)-1; index++ {
		children, truncated, err := s.Children(treeAncestors[index])
		if err != nil {
			return Tree{}, err
		}
		tree.Truncated = tree.Truncated || truncated
		tree.appendBranchChildren(
			children,
			treeAncestors[index+1],
			selected.Path,
		)
	}

	return tree, nil
}

func (tree *Tree) appendBranchChildren(
	children []Entry,
	nextAncestor string,
	selectedPath string,
) {
	for _, child := range children {
		if child.Path == nextAncestor {
			child.Expanded = child.Path != selectedPath
			child.Selected = child.Path == selectedPath
		}
		tree.append(child)
	}
}
