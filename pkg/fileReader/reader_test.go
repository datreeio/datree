package fileReader

import (
	"io/fs"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ioMock struct {
	mock.Mock
}

func (c *ioMock) ReadFile(filename string) ([]byte, error) {
	args := c.Called(filename)
	return args.Get(0).([]byte), args.Error(1)
}

type readFileContentTestCase struct {
	name string
	args struct {
		path string
	}
	mock struct {
		readFilFn struct {
			response []byte
			err      error
		}
	}
	expected struct {
		calledWith string
		isCalled   bool
		response   string
		err        error
	}
}

func TestReadFileContent(t *testing.T) {
	io := ioMock{}

	tests := []readFileContentTestCase{
		{
			name: "success - should override with opts.methods and return response",
			args: struct{ path string }{
				path: "path/file.yaml",
			},
			mock: struct {
				readFilFn struct {
					response []byte
					err      error
				}
			}{
				readFilFn: struct {
					response []byte
					err      error
				}{
					response: []byte("content in file"),
					err:      nil,
				},
			},
			expected: struct {
				calledWith string
				isCalled   bool
				response   string
				err        error
			}{
				calledWith: "path/file.yaml",
				isCalled:   true,
				response:   "content in file",
				err:        nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io.On("ReadFile", mock.Anything).Return(tt.mock.readFilFn.response, tt.mock.readFilFn.err)
			fileReader := &FileReader{
				glob:     nil,
				readFile: io.ReadFile,
			}

			res, err := fileReader.ReadFileContent(tt.args.path)
			io.AssertCalled(t, "ReadFile", tt.expected.calledWith)
			assert.Equal(t, tt.expected.response, res)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}

type globMock struct {
	mock.Mock
}

func (c *globMock) Glob(pattern string) (matches []string, err error) {
	args := c.Called(pattern)
	return args.Get(0).([]string), args.Error(1)
}

type statMock struct {
	mock.Mock
}

func (c *statMock) Stat(name string) (os.FileInfo, error) {
	args := c.Called(name)
	return args.Get(0).(os.FileInfo), args.Error(1)
}

type absMock struct {
	mock.Mock
}

func (c *absMock) Abs(path string) (string, error) {
	args := c.Called(path)
	return args.Get(0).(string), args.Error(1)
}

type filterFilesTestCase struct {
	name string
	args struct {
		paths          []string
		excludePattern string
	}
	mock struct {
		stat struct {
			response fs.FileInfo
			err      error
		}
	}
	expected struct {
		filePaths []string
	}
}

func TestFilterFiles(t *testing.T) {
	stat := statMock{}
	fileInfo := &MockFileInfo{}

	tests := []filterFilesTestCase{
		{
			name: "success",
			args: struct {
				paths          []string
				excludePattern string
			}{
				paths:          []string{"file1.yaml", "file2.yaml", "file3-exclude.yaml", "file4.yaml", "file5-exclude.yaml"},
				excludePattern: "exclude.yaml",
			},
			mock: struct {
				stat struct {
					response fs.FileInfo
					err      error
				}
			}{
				stat: struct {
					response fs.FileInfo
					err      error
				}{
					response: fileInfo,
					err:      nil,
				},
			},
			expected: struct{ filePaths []string }{
				filePaths: []string{"file1.yaml", "file2.yaml", "file4.yaml"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stat.On("Stat", mock.Anything).Return(tt.mock.stat.response, tt.mock.stat.err)
			fileInfo.On("IsDir", mock.Anything).Return(false)
			fileReader := &FileReader{
				glob: nil,
				stat: stat.Stat,
			}

			filteredfiles, _ := fileReader.FilterFiles(tt.args.paths, tt.args.excludePattern)
			stat.AssertCalled(t, "Stat", tt.expected.filePaths[0])
			fileInfo.AssertCalled(t, "IsDir")
			assert.Equal(t, tt.expected.filePaths, filteredfiles)
		})
	}
}

func TestCreateFileReader(t *testing.T) {
	glob := globMock{}
	io := ioMock{}
	stat := statMock{}
	abs := absMock{}

	opt := &FileReaderOptions{
		Glob:     glob.Glob,
		ReadFile: io.ReadFile,
		Stat:     stat.Stat,
		Abs:      abs.Abs,
	}

	fileReader := CreateFileReader(opt)

	expectedGlobFnValue := reflect.ValueOf(glob.Glob)
	actualGlobFnValue := reflect.ValueOf(fileReader.glob)

	assert.Equal(t, expectedGlobFnValue.Pointer(), actualGlobFnValue.Pointer())
}

type getFilesPathsTestCase struct {
	name string
	args struct {
		paths []string
	}
	mock struct {
		stat struct {
			response os.FileInfo
			err      error
		}
		abs struct {
			response string
			err      error
		}
	}
	expected struct {
		stat struct {
			calledWith string
			isCalled   bool
		}
		abs struct {
			calledWith string
			isCalled   bool
		}
		response []string
		err      []error
	}
}

func getFilesPaths_noMatchesTestCase() getFilesPathsTestCase {
	return getFilesPathsTestCase{
		name: "success no matches",
		args: struct{ paths []string }{
			paths: []string{},
		},
		mock: struct {
			stat struct {
				response os.FileInfo
				err      error
			}
			abs struct {
				response string
				err      error
			}
		}{
			stat: struct {
				response os.FileInfo
				err      error
			}{
				response: nil,
				err:      nil,
			},
			abs: struct {
				response string
				err      error
			}{
				response: "",
				err:      nil,
			},
		},
		expected: struct {
			stat struct {
				calledWith string
				isCalled   bool
			}
			abs struct {
				calledWith string
				isCalled   bool
			}
			response []string
			err      []error
		}{
			stat: struct {
				calledWith string
				isCalled   bool
			}{
				calledWith: "",
				isCalled:   false,
			},
			abs: struct {
				calledWith string
				isCalled   bool
			}{
				calledWith: "",
				isCalled:   false,
			},
			response: nil,
			err:      nil,
		},
	}
}

type MockFileInfo struct {
	mock.Mock
	FileName    string
	IsDirectory bool
}

func (mfi MockFileInfo) Name() string       { return mfi.FileName }
func (mfi MockFileInfo) Size() int64        { return int64(8) }
func (mfi MockFileInfo) Mode() os.FileMode  { return os.ModePerm }
func (mfi MockFileInfo) ModTime() time.Time { return time.Now() }
func (mfi MockFileInfo) Sys() interface{}   { return nil }

func (c *MockFileInfo) IsDir() bool {
	args := c.Called()
	return args.Get(0).(bool)
}

func getFilesPaths_withMatchesTestCase() getFilesPathsTestCase {
	mockFileInfo := &MockFileInfo{}
	mockFileInfo.On("IsDir").Return(false)

	return getFilesPathsTestCase{
		name: "success with matches",
		args: struct{ paths []string }{
			paths: []string{"./fail-30.yaml"},
		},
		mock: struct {
			stat struct {
				response os.FileInfo
				err      error
			}
			abs struct {
				response string
				err      error
			}
		}{
			stat: struct {
				response os.FileInfo
				err      error
			}{
				response: mockFileInfo,
				err:      nil,
			},
			abs: struct {
				response string
				err      error
			}{
				response: "path",
				err:      nil,
			},
		},
		expected: struct {
			stat struct {
				calledWith string
				isCalled   bool
			}
			abs struct {
				calledWith string
				isCalled   bool
			}
			response []string
			err      []error
		}{
			stat: struct {
				calledWith string
				isCalled   bool
			}{
				calledWith: "./fail-30.yaml",
				isCalled:   true,
			},
			abs: struct {
				calledWith string
				isCalled   bool
			}{
				calledWith: "./fail-30.yaml",
				isCalled:   true,
			},
			response: []string{"path"},
			err:      nil,
		},
	}
}
