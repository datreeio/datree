package evaluation

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsWritableDirectory(directoryPath string) bool {
	// Check if the file exists
	if fileInfo, err := os.Stat(directoryPath); os.IsNotExist(err) {
		extension := filepath.Ext(directoryPath)
		if extension != ".json" {
			fmt.Println("fuck4")
			return false
		}
		_, err := os.Create(directoryPath)
		if err != nil {
			fmt.Println("create", err)
			return false
		}
		err = os.Remove(directoryPath)
		if err != nil {
			return false
		}
	} else {
		// Check if the file has write permissions
		if fileInfo.Mode().Perm()&0222 == 0 {
			return false
		}
	}
	return true
}
