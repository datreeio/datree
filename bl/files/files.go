package files

import (
	"os"
	"path/filepath"
	"sync"
)

func toAbsolutePath(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil && !fileInfo.IsDir() {
		return filepath.Abs(path)
	}

	return "", err
}

func ToAbsolutePaths(paths []string) (<-chan string, <-chan error) {
	errorChan := make(chan error, 100)
	pathsChan := make(chan string, 100)
	
	conc := 10
	wg := sync.WaitGroup{}
	wg.Add(conc)

	go func() {
		for i := 0; i < conc; i++ {
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

			wg.Done()
			}()
		}

		wg.Wait()
		close(pathsChan)
		close(errorChan)

	}()

	return pathsChan, errorChan
}
