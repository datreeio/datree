package files

import (
	"os"
	"path/filepath"
)

func toAbsolutePath(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil && fileInfo != nil && !fileInfo.IsDir() {
		return filepath.Abs(path)
	}

	return "", err
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
