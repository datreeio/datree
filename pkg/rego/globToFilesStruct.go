package rego

import (
	"github.com/datreeio/datree/pkg/extractor"
	"path/filepath"
)

type FilesAsStruct map[string]string

func GlobToFilesStruct(globPattern string) FilesAsStruct {
	paths, err := filepath.Glob(globPattern)

	if err != nil {
		panic(err)
	}

	files := make(map[string]string)

	for _, path := range paths {
		fileContent, err := extractor.ReadFileContent(path)
		if err != nil {
			panic(err)
		}
		files[path] = fileContent
	}

	return files
}
