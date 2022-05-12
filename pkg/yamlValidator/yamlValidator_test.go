package yamlValidator

import (
	"testing"

	pkgExtractor "github.com/datreeio/datree/pkg/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestValidateFiles(t *testing.T) {
	extractor := &MockExtractor{}

	var yamlValidator = New(&YamlValidatorOptions{
		Extractor: extractor,
	})

	t.Run("invalid yaml files", func(t *testing.T) {
		invalidFilePath := "invalid.not-yaml"
		extractor.On("ExtractConfigurationsFromYamlFile", mock.Anything).Once().Return(nil, "", &pkgExtractor.InvalidFile{Path: invalidFilePath})
		invalidFiles := yamlValidator.ValidateFiles([]string{invalidFilePath})
		assert.Equal(t, invalidFilePath, invalidFiles[0].Path, "invalid file path should be returned")
	})

	t.Run("valid yaml files", func(t *testing.T) {
		validFilePath := "valid.yaml"
		extractor.On("ExtractConfigurationsFromYamlFile", mock.Anything).Once().Return(nil, "", nil)
		invalidFiles := yamlValidator.ValidateFiles([]string{validFilePath})
		assert.Equal(t, 0, len(invalidFiles), "invalid files should be empty")
	})
}
