package yamlValidator

import (
	"github.com/datreeio/datree/bl/files"
	pkgExtractor "github.com/datreeio/datree/pkg/extractor"
)

type YamlValidatorOptions struct {
	Extractor IExtractor
}
type YamlValidator struct {
	extractor IExtractor
}

type IExtractor interface {
	ExtractConfigurationsFromYamlFile(path string) (*[]pkgExtractor.Configuration, string, *pkgExtractor.InvalidFile)
}

func New(options *YamlValidatorOptions) *YamlValidator {

	if options != nil {
		return &YamlValidator{
			options.Extractor,
		}
	}

	return &YamlValidator{
		extractor: files.New(),
	}
}

func (yv *YamlValidator) ValidateFiles(filesPaths []string) []*pkgExtractor.InvalidFile {
	var invalidYamlFiles []*pkgExtractor.InvalidFile
	for _, filePath := range filesPaths {
		_, _, invalidYamlFile := yv.extractor.ExtractConfigurationsFromYamlFile(filePath)
		if invalidYamlFile != nil {
			invalidYamlFiles = append(invalidYamlFiles, invalidYamlFile)
		}
	}

	return invalidYamlFiles
}
