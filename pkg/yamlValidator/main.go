package yamlValidator

import pkgExtractor "github.com/datreeio/datree/pkg/extractor"

type IExtractor interface {
	ExtractConfigurationsFromYamlFile(path string) (*[]pkgExtractor.Configuration, string, *pkgExtractor.InvalidFile)
}

type IPrinter interface {
	PrintFilename(title string)
	PrintYamlValidationErrors(validationErrors []error)
	PrintYamlValidationSummary(passedFiles int, allFiles int)
	PrintMessage(messageText string, messageColor string)
}

func ValidateFiles(extractor IExtractor, filesPaths []string) []*pkgExtractor.InvalidFile {
	var invalidYamlFiles []*pkgExtractor.InvalidFile
	for _, filePath := range filesPaths {
		_, _, invalidYamlFile := extractor.ExtractConfigurationsFromYamlFile(filePath)
		if invalidYamlFile != nil {
			invalidYamlFiles = append(invalidYamlFiles, invalidYamlFile)
		}
	}

	return invalidYamlFiles
}

func PrintValidationResults(printer IPrinter, invalidFiles []*pkgExtractor.InvalidFile, filesCount int) {
	for _, invalidFile := range invalidFiles {
		printer.PrintFilename(invalidFile.Path)
		printer.PrintYamlValidationErrors(invalidFile.ValidationErrors)
	}

	// print summary
	validFilesCount := filesCount - len(invalidFiles)
	printer.PrintYamlValidationSummary(validFilesCount, filesCount)
}
