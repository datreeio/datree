package files

import (
	"fmt"
	"os"
	"path/filepath"
)

func toAbsolutePath(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	fileInfo, _ := os.Stat(absolutePath)
	if fileInfo != nil && !fileInfo.IsDir() {
		return filepath.Abs(absolutePath)
	}

	return "", fmt.Errorf("failed parsing absolute path %s", path)
}

func ToAbsolutePaths(paths []string) chan string {
	pathsChan := make(chan string)

	go func() {
		for _, p := range paths {
			absolutePath, err := toAbsolutePath(p)
			if err == nil {
				pathsChan <- absolutePath
			}
		}
		close(pathsChan)
	}()

	return pathsChan
}
