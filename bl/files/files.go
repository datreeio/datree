package files

import (
	"fmt"
	"os"
	"path/filepath"
)

func toAbsolutePath(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", nil
	}

	fileInfo, err := os.Stat(absolutePath)
	if fileInfo != nil && !fileInfo.IsDir() {
		return filepath.Abs(absolutePath)
	}

	return "", fmt.Errorf("failed parsing absolute path %s", path)
}

func ToAbsolutePaths(paths []string) (<-chan string, <-chan error) {
	errorChan := make(chan error, 100)
	pathsChan := make(chan string, 100)

	go func() {
		for _, p := range paths {
			absolutePath, err := toAbsolutePath(p)
			if err != nil {
				errorChan <- err
				continue
			} else {
				pathsChan <- absolutePath
			}
		}

		close(pathsChan)
		close(errorChan)
	}()

	return pathsChan, errorChan
}
