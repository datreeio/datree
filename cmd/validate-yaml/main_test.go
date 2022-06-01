package validate_yaml

import (
	"bytes"
	"errors"
	"testing"

	"github.com/datreeio/datree/pkg/cliClient"
	pkgExtractor "github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/yamlValidator"
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

func (p *MockPrinter) GetFileNameText(title string) string {
	p.Called(title)
	return ""
}

func (p *MockPrinter) GetYamlValidationErrorsText(validationErrors []error) string {
	p.Called(validationErrors)
	return ""
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
	var configurations *[]pkgExtractor.Configuration
	var invalidFile *pkgExtractor.InvalidFile

	if args.Get(0) == nil {
		configurations = nil
	} else {
		configurations = args.Get(0).(*[]pkgExtractor.Configuration)
	}

	if args.Get(2) == nil {
		invalidFile = nil
	} else {
		invalidFile = args.Get(2).(*pkgExtractor.InvalidFile)
	}

	return configurations, args.String(1), invalidFile
}

type MockCliClient struct {
	mock.Mock
}

func (cc *MockCliClient) SendValidateYamlResult(request *cliClient.ValidatedYamlResult) {
	cc.Called(request)
}

type MockLocalConfig struct {
	mock.Mock
}

func (lc *MockLocalConfig) GetLocalConfiguration() (*localConfig.LocalConfig, error) {
	lc.Called()
	localConfig := &localConfig.LocalConfig{}
	return localConfig, nil
}

func createMocks() (*MockFileReader, *MockPrinter, *MockExtractor, *MockCliClient, *MockLocalConfig) {
	mockedFileReader := &MockFileReader{}
	mockedFileReader.On("FilterFiles", mock.Anything).Return([]string{"."}, nil)

	mockedPrinter := &MockPrinter{}
	mockedPrinter.On("GetFileNameText", mock.Anything).Return()
	mockedPrinter.On("GetYamlValidationErrorsText", mock.Anything).Return()
	mockedPrinter.On("PrintYamlValidationSummary", mock.Anything, mock.Anything).Return()
	mockedPrinter.On("PrintMessage", mock.Anything, mock.Anything).Return()

	mockedExtractor := &MockExtractor{}
	mockedExtractor.On("ExtractConfigurationsFromYamlFile", mock.Anything).Return(nil, "", nil)

	mockedCliClient := &MockCliClient{}
	mockedCliClient.On("SendValidateYamlResult", mock.Anything).Return()

	mockedLocalConfig := &MockLocalConfig{}
	mockedLocalConfig.On("GetLocalConfiguration").Return(&localConfig.LocalConfig{}, nil)

	return mockedFileReader, mockedPrinter, mockedExtractor, mockedCliClient, mockedLocalConfig
}

func createCommand(reader IReader, printer IPrinter, extractor yamlValidator.IExtractor, CliClient ICliClient, LocalConfig ILocalConfig, cliVersion string) *cobra.Command {
	cmd := New(&ValidateYamlCommandContext{
		Reader:      reader,
		Printer:     printer,
		Extractor:   extractor,
		CliClient:   CliClient,
		LocalConfig: LocalConfig,
		CliVersion:  cliVersion,
	})
	return cmd
}

func TestNoFilesSelected(t *testing.T) {
	mockedFileReader, mockedPrinter, mockedExtractor, mockedCliClient, mockedLocalConfig := createMocks()
	cmd := createCommand(mockedFileReader, mockedPrinter, mockedExtractor, mockedCliClient, mockedLocalConfig, "")

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

	mockedFileReader, mockedPrinter, _, mockedCliClient, mockedLocalConfig := createMocks()
	cmd := createCommand(mockedFileReader, mockedPrinter, mockedExtractor, mockedCliClient, mockedLocalConfig, "")

	actual := new(bytes.Buffer)
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs([]string{"."})
	output := cmd.Execute()

	assert.Equal(t, YamlNotValidError, output, "expected YamlNotValidError error")
}

func TestValidYaml(t *testing.T) {
	mockedFileReader, mockedPrinter, mockedExtractor, mockedCliClient, mockedLocalConfig := createMocks()
	cmd := createCommand(mockedFileReader, mockedPrinter, mockedExtractor, mockedCliClient, mockedLocalConfig, "")

	actual := new(bytes.Buffer)
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs([]string{"."})
	output := cmd.Execute()

	assert.Equal(t, nil, output, "expected nil")
}
