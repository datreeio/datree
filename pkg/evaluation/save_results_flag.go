package evaluation

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func IsWritableDirectory(filePath string) bool {
	dirPath := filepath.Dir(filePath)

	// Check if the directory exists
	_, err := os.Lstat(dirPath)
	if os.IsNotExist(err) {
		return false
	}

	// Check if you can edit the file
	if unix.Access(filePath, unix.W_OK) == nil {
		return true
	}

	// Check if you can create a file in the directory
	dirInfo, err := os.Stat(dirPath)
	if err == nil && dirInfo.Mode().Perm()&0200 != 0 {
		return true
	}

	extension := filepath.Ext(filePath)
	if extension != ".json" {
		return false
	}

	return false
}
