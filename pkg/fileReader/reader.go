package fileReader

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v2"
)

type ReadFileFn = func(filename string) ([]byte, error)
type GlobFn = func(pattern string) ([]string, error)
type StatFn = func(name string) (os.FileInfo, error)
type AbsFn = func(path string) (string, error)

type FileReader struct {
	readFile ReadFileFn
	glob     GlobFn
	stat     StatFn
	abs      AbsFn
}

type FileReaderOptions struct {
	ReadFile ReadFileFn
	Glob     GlobFn
	Stat     StatFn
	Abs      AbsFn
}

func CreateFileReader(opts *FileReaderOptions) *FileReader {
	fileReader := &FileReader{
		readFile: ioutil.ReadFile,
		glob:     doublestar.Glob,
		abs:      filepath.Abs,
		stat:     os.Stat,
	}

	if opts != nil {
		if opts.ReadFile != nil {
			fileReader.readFile = opts.ReadFile
		}

		if opts.Glob != nil {
			fileReader.glob = opts.Glob
		}

		if opts.Stat != nil {
			fileReader.stat = opts.Stat
		}

		if opts.Abs != nil {
			fileReader.abs = opts.Abs
		}
	}

	return fileReader
}

func (fr *FileReader) FilterFiles(paths []string) []string {
	var filePaths []string

	for _, path := range paths {
		stat, _ := fr.stat(path)

		if !stat.IsDir() {
			filePaths = append(filePaths, path)
		}
	}

	return filePaths
}

func (fr *FileReader) ReadFileContent(filepath string) (string, error) {
	dat, err := fr.readFile(filepath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}
