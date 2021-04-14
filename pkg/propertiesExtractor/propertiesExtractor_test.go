package propertiesExtractor

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockFileReader struct {
	mock.Mock
}

func (m *mockFileReader) ReadFileContent(filepath string) (string, error) {
	args := m.Called(filepath)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockFileReader) GetFilesPaths(pattern string) (chan string, chan error) {
	args := m.Called(pattern)
	return args.Get(0).(chan string), args.Get(1).(chan error)
}

type newPropertiesExtractorTestCase struct {
	name     string
	arg      FileReader
	expected FileReader
}

func TestNewPropertiesExtractor(t *testing.T) {
	expectedFileReader := &fileReader.FileReader{}

	tests := []newPropertiesExtractorTestCase{
		{
			name:     "override with argument",
			arg:      expectedFileReader,
			expected: expectedFileReader,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			propertiesExtractor := NewPropertiesExtractor(tt.arg)
			expectedGlobFnValue := reflect.ValueOf(tt.expected)
			actualGlobFnValue := reflect.ValueOf(propertiesExtractor.reader)

			assert.Equal(t, expectedGlobFnValue.Pointer(), actualGlobFnValue.Pointer())
		})
	}
}

type readFilesFromPatternTestCase struct {
	name string
	args struct {
		conc    int
		pattern string
	}
	mock struct {
		readFileContentResponse struct {
			data   string
			errors error
		}
		getFilesPathsResponse struct {
			data   chan string
			errors chan error
		}
	}
	expected struct {
		readFileContentCalledWith []string
		filesProperties           []*FileProperties
		fileErrors                []FileError
		errors                    []error
	}
}

func TestReadFilesFromPattern(t *testing.T) {
	tests := []readFilesFromPatternTestCase{
		test_readFileFromPattern_success(),
		test_readFileFromPattern_invalidYaml(),
		test_readFileFromPattern_invalidFile(),
	}

	fileReader := &mockFileReader{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileReader.On("GetFilesPaths", mock.Anything).Return(tt.mock.getFilesPathsResponse.data, tt.mock.getFilesPathsResponse.errors)
			fileReader.On("ReadFileContent", mock.Anything).Return(tt.mock.readFileContentResponse.data, tt.mock.readFileContentResponse.errors)

			propertiesExtractor := &PropertiesExtractor{
				reader: fileReader,
			}
			actualFilesProperties, actualFileErrs, actualErrors := propertiesExtractor.ReadFilesFromPattern(tt.args.pattern, tt.args.conc)

			fileReader.AssertCalled(t, "GetFilesPaths", tt.args.pattern)
			fileReader.AssertCalled(t, "ReadFileContent", tt.expected.readFileContentCalledWith[0])

			for i, prop := range actualFilesProperties {
				assert.Equal(t, tt.expected.filesProperties[i], prop)
			}

			for i, err := range actualFileErrs {
				assert.Equal(t, tt.expected.fileErrors[i], err)
			}

			for i, err := range actualErrors {
				assert.Equal(t, tt.expected.errors[i], err)
			}
		})
	}
}

func test_readFileFromPattern_success() readFilesFromPatternTestCase {
	return readFilesFromPatternTestCase{
		name: "success",
		args: struct {
			conc    int
			pattern string
		}{
			conc:    5,
			pattern: "pattern/*",
		},
		mock: struct {
			readFileContentResponse struct {
				data   string
				errors error
			}
			getFilesPathsResponse struct {
				data   chan string
				errors chan error
			}
		}{
			readFileContentResponse: struct {
				data   string
				errors error
			}{
				data:   "apiVersion: extensions/v1beta1",
				errors: nil,
			},
			getFilesPathsResponse: struct {
				data   chan string
				errors chan error
			}{
				data:   mock_createFilesPathsChannel("path1/path2/file.yaml"),
				errors: mock_createErrorsChannel(nil),
			},
		},
		expected: struct {
			readFileContentCalledWith []string
			filesProperties           []*FileProperties
			fileErrors                []FileError
			errors                    []error
		}{
			filesProperties:           []*FileProperties{{FileName: "path1/path2/file.yaml", Configurations: []K8sConfiguration{{"apiVersion": "extensions/v1beta1"}}}},
			fileErrors:                []FileError{},
			errors:                    []error{},
			readFileContentCalledWith: []string{"path1/path2/file.yaml"},
		},
	}
}

func test_readFileFromPattern_invalidYaml() readFilesFromPatternTestCase {
	return readFilesFromPatternTestCase{
		name: "invalid yaml",
		args: struct {
			conc    int
			pattern string
		}{
			conc:    5,
			pattern: "pattern/*",
		},
		mock: struct {
			readFileContentResponse struct {
				data   string
				errors error
			}
			getFilesPathsResponse struct {
				data   chan string
				errors chan error
			}
		}{
			readFileContentResponse: struct {
				data   string
				errors error
			}{
				data:   "invalid data",
				errors: nil,
			},
			getFilesPathsResponse: struct {
				data   chan string
				errors chan error
			}{
				data:   mock_createFilesPathsChannel("path1/path2/file.yaml"),
				errors: mock_createErrorsChannel(fmt.Errorf("unmarshal yaml")),
			},
		},
		expected: struct {
			readFileContentCalledWith []string
			filesProperties           []*FileProperties
			fileErrors                []FileError
			errors                    []error
		}{
			filesProperties:           []*FileProperties{},
			fileErrors:                []FileError{},
			errors:                    []error{fmt.Errorf("unmarshal yaml")},
			readFileContentCalledWith: []string{"path1/path2/file.yaml"},
		},
	}
}

func test_readFileFromPattern_invalidFile() readFilesFromPatternTestCase {
	return readFilesFromPatternTestCase{
		name: "invalid yaml",
		args: struct {
			conc    int
			pattern string
		}{
			conc:    5,
			pattern: "pattern/*",
		},
		mock: struct {
			readFileContentResponse struct {
				data   string
				errors error
			}
			getFilesPathsResponse struct {
				data   chan string
				errors chan error
			}
		}{
			readFileContentResponse: struct {
				data   string
				errors error
			}{
				data:   "invalid data",
				errors: fmt.Errorf("invalid file"),
			},
			getFilesPathsResponse: struct {
				data   chan string
				errors chan error
			}{
				data:   mock_createFilesPathsChannel("path1/path2/file.yaml"),
				errors: mock_createErrorsChannel(fmt.Errorf("unmarshal yaml")),
			},
		},
		expected: struct {
			readFileContentCalledWith []string
			filesProperties           []*FileProperties
			fileErrors                []FileError
			errors                    []error
		}{
			filesProperties:           []*FileProperties{},
			fileErrors:                []FileError{{Message: "invalid file", Filename: "path1/path2/file.yaml"}},
			errors:                    []error{fmt.Errorf("unmarshal yaml")},
			readFileContentCalledWith: []string{"path1/path2/file.yaml"},
		},
	}
}

func mock_createFilesPathsChannel(filepath string) chan string {
	pathsChan := make(chan string, 1)
	pathsChan <- filepath
	close(pathsChan)

	return pathsChan
}

func mock_createErrorsChannel(err error) chan error {
	errorsChan := make(chan error, 1)

	if err != nil {
		errorsChan <- err
	}

	close(errorsChan)

	return errorsChan
}
