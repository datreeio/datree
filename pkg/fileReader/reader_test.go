package fileReader

import (
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
		pattern string
	}
	mock struct {
		glob struct {
			response []string
			err      error
		}
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
		glob struct {
			calledWith string
			isCalled   bool
		}
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

func TestGetFilesPaths(t *testing.T) {
	tests := []getFilesPathsTestCase{
		getFilesPaths_noMatchesTestCase(),
		getFilesPaths_withMatchesTestCase(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			glob := globMock{}
			stat := statMock{}
			abs := absMock{}

			glob.On("Glob", mock.Anything).Return(tt.mock.glob.response, tt.mock.glob.err)
			abs.On("Abs", mock.Anything).Return(tt.mock.abs.response, tt.mock.abs.err)
			stat.On("Stat", mock.Anything).Return(tt.mock.stat.response, tt.mock.stat.err)

			fileReader := &FileReader{
				glob:     glob.Glob,
				readFile: nil,
				stat:     stat.Stat,
				abs:      abs.Abs,
			}

			res, err := fileReader.GetFilesPaths(tt.args.pattern)

			var actualRes []string
			for m := range res {
				actualRes = append(actualRes, m)
			}

			var actualErrs []error
			for e := range err {
				actualErrs = append(actualErrs, e)
			}

			assert.Equal(t, tt.expected.response, actualRes)
			assert.Equal(t, tt.expected.err, actualErrs)

			if tt.expected.glob.isCalled {
				glob.AssertCalled(t, "Glob", tt.expected.glob.calledWith)
			}

			if tt.expected.abs.isCalled {
				abs.AssertCalled(t, "Abs", tt.expected.abs.calledWith)
			}

			if tt.expected.stat.isCalled {
				stat.AssertCalled(t, "Stat", tt.expected.stat.calledWith)
			}

		})
	}
}

func getFilesPaths_noMatchesTestCase() getFilesPathsTestCase {
	return getFilesPathsTestCase{
		name: "success no matches",
		args: struct{ pattern string }{
			pattern: "*",
		},
		mock: struct {
			glob struct {
				response []string
				err      error
			}
			stat struct {
				response os.FileInfo
				err      error
			}
			abs struct {
				response string
				err      error
			}
		}{
			glob: struct {
				response []string
				err      error
			}{
				response: make([]string, 0),
				err:      nil,
			},
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
			glob struct {
				calledWith string
				isCalled   bool
			}
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
			glob: struct {
				calledWith string
				isCalled   bool
			}{
				calledWith: "*",
				isCalled:   true,
			},
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
		args: struct{ pattern string }{
			pattern: "./mock/*",
		},
		mock: struct {
			glob struct {
				response []string
				err      error
			}
			stat struct {
				response os.FileInfo
				err      error
			}
			abs struct {
				response string
				err      error
			}
		}{
			glob: struct {
				response []string
				err      error
			}{
				response: []string{"./fail-30.yaml"},
				err:      nil,
			},
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
			glob struct {
				calledWith string
				isCalled   bool
			}
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
			glob: struct {
				calledWith string
				isCalled   bool
			}{
				calledWith: "./mock/*",
				isCalled:   true,
			},
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
