package test

import (
	"fmt"

	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/extractor"
)

type ValidationManager struct {
	invalidYamlFiles                 []*extractor.InvalidFile
	invalidK8sFiles                  []*extractor.InvalidFile
	validK8sFilesConfigurations      []*extractor.FileConfigurations
	k8sValidationWarningPerValidFile *validation.K8sValidationWarningPerValidFile
	ignoredFiles                     []extractor.FileConfigurations
}

func (v *ValidationManager) AggregateInvalidYamlFiles(invalidFilesChan chan *extractor.InvalidFile) {
	for invalidFile := range invalidFilesChan {
		v.invalidYamlFiles = append(v.invalidYamlFiles, invalidFile)
	}
}

func (v *ValidationManager) InvalidYamlFiles() []*extractor.InvalidFile {
	return v.invalidYamlFiles
}

func (v *ValidationManager) InvalidYamlFilesCount() int {
	return len(v.invalidYamlFiles)
}

func (v *ValidationManager) AggregateInvalidK8sFiles(invalidFilesChan chan *extractor.InvalidFile) {
	for invalidFile := range invalidFilesChan {
		v.invalidK8sFiles = append(v.invalidK8sFiles, invalidFile)
	}
}

func (v *ValidationManager) InvalidK8sFiles() []*extractor.InvalidFile {
	return v.invalidK8sFiles
}

func (v *ValidationManager) InvalidK8sFilesCount() int {
	return len(v.invalidK8sFiles)
}

func (v *ValidationManager) AggregateValidK8sFiles(validK8sFilesConfigurationsChan chan *extractor.FileConfigurations) {
	for fileConfigurations := range validK8sFilesConfigurationsChan {
		v.validK8sFilesConfigurations = append(v.validK8sFilesConfigurations, fileConfigurations)
	}
}

func (v *ValidationManager) ValidK8sFilesConfigurations() []*extractor.FileConfigurations {
	return v.validK8sFilesConfigurations
}

func (v *ValidationManager) GetK8sValidationSummaryStr(filesCount int) string {
	if v.hasFilesWithWarnings() {
		return "skipped since there is no internet connection"
	}

	return fmt.Sprintf("%v/%v", v.ValidK8sFilesConfigurationsCount(), filesCount)
}

func (v *ValidationManager) hasFilesWithWarnings() bool {
	for _, value := range *v.k8sValidationWarningPerValidFile {
		if value != nil {
			return true
		}
	}

	return false
}

func (v *ValidationManager) ValidK8sFilesConfigurationsCount() int {
	return len(v.validK8sFilesConfigurations)
}

func (v *ValidationManager) ValidK8sConfigurationsCount() int {
	totalConfigs := 0

	for _, fileConfiguration := range v.validK8sFilesConfigurations {
		totalConfigs += len(fileConfiguration.Configurations)
	}

	return totalConfigs
}

func (v *ValidationManager) SaveK8sValidationWarningPerValidFile(k8sValidationWarningPerValidFile *validation.K8sValidationWarningPerValidFile) {
	v.k8sValidationWarningPerValidFile = k8sValidationWarningPerValidFile
}

func (v *ValidationManager) GetK8sValidationWarningPerValidFile() *validation.K8sValidationWarningPerValidFile {
	return v.k8sValidationWarningPerValidFile
}

func (v *ValidationManager) AggregateIgnoredYamlFiles(ignoredFilesChan chan *extractor.FileConfigurations) {
	for ignoredFile := range ignoredFilesChan {
		v.ignoredFiles = append(v.ignoredFiles, *ignoredFile)
	}
}

func (v *ValidationManager) IgnoredFiles() []*extractor.FileConfigurations {
	return v.validK8sFilesConfigurations
}

func (v *ValidationManager) IgnoredFilesCount() int {
	return len(v.ignoredFiles)
}
