package test

import (
	"fmt"
	"sync"

	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/extractor"
)

type ValidationManager struct {
	invalidYamlFiles                 []*extractor.InvalidFile
	invalidK8sFiles                  []*extractor.InvalidFile
	validK8sFilesConfigurations      []*extractor.FileConfigurations
	k8sValidationWarningPerValidFile validation.K8sValidationWarningPerValidFile
	ignoredFiles                     []extractor.FileConfigurations
}

func NewValidationManager() *ValidationManager {
	return &ValidationManager{
		k8sValidationWarningPerValidFile: make(validation.K8sValidationWarningPerValidFile),
	}
}

func (v *ValidationManager) AggregateInvalidYamlFiles(invalidFilesChan chan *extractor.InvalidFile, wg *sync.WaitGroup) {
	for invalidFile := range invalidFilesChan {
		v.invalidYamlFiles = append(v.invalidYamlFiles, invalidFile)
	}
	wg.Done()
}

func (v *ValidationManager) InvalidYamlFiles() []*extractor.InvalidFile {
	return v.invalidYamlFiles
}

func (v *ValidationManager) InvalidYamlFilesCount() int {
	return len(v.invalidYamlFiles)
}

func (v *ValidationManager) AggregateInvalidK8sFiles(invalidFilesChan chan *extractor.InvalidFile, wg *sync.WaitGroup) {
	for invalidFile := range invalidFilesChan {
		v.invalidK8sFiles = append(v.invalidK8sFiles, invalidFile)
	}
	wg.Done()
}

func (v *ValidationManager) InvalidK8sFiles() []*extractor.InvalidFile {
	return v.invalidK8sFiles
}

func (v *ValidationManager) InvalidK8sFilesCount() int {
	return len(v.invalidK8sFiles)
}

func (v *ValidationManager) AggregateValidK8sFiles(validK8sFilesConfigurationsChan chan *extractor.FileConfigurations, wg *sync.WaitGroup) {
	for fileConfigurations := range validK8sFilesConfigurationsChan {
		v.validK8sFilesConfigurations = append(v.validK8sFilesConfigurations, fileConfigurations)
	}
	wg.Done()
}

func (v *ValidationManager) ValidK8sFilesConfigurations() []*extractor.FileConfigurations {
	return v.validK8sFilesConfigurations
}

func (v *ValidationManager) GetK8sValidationSummaryStr(filesCount int) string {
	if v.hasFilesWithWarningsOfKind(validation.NetworkError) {
		return "skipped since there is no internet connection"
	}

	return fmt.Sprintf("%v/%v", v.ValidK8sFilesConfigurationsCount()-v.countFilesWithWarningsOfKind(validation.Skipped), filesCount)
}

func (v *ValidationManager) hasFilesWithWarningsOfKind(warningKind validation.WarningKind) bool {
	for _, value := range v.k8sValidationWarningPerValidFile {
		if value.WarningKind == warningKind {
			return true
		}
	}
	return false
}

func (v *ValidationManager) countFilesWithWarningsOfKind(warningKind validation.WarningKind) int {
	count := 0
	for _, value := range v.k8sValidationWarningPerValidFile {
		if value.WarningKind == warningKind {
			count++
		}
	}
	return count
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

func (v *ValidationManager) AggregateK8sValidationWarningsPerValidFile(filesWithWarningsChan chan *validation.FileWithWarning, wg *sync.WaitGroup) {
	for fileWithWarning := range filesWithWarningsChan {
		if fileWithWarning != nil {
			v.k8sValidationWarningPerValidFile[fileWithWarning.Filename] = *fileWithWarning
		}
	}
	wg.Done()
}

func (v *ValidationManager) GetK8sValidationWarningPerValidFile() validation.K8sValidationWarningPerValidFile {
	return v.k8sValidationWarningPerValidFile
}

func (v *ValidationManager) AggregateIgnoredYamlFiles(ignoredFilesChan chan *extractor.FileConfigurations, wg *sync.WaitGroup) {
	for ignoredFile := range ignoredFilesChan {
		v.ignoredFiles = append(v.ignoredFiles, *ignoredFile)
	}
	wg.Done()
}

func (v *ValidationManager) IgnoredFiles() []*extractor.FileConfigurations {
	return v.validK8sFilesConfigurations
}

func (v *ValidationManager) IgnoredFilesCount() int {
	return len(v.ignoredFiles)
}
