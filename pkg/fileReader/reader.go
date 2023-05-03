package fileReader

import (
	"os"
	"path/filepath"
	"regexp"

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
		readFile: os.ReadFile,
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

func (fr *FileReader) FilterFiles(paths []string, excludePattern string) ([]string, error) {
	var filePaths []string

	excludeRe, err := regexp.Compile(excludePattern)
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		stat, err := fr.stat(path)
		if err != nil {
			return []string{}, err
		}

		isMatched := excludePattern != "" && excludeRe.MatchString(path)
		if !stat.IsDir() && !isMatched {
			filePaths = append(filePaths, path)
		}
	}

	return filePaths, nil
}

func (fr *FileReader) ReadFileContent(filepath string) (string, error) {
	dat, err := fr.readFile(filepath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

func (fr *FileReader) ReadDir(dirpath string) ([]string, error) {
	stat, err := fr.stat(dirpath)
	if err != nil {
		return []string{}, err
	}
	if !stat.IsDir() {
		return []string{}, nil
	}

	files, err := fr.glob(dirpath + "/**/*")
	if err != nil {
		return []string{}, err
	}

	return files, nil
}

func (fr *FileReader) GetFilename(path string) (string, error) {
	stat, err := fr.stat(path)
	if err != nil {
		return "", err
	}

	if stat.IsDir() {
		return "", nil
	}

	return stat.Name(), nil
}
