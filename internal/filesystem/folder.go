package filesystem

import folderfs "serverigh/internal/filesystem/folder"

type FolderTree = folderfs.Tree

type FolderTreeEntry = folderfs.Entry

func (s Service) FolderBranch(selectedPath string) (FolderTree, error) {
	return s.folderService().Branch(selectedPath)
}

func (s Service) FolderChildren(
	parentPath string,
) ([]FolderTreeEntry, bool, error) {
	return s.folderService().Children(parentPath)
}

func (s Service) folderService() folderfs.Service {
	return folderfs.NewService(s.root, s.showHidden)
}
