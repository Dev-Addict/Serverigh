package fileinfo

import "os"

func Kind(info os.FileInfo) string {
	if info.IsDir() {
		return "folder"
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return "symlink"
	}

	return "file"
}
