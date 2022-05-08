package validate_yaml

import (
	"bytes"
	"errors"
	"testing"

	pkgExtractor "github.com/datreeio/datree/pkg/extractor"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileReader struct {
	mock.Mock
}

func (fr *MockFileReader) FilterFiles(paths []string) ([]string, error) {
	fr.Called(paths)
	return paths, nil
}

type MockPrinter struct {
	mock.Mock
}

func (p *MockPrinter) PrintFilename(title string) {
	p.Called(title)
}

func (p *MockPrinter) PrintYamlValidationErrors(validationErrors []error) {
	p.Called(validationErrors)
}

func (p *MockPrinter) PrintYamlValidationSummary(passedFiles int, allFiles int) {
	p.Called(passedFiles, allFiles)
}

func (p *MockPrinter) PrintMessage(messageText string, messageColor string) {
	p.Called(messageText, messageColor)
}

type MockExtractor struct {
	mock.Mock
}

func (e *MockExtractor) ExtractConfigurationsFromYamlFile(path string) (*[]pkgExtractor.Configuration, string, *pkgExtractor.InvalidFile) {
	args := e.Called(path)
	var firstArgument *[]pkgExtractor.Configuration
	var thirdArgument *pkgExtractor.InvalidFile

	if args.Get(0) == nil {
		firstArgument = nil
	} else {
		firstArgument = args.Get(0).(*[]pkgExtractor.Configuration)
	}

	if args.Get(2) == nil {
		thirdArgument = nil
	} else {
		thirdArgument = args.Get(2).(*pkgExtractor.InvalidFile)
	}

	return firstArgument, args.String(1), thirdArgument
}

func createMocks() (*MockFileReader, *MockPrinter, *MockExtractor) {
	mockedFileReader := &MockFileReader{}
	mockedFileReader.On("FilterFiles", mock.Anything).Return([]string{"."}, nil)

	mockedPrinter := &MockPrinter{}
	mockedPrinter.On("PrintFilename", mock.Anything).Return()
	mockedPrinter.On("PrintYamlValidationErrors", mock.Anything).Return()
	mockedPrinter.On("PrintYamlValidationSummary", mock.Anything, mock.Anything).Return()
	mockedPrinter.On("PrintMessage", mock.Anything, mock.Anything).Return()

	mockedExtractor := &MockExtractor{}
	mockedExtractor.On("ExtractConfigurationsFromYamlFile", mock.Anything).Return(nil, "", nil)

	return mockedFileReader, mockedPrinter, mockedExtractor
}

func createCommand(reader IReader, printer IPrinter, extractor IExtractor) *cobra.Command {
	cmd := New(&ValidateYamlCommandContext{
		Reader:    reader,
		Printer:   printer,
		Extractor: extractor,
	})
	return cmd
}

func TestNoFilesSelected(t *testing.T) {
	mockedFileReader, mockedPrinter, mockedExtractor := createMocks()
	cmd := createCommand(mockedFileReader, mockedPrinter, mockedExtractor)

	cmd.SetArgs([]string{})
	output := cmd.Execute()

	assert.Equal(t, errors.New("No files detected"), output, "expected No files detected error")
}

func TestInvalidYaml(t *testing.T) {
	mockedExtractor := &MockExtractor{}
	mockedExtractor.On("ExtractConfigurationsFromYamlFile", mock.Anything).Return(nil, "", &pkgExtractor.InvalidFile{
		Path: ".gitignore",
		ValidationErrors: []error{
			&pkgExtractor.InvalidYamlError{
				ErrorMessage: "file content is not valid yaml",
			},
		},
	})

	mockedFileReader, mockedPrinter, _ := createMocks()
	cmd := createCommand(mockedFileReader, mockedPrinter, mockedExtractor)

	actual := new(bytes.Buffer)
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs([]string{"."})
	output := cmd.Execute()

	assert.Equal(t, YamlNotValidError, output, "expected YamlNotValidError error")
}

func TestValidYaml(t *testing.T) {
	// mockedExtractor := &MockExtractor{}
	// mockedExtractor.On("ExtractConfigurationsFromYamlFile", mock.Anything).Return(nil, "", nil)

	mockedFileReader, mockedPrinter, mockedExtractor := createMocks()
	cmd := createCommand(mockedFileReader, mockedPrinter, mockedExtractor)

	actual := new(bytes.Buffer)
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs([]string{"."})
	output := cmd.Execute()

	assert.Equal(t, nil, output, "expected nil")
}
