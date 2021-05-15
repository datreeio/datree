package propertiesExtractor

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/datreeio/datree/pkg/fileReader"
	"gopkg.in/yaml.v3"
)

type FileReader interface {
	ReadFileContent(filepath string) (string, error)
	GetFilesPaths(paths []string) (chan string, chan error)
}

type PropertiesExtractor struct {
	reader FileReader
}

type FileError struct {
	Message  string
	Filename string
}

type K8sConfiguration map[string]interface{}

type FileProperties struct {
	FileName       string             `json:"fileName"`
	Configurations []K8sConfiguration `json:"configurations"`
}

func NewPropertiesExtractor(fr FileReader) *PropertiesExtractor {
	if fr != nil {
		return &PropertiesExtractor{
			reader: fr,
		}
	}
	return &PropertiesExtractor{
		reader: fileReader.CreateFileReader(nil),
	}
}

func (e *PropertiesExtractor) ReadFilesFromPaths(paths []string, conc int) ([]*FileProperties, []FileError, []error) {
	filePathsChan, getFilesPathsErrorsChan := e.reader.GetFilesPaths(paths)
	filesPropertiesChan, extractPropertiesErrorsChan := e.extractFilesProperties(filePathsChan, conc)

	var files []*FileProperties
	for fileProperties := range filesPropertiesChan {
		files = append(files, fileProperties)
	}

	fileErrors := make([]FileError, 0)
	for fileError := range extractPropertiesErrorsChan {
		fileErrors = append(fileErrors, fileError)
	}

	errors := make([]error, 0)
	for err := range getFilesPathsErrorsChan {
		errors = append(errors, fmt.Errorf("error on get file paths, err:%s", err.Error()))
	}

	return files, fileErrors, errors
}

func (e *PropertiesExtractor) extractFilesProperties(filePathsChan <-chan string, conc int) (chan *FileProperties, chan FileError) {
	filesPropertiesChan := make(chan *FileProperties, 100)
	errorsChan := make(chan FileError, 100)

	wg := sync.WaitGroup{}
	wg.Add(conc)

	go func() {
		for i := 0; i < conc; i++ {
			go func() {
			readFilesPathsChan:
				for {
					select {
					case path, ok := <-filePathsChan:
						if !ok {
							break readFilesPathsChan
						}
						content, err := e.reader.ReadFileContent(path)
						if err != nil {
							errorsChan <- FileError{Filename: path, Message: err.Error()}
							continue
						}

						configurations, err := e.yamlToK8sConfigurations(content)
						if err != nil {
							errorsChan <- FileError{Filename: path, Message: err.Error()}
							continue
						} else {
							filesPropertiesChan <- &FileProperties{
								Configurations: configurations,
								FileName:       path,
							}
						}
					}
				}
				wg.Done()
			}()
		}

		wg.Wait()
		close(filesPropertiesChan)
		close(errorsChan)
	}()

	return filesPropertiesChan, errorsChan
}

func (e *PropertiesExtractor) yamlToK8sConfigurations(content string) ([]K8sConfiguration, error) {
	var configurations []K8sConfiguration

	yamlDecoder := yaml.NewDecoder(bytes.NewReader([]byte(content)))

	var err error
	for {
		var doc K8sConfiguration
		err = yamlDecoder.Decode(&doc)
		if err != nil {
			break
		}
		configurations = append(configurations, doc)
	}

	if err == io.EOF {
		err = nil
	}

	return configurations, err
}
