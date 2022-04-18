package test

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/bl/messager"
	policy_factory "github.com/datreeio/datree/bl/policy"
	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/fileReader"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/pkg/errors"

	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEvaluator struct {
	mock.Mock
}

// struct type for testing Test command flow
type TestFlowTestCase struct {
	name string
	args struct {
		path []string
	}
	mock struct {
		ExtractFilesConfigurations struct {
			filesConfigurationsChan chan *extractor.FileConfigurations
			invalidFilesChan        chan *extractor.InvalidFile
		}
		ValidateResources struct {
			k8sFilesConfigurationsChan chan *extractor.FileConfigurations
			k8sInvalidFilesChan        chan *extractor.InvalidFile
			filesWithWarningsChan      chan *validation.FileWithWarning
		}
		Evaluate struct {
			policyCheckResultData evaluation.PolicyCheckResultData
			err                   error
		}
		SendEvaluationResult struct {
			sendEvaluationResultsResponse *cliClient.SendEvaluationResultsResponse
			err                           error
		}
	}
	expected struct {
		ValidateResources struct {
			filesConfigurationsChan chan *extractor.FileConfigurations
		}
		Evaluate struct {
			evaluationData evaluation.PolicyCheckData
		}
		SendEvaluationResult struct {
			evaluationRequestData evaluation.EvaluationRequestData
		}
		PrintWarnings     []printer.Warning
		EvaluationSummary printer.EvaluationSummary
		err               error
	}
}

func (m *mockEvaluator) Evaluate(evaluationData evaluation.PolicyCheckData) (evaluation.PolicyCheckResultData, error) {
	args := m.Called(evaluationData)
	return args.Get(0).(evaluation.PolicyCheckResultData), args.Error(1)
}

func (m *mockEvaluator) SendEvaluationResult(evaluationRequestData evaluation.EvaluationRequestData) (*cliClient.SendEvaluationResultsResponse, error) {
	args := m.Called(evaluationRequestData)
	return args.Get(0).(*cliClient.SendEvaluationResultsResponse), args.Error(1)
}

func (m *mockEvaluator) RequestEvaluationPrerunData(token string) (*cliClient.EvaluationPrerunDataResponse, int, error) {
	args := m.Called(token)
	return args.Get(0).(*cliClient.EvaluationPrerunDataResponse), args.Get(1).(int), args.Error(2)
}

type mockMessager struct {
	mock.Mock
}

func (m *mockMessager) LoadVersionMessages(cliVersion string) chan *messager.VersionMessage {
	messages := make(chan *messager.VersionMessage, 1)
	go func() {
		messages <- &messager.VersionMessage{
			CliVersion:   "1.2.3",
			MessageText:  "version message mock",
			MessageColor: "green"}
		close(messages)
	}()

	m.Called(cliVersion)
	return messages
}

func (m *mockMessager) HandleVersionMessage(messageChannel <-chan *messager.VersionMessage) {
	m.Called(messageChannel)
}

type FilesExtractorMock struct {
	mock.Mock
}

func (fe *FilesExtractorMock) ExtractFilesConfigurations(paths []string, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile) {
	args := fe.Called(paths, concurrency)
	return args.Get(0).(chan *extractor.FileConfigurations), args.Get(1).(chan *extractor.InvalidFile)
}

func (fe *FilesExtractorMock) ExtractYamlFileToUnknownStruct(path string) (files.UnknownStruct, error) {
	args := fe.Called(path)
	return args.Get(0).(files.UnknownStruct), args.Error(1)
}

type K8sValidatorMock struct {
	mock.Mock
}

func (kv *K8sValidatorMock) ValidateResources(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile, chan *validation.FileWithWarning) {
	args := kv.Called(filesConfigurationsChan, concurrency)
	return args.Get(0).(chan *extractor.FileConfigurations), args.Get(1).(chan *extractor.InvalidFile), args.Get(2).(chan *validation.FileWithWarning)
}

func (kv *K8sValidatorMock) GetK8sFiles(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.FileConfigurations) {
	args := kv.Called(filesConfigurationsChan, concurrency)
	return args.Get(0).(chan *extractor.FileConfigurations), args.Get(1).(chan *extractor.FileConfigurations)
}

func (kv *K8sValidatorMock) InitClient(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string) {
}

type PrinterMock struct {
	mock.Mock
}

func (p *PrinterMock) PrintWarnings(warnings []printer.Warning) {
	p.Called(warnings)
}

func (p *PrinterMock) PrintSummaryTable(summary printer.Summary) {
	p.Called(summary)
}

func (p *PrinterMock) PrintEvaluationSummary(evaluationSummary printer.EvaluationSummary, k8sVersion string) {
	p.Called(evaluationSummary)
}

func (p *PrinterMock) PrintMessage(messageText string, messageColor string) {
	p.Called(messageText, messageColor)
}

func (p *PrinterMock) PrintPromptMessage(promptMessage string) {
	p.Called(promptMessage)
}

func (p *PrinterMock) SetTheme(theme *printer.Theme) {
	p.Called(theme)
}

type CliClientMock struct {
	mock.Mock
}

func (c *CliClientMock) RequestEvaluationPrerunData(token string) (*cliClient.EvaluationPrerunDataResponse, error) {
	args := c.Called(token)
	return args.Get(0).(*cliClient.EvaluationPrerunDataResponse), nil
}

type ReaderMock struct {
	mock.Mock
}

func (rm *ReaderMock) FilterFiles(paths []string) ([]string, error) {
	args := rm.Called(paths)
	return args.Get(0).([]string), nil
}

type LocalConfigMock struct {
	mock.Mock
}

func (lc *LocalConfigMock) GetLocalConfiguration() (*localConfig.LocalConfig, error) {
	lc.Called()
	return &localConfig.LocalConfig{Token: "134kh"}, nil
}

var filesConfigurations []*extractor.FileConfigurations
var evaluationId int
var ctx *TestCommandContext
var testingPolicy policy_factory.Policy

// mock instances
var k8sValidatorMock *K8sValidatorMock
var mockedEvaluator *mockEvaluator
var localConfigMock *LocalConfigMock
var messagerMock *mockMessager
var readerMock *ReaderMock
var mockedCliClient *CliClientMock

func pathFromRoot(path string) string {
	_, filename, _, _ := runtime.Caller(0)
	path = filepath.FromSlash(path)
	filename = filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filename))), path)
	return filename
}

func TestTestFlow(t *testing.T) {
	tests := []*TestFlowTestCase{
		test_no_k8s_resources_found(),
		test_all_k8s_resources_tested(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluatorMock := &mockEvaluator{}
			readerMock := &ReaderMock{}
			messager := &mockMessager{}
			messager.On("LoadVersionMessages", mock.Anything)
			k8sValidatorMock := &K8sValidatorMock{}
			printerMock := &PrinterMock{}
			localConfigMock := &LocalConfigMock{}
			localConfigMock.On("GetLocalConfiguration").Return(&localConfig.LocalConfig{Token: "134kh"}, nil)
			filesExtractorMock := &FilesExtractorMock{}

			readerMock.On("FilterFiles", mock.Anything).Return(tt.args.path, nil)
			filesExtractorMock.On("ExtractFilesConfigurations", mock.Anything, 100).Return(tt.mock.ExtractFilesConfigurations.filesConfigurationsChan, tt.mock.ExtractFilesConfigurations.invalidFilesChan)
			k8sValidatorMock.On("ValidateResources", mock.Anything, 100).Return(tt.mock.ValidateResources.k8sFilesConfigurationsChan, tt.mock.ValidateResources.k8sInvalidFilesChan, tt.mock.ValidateResources.filesWithWarningsChan)
			k8sValidatorMock.On("InitClient", mock.Anything, mock.Anything, mock.Anything).Return()
			evaluatorMock.On("Evaluate", mock.Anything, mock.Anything, mock.Anything).Return(tt.mock.Evaluate.policyCheckResultData, tt.mock.Evaluate.err)
			evaluatorMock.On("SendEvaluationResult", mock.Anything).Return(tt.mock.SendEvaluationResult.sendEvaluationResultsResponse, tt.mock.SendEvaluationResult.err)

			printerMock.On("PrintWarnings", mock.Anything)
			printerMock.On("PrintSummaryTable", mock.Anything)
			printerMock.On("PrintEvaluationSummary", mock.Anything, mock.Anything)
			printerMock.On("PrintMessage", mock.Anything, mock.Anything)
			printerMock.On("PrintPromptMessage", mock.Anything)
			printerMock.On("SetTheme", mock.Anything)
			ctx := &TestCommandContext{
				K8sValidator:   k8sValidatorMock,
				Evaluator:      evaluatorMock,
				LocalConfig:    localConfigMock,
				Messager:       messager,
				Printer:        printerMock,
				Reader:         readerMock,
				FilesExtractor: filesExtractorMock,
			}

			err := Test(ctx, tt.args.path, &TestCommandData{K8sVersion: "1.18.0", Output: "", Policy: tt.expected.Evaluate.evaluationData.Policy, Token: "134kh"})
			if tt.expected.err != nil {
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.Equal(t, err, tt.expected.err)
			}
			readerMock.AssertCalled(t, "FilterFiles", tt.args.path)
			printerMock.AssertNotCalled(t, "SetTheme", mock.Anything)
			filesExtractorMock.AssertCalled(t, "ExtractFilesConfigurations", tt.args.path, 100)
			k8sValidatorMock.AssertCalled(t, "ValidateResources", tt.expected.ValidateResources.filesConfigurationsChan, 100)
			evaluatorMock.AssertCalled(t, "Evaluate", mock.MatchedBy(func(policyCheckData evaluation.PolicyCheckData) bool {
				expectedPolicyCheckData := tt.expected.Evaluate.evaluationData
				if len(policyCheckData.FilesConfigurations) == len(expectedPolicyCheckData.FilesConfigurations) {
					for index, validK8sFilesConfiguration := range policyCheckData.FilesConfigurations {
						if validK8sFilesConfiguration.FileName != expectedPolicyCheckData.FilesConfigurations[index].FileName {
							return false
						}
					}
				} else {
					return false
				}

				if (policyCheckData.IsInteractiveMode != expectedPolicyCheckData.IsInteractiveMode) ||
					(policyCheckData.PolicyName != expectedPolicyCheckData.PolicyName) {
					return false
				}
				return true
			}))
			evaluatorMock.AssertCalled(t, "SendEvaluationResult", mock.MatchedBy(func(evaluationRequestData evaluation.EvaluationRequestData) bool {
				expectedEvaluationRequestData := tt.expected.SendEvaluationResult.evaluationRequestData
				if expectedEvaluationRequestData.PolicyName != evaluationRequestData.PolicyName {
					return false
				}

				if len(evaluationRequestData.FilesData) == len(expectedEvaluationRequestData.FilesData) {
					for index, value := range evaluationRequestData.FilesData {
						if value.FilePath != expectedEvaluationRequestData.FilesData[index].FilePath {
							return false
						}

						if value.ConfigurationsCount != expectedEvaluationRequestData.FilesData[index].ConfigurationsCount {
							return false
						}
					}
				} else {
					return false
				}

				if len(evaluationRequestData.FailedK8sFiles) == len(expectedEvaluationRequestData.FailedK8sFiles) {
					for index, value := range evaluationRequestData.FailedK8sFiles {
						if value != expectedEvaluationRequestData.FailedK8sFiles[index] {
							return false
						}
					}
				} else {
					return false
				}
				return true
			}))
			printerMock.AssertCalled(t, "PrintWarnings", mock.MatchedBy(func(warnings []printer.Warning) bool {
				return len(warnings) == len(tt.expected.PrintWarnings)
			}))
			printerMock.AssertCalled(t, "PrintEvaluationSummary", mock.MatchedBy(func(evaluationSummary printer.EvaluationSummary) bool {
				expected := tt.expected.EvaluationSummary
				if (evaluationSummary.ConfigsCount == expected.ConfigsCount) && (evaluationSummary.RulesCount == expected.RulesCount) &&
					(evaluationSummary.FilesCount == expected.FilesCount) && (evaluationSummary.PassedYamlValidationCount == expected.PassedYamlValidationCount) &&
					(evaluationSummary.K8sValidation == expected.K8sValidation) && (evaluationSummary.PassedPolicyCheckCount == expected.PassedPolicyCheckCount) {
					return true
				}
				return false
			}), mock.Anything)
		})
	}
}

func test_all_k8s_resources_tested() *TestFlowTestCase {
	validPath := pathFromRoot("internal/fixtures/kube/pass-all.yaml")
	filesConfigurationsChan := newFilesConfigurationsChan(validPath)
	invalidFilesChan := make(chan *extractor.InvalidFile, 100)
	var validK8sFilesConfigurations []*extractor.FileConfigurations
	for fileConfigurations := range newFilesConfigurationsChan(validPath) {
		validK8sFilesConfigurations = append(validK8sFilesConfigurations, fileConfigurations)
	}
	preRunData := mockGetPreRunData()
	policy, _ := policy_factory.CreatePolicy(preRunData.PoliciesJson, "", preRunData.RegistrationURL)
	close(invalidFilesChan)

	return &TestFlowTestCase{
		name: "all valid k8s yaml files passed as path",
		args: struct {
			path []string
		}{
			path: []string{validPath},
		},
		mock: struct {
			ExtractFilesConfigurations struct {
				filesConfigurationsChan chan *extractor.FileConfigurations
				invalidFilesChan        chan *extractor.InvalidFile
			}
			ValidateResources struct {
				k8sFilesConfigurationsChan chan *extractor.FileConfigurations
				k8sInvalidFilesChan        chan *extractor.InvalidFile
				filesWithWarningsChan      chan *validation.FileWithWarning
			}
			Evaluate struct {
				policyCheckResultData evaluation.PolicyCheckResultData
				err                   error
			}
			SendEvaluationResult struct {
				sendEvaluationResultsResponse *cliClient.SendEvaluationResultsResponse
				err                           error
			}
		}{
			ExtractFilesConfigurations: struct {
				filesConfigurationsChan chan *extractor.FileConfigurations
				invalidFilesChan        chan *extractor.InvalidFile
			}{
				filesConfigurationsChan: filesConfigurationsChan,
				invalidFilesChan:        invalidFilesChan,
			},
			ValidateResources: struct {
				k8sFilesConfigurationsChan chan *extractor.FileConfigurations
				k8sInvalidFilesChan        chan *extractor.InvalidFile
				filesWithWarningsChan      chan *validation.FileWithWarning
			}{
				k8sFilesConfigurationsChan: newFilesConfigurationsChan(validPath),
				k8sInvalidFilesChan:        invalidFilesChan,
				filesWithWarningsChan:      newK8sValidationWarningsChan(),
			},
			Evaluate: struct {
				policyCheckResultData evaluation.PolicyCheckResultData
				err                   error
			}{
				policyCheckResultData: evaluation.PolicyCheckResultData{
					FormattedResults: evaluation.FormattedResults{
						EvaluationResults: &evaluation.EvaluationResults{
							FileNameRuleMapper: make(map[string]map[string]*evaluation.Rule),
							Summary: evaluation.EvaluationResultsSummery{
								TotalFailedRules:  0,
								TotalSkippedRules: 0,
								FilesCount:        1,
								FilesPassedCount:  1,
							},
						},
					},
					RulesData:  []cliClient.RuleData{},
					FilesData:  []cliClient.FileData{{FilePath: validPath, ConfigurationsCount: 5}},
					RawResults: evaluation.FailedRulesByFiles{},
					RulesCount: len(policy.Rules),
				},
				err: nil,
			},
			SendEvaluationResult: struct {
				sendEvaluationResultsResponse *cliClient.SendEvaluationResultsResponse
				err                           error
			}{
				sendEvaluationResultsResponse: &cliClient.SendEvaluationResultsResponse{
					EvaluationId:  1,
					PromptMessage: "",
				},
			},
		},
		expected: struct {
			ValidateResources struct {
				filesConfigurationsChan chan *extractor.FileConfigurations
			}
			Evaluate struct {
				evaluationData evaluation.PolicyCheckData
			}
			SendEvaluationResult struct {
				evaluationRequestData evaluation.EvaluationRequestData
			}
			PrintWarnings     []printer.Warning
			EvaluationSummary printer.EvaluationSummary
			err               error
		}{
			ValidateResources: struct {
				filesConfigurationsChan chan *extractor.FileConfigurations
			}{
				filesConfigurationsChan: filesConfigurationsChan,
			},
			Evaluate: struct {
				evaluationData evaluation.PolicyCheckData
			}{
				evaluationData: evaluation.PolicyCheckData{
					FilesConfigurations: validK8sFilesConfigurations,
					IsInteractiveMode:   true,
					PolicyName:          policy.Name,
					Policy:              policy,
				},
			},
			SendEvaluationResult: struct {
				evaluationRequestData evaluation.EvaluationRequestData
			}{
				evaluationRequestData: evaluation.EvaluationRequestData{
					PolicyName: policy.Name,
					FilesData:  []cliClient.FileData{{FilePath: validPath, ConfigurationsCount: 5}},
				},
			},
			PrintWarnings: []printer.Warning{},
			EvaluationSummary: printer.EvaluationSummary{
				ConfigsCount:              0,
				RulesCount:                len(policy.Rules),
				FilesCount:                1,
				PassedYamlValidationCount: 1,
				K8sValidation:             "1/1",
				PassedPolicyCheckCount:    1,
			},
			err: nil,
		},
	}
}

func test_no_k8s_resources_found() *TestFlowTestCase {
	root := pathFromRoot("internal/fixtures/nonKube/")
	preRunData := mockGetPreRunData()
	policy, _ := policy_factory.CreatePolicy(preRunData.PoliciesJson, "", preRunData.RegistrationURL)
	paths := []string{root + "/docker-compose-config.yaml", root + "/simple.json", root + "/simple.yaml"}
	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 3)
	go func() {
		defer func() {
			close(filesConfigurationsChan)
		}()

		filesConfigurationsChan <- &extractor.FileConfigurations{FileName: paths[0]}
		filesConfigurationsChan <- &extractor.FileConfigurations{FileName: paths[1]}
		filesConfigurationsChan <- &extractor.FileConfigurations{FileName: paths[2]}
	}()
	invalidFilesChan := make(chan *extractor.InvalidFile)
	filesk8sConfigurationsChan := make(chan *extractor.FileConfigurations)
	invalidK8sFilesChan := make(chan *extractor.InvalidFile, 3)
	go func() {
		defer func() {
			close(invalidK8sFilesChan)
		}()
		invalidK8sFilesChan <- &extractor.InvalidFile{Path: paths[0]}
		invalidK8sFilesChan <- &extractor.InvalidFile{Path: paths[1]}
		invalidK8sFilesChan <- &extractor.InvalidFile{Path: paths[2]}
	}()
	close(invalidFilesChan)
	close(filesk8sConfigurationsChan)
	return &TestFlowTestCase{
		name: "files are valid yaml but not valid k8s manifest files",
		args: struct {
			path []string
		}{
			path: paths,
		},
		mock: struct {
			ExtractFilesConfigurations struct {
				filesConfigurationsChan chan *extractor.FileConfigurations
				invalidFilesChan        chan *extractor.InvalidFile
			}
			ValidateResources struct {
				k8sFilesConfigurationsChan chan *extractor.FileConfigurations
				k8sInvalidFilesChan        chan *extractor.InvalidFile
				filesWithWarningsChan      chan *validation.FileWithWarning
			}
			Evaluate struct {
				policyCheckResultData evaluation.PolicyCheckResultData
				err                   error
			}
			SendEvaluationResult struct {
				sendEvaluationResultsResponse *cliClient.SendEvaluationResultsResponse
				err                           error
			}
		}{
			ExtractFilesConfigurations: struct {
				filesConfigurationsChan chan *extractor.FileConfigurations
				invalidFilesChan        chan *extractor.InvalidFile
			}{
				filesConfigurationsChan: filesConfigurationsChan,
				invalidFilesChan:        invalidFilesChan,
			},
			ValidateResources: struct {
				k8sFilesConfigurationsChan chan *extractor.FileConfigurations
				k8sInvalidFilesChan        chan *extractor.InvalidFile
				filesWithWarningsChan      chan *validation.FileWithWarning
			}{
				k8sFilesConfigurationsChan: filesk8sConfigurationsChan,
				k8sInvalidFilesChan:        invalidK8sFilesChan,
				filesWithWarningsChan:      newK8sValidationWarningsChan(),
			},
			Evaluate: struct {
				policyCheckResultData evaluation.PolicyCheckResultData
				err                   error
			}{
				policyCheckResultData: evaluation.PolicyCheckResultData{
					FormattedResults: evaluation.FormattedResults{},
					RulesData:        []cliClient.RuleData{},
					FilesData:        []cliClient.FileData{},
					RawResults:       evaluation.FailedRulesByFiles{},
					RulesCount:       len(policy.Rules),
				},
				err: nil,
			},
			SendEvaluationResult: struct {
				sendEvaluationResultsResponse *cliClient.SendEvaluationResultsResponse
				err                           error
			}{
				sendEvaluationResultsResponse: &cliClient.SendEvaluationResultsResponse{
					EvaluationId:  1,
					PromptMessage: "",
				},
				err: nil,
			},
		},
		expected: struct {
			ValidateResources struct {
				filesConfigurationsChan chan *extractor.FileConfigurations
			}
			Evaluate struct {
				evaluationData evaluation.PolicyCheckData
			}
			SendEvaluationResult struct {
				evaluationRequestData evaluation.EvaluationRequestData
			}
			PrintWarnings     []printer.Warning
			EvaluationSummary printer.EvaluationSummary
			err               error
		}{
			ValidateResources: struct {
				filesConfigurationsChan chan *extractor.FileConfigurations
			}{
				filesConfigurationsChan: filesConfigurationsChan,
			},
			Evaluate: struct {
				evaluationData evaluation.PolicyCheckData
			}{
				evaluationData: evaluation.PolicyCheckData{
					FilesConfigurations: []*extractor.FileConfigurations{},
					IsInteractiveMode:   true,
					PolicyName:          policy.Name,
					Policy:              policy,
				},
			},
			PrintWarnings: make([]printer.Warning, 3),
			EvaluationSummary: printer.EvaluationSummary{
				ConfigsCount:              0,
				RulesCount:                len(policy.Rules),
				FilesCount:                3,
				PassedYamlValidationCount: 3,
				K8sValidation:             "0/3",
				PassedPolicyCheckCount:    0,
			},
			SendEvaluationResult: struct {
				evaluationRequestData evaluation.EvaluationRequestData
			}{
				evaluationRequestData: evaluation.EvaluationRequestData{
					PolicyName:     policy.Name,
					FilesData:      []cliClient.FileData{},
					FailedK8sFiles: paths,
				},
			},
			err: errors.New(""),
		},
	}
}

func setup() {
	evaluationId = 444

	prerunData := mockGetPreRunData()

	formattedResults := evaluation.FormattedResults{}

	policyCheckResultData := evaluation.PolicyCheckResultData{
		FormattedResults: formattedResults,
		RulesData:        []cliClient.RuleData{},
		FilesData:        []cliClient.FileData{},
		RawResults:       nil,
		RulesCount:       0,
	}

	formattedResults.EvaluationResults = &evaluation.EvaluationResults{
		FileNameRuleMapper: map[string]map[string]*evaluation.Rule{},
		Summary: evaluation.EvaluationResultsSummery{
			TotalFailedRules:  0,
			TotalSkippedRules: 0,
			FilesCount:        0,
			FilesPassedCount:  1,
		},
	}

	sendEvaluationResultsResponse := &cliClient.SendEvaluationResultsResponse{
		EvaluationId:  1,
		PromptMessage: "",
	}

	mockedEvaluator = &mockEvaluator{}
	mockedEvaluator.On("Evaluate", mock.Anything).Return(policyCheckResultData, nil)
	mockedEvaluator.On("SendEvaluationResult", mock.Anything).Return(sendEvaluationResultsResponse, nil)

	messagerMock = &mockMessager{}
	messagerMock.On("LoadVersionMessages", mock.Anything)

	k8sValidatorMock = &K8sValidatorMock{}

	path := "valid/path"
	filesConfigurationsChan := newFilesConfigurationsChan(path)
	filesConfigurations = newFilesConfigurations(path)

	invalidFilesChan := make(chan *extractor.InvalidFile, 100)
	close(invalidFilesChan)
	invalidK8sFilesChan := newInvalidK8sFilesChan()
	ignoredFilesChan := newIgnoredYamlFilesChan()
	k8sValidationWarningsChan := newK8sValidationWarningsChan()

	k8sValidatorMock.On("ValidateResources", mock.Anything, mock.Anything).Return(filesConfigurationsChan, invalidK8sFilesChan, k8sValidationWarningsChan, newErrorsChan())
	k8sValidatorMock.On("GetK8sFiles", mock.Anything, mock.Anything).Return(filesConfigurationsChan, ignoredFilesChan, newErrorsChan())
	k8sValidatorMock.On("InitClient", mock.Anything, mock.Anything, mock.Anything).Return()

	filesExtractorMock := &FilesExtractorMock{}
	filesExtractorMock.On("ExtractFilesConfigurations", mock.Anything, 100).Return(filesConfigurationsChan, invalidFilesChan)
	printerMock := &PrinterMock{}
	printerMock.On("PrintWarnings", mock.Anything)
	printerMock.On("PrintSummaryTable", mock.Anything)
	printerMock.On("PrintEvaluationSummary", mock.Anything, mock.Anything)
	printerMock.On("PrintMessage", mock.Anything, mock.Anything)
	printerMock.On("PrintPromptMessage", mock.Anything)
	printerMock.On("SetTheme", mock.Anything)

	readerMock = &ReaderMock{}
	readerMock.On("FilterFiles", []string{"8/*"}).Return([]string{"file/path"}, nil)
	readerMock.On("FilterFiles", []string{"valid/path"}).Return([]string{"file/path"}, nil)

	localConfigMock = &LocalConfigMock{}
	localConfigMock.On("GetLocalConfiguration").Return(&localConfig.LocalConfig{Token: "134kh"}, nil)

	mockedCliClient = &CliClientMock{}
	mockedCliClient.On("RequestEvaluationPrerunData", mock.Anything).Return(prerunData, nil)

	ctx = &TestCommandContext{
		K8sValidator:   k8sValidatorMock,
		Evaluator:      mockedEvaluator,
		LocalConfig:    localConfigMock,
		Messager:       messagerMock,
		Printer:        printerMock,
		Reader:         readerMock,
		FilesExtractor: filesExtractorMock,
		CliClient:      mockedCliClient,
	}

	testingPolicy, _ = policy_factory.CreatePolicy(prerunData.PoliciesJson, "", prerunData.RegistrationURL)
}

func TestTestCommandFlagsValidation(t *testing.T) {
	setup()
	test_testCommand_output_flags_validation(t, ctx)
	test_testCommand_version_flags_validation(t, ctx)
	test_testCommand_no_record_flag(t, ctx)
}

func TestTestCommandEmptyDir(t *testing.T) {
	setup()
	emptyDir := t.TempDir()
	emptyDirPaths := filepath.Join(emptyDir, "*.yaml")

	readerMock.On("FilterFiles", []string{emptyDirPaths}).Return([]string{}, nil)
	err := Test(ctx, []string{emptyDirPaths}, &TestCommandData{K8sVersion: "1.18.0", Output: "", Policy: testingPolicy, Token: "134kh"})

	assert.EqualError(t, err, "No files detected")
}
func TestTestCommandNoFlags(t *testing.T) {
	setup()
	_ = Test(ctx, []string{"8/*"}, &TestCommandData{K8sVersion: "1.18.0", Output: "", Policy: testingPolicy, Token: "134kh"})

	policyCheckData := evaluation.PolicyCheckData{
		FilesConfigurations: filesConfigurations,
		IsInteractiveMode:   true,
		PolicyName:          testingPolicy.Name,
		Policy:              testingPolicy,
	}

	k8sValidatorMock.AssertCalled(t, "ValidateResources", mock.Anything, 100)
	mockedEvaluator.AssertCalled(t, "Evaluate", policyCheckData)
}

func TestTestCommandJsonOutput(t *testing.T) {
	setup()
	_ = Test(ctx, []string{"valid/path"}, &TestCommandData{Output: "json", Policy: testingPolicy})

	policyCheckData := evaluation.PolicyCheckData{
		FilesConfigurations: filesConfigurations,
		IsInteractiveMode:   false,
		PolicyName:          testingPolicy.Name,
		Policy:              testingPolicy,
	}

	k8sValidatorMock.AssertCalled(t, "ValidateResources", mock.Anything, 100)
	mockedEvaluator.AssertCalled(t, "Evaluate", policyCheckData)
}

func TestTestCommandYamlOutput(t *testing.T) {
	setup()
	_ = Test(ctx, []string{"8/*"}, &TestCommandData{Output: "yaml", Policy: testingPolicy})

	policyCheckData := evaluation.PolicyCheckData{
		FilesConfigurations: filesConfigurations,
		IsInteractiveMode:   false,
		PolicyName:          testingPolicy.Name,
		Policy:              testingPolicy,
	}

	k8sValidatorMock.AssertCalled(t, "ValidateResources", mock.Anything, 100)
	mockedEvaluator.AssertCalled(t, "Evaluate", policyCheckData)
}

func TestTestCommandXmlOutput(t *testing.T) {
	setup()
	_ = Test(ctx, []string{"valid/path"}, &TestCommandData{Output: "xml", Policy: testingPolicy})

	policyCheckData := evaluation.PolicyCheckData{
		FilesConfigurations: filesConfigurations,
		IsInteractiveMode:   false,
		PolicyName:          testingPolicy.Name,
		Policy:              testingPolicy,
	}

	k8sValidatorMock.AssertCalled(t, "ValidateResources", mock.Anything, 100)
	mockedEvaluator.AssertCalled(t, "Evaluate", policyCheckData)
}

func TestTestCommandOnlyK8sFiles(t *testing.T) {
	setup()
	_ = Test(ctx, []string{"8/*"}, &TestCommandData{OnlyK8sFiles: true})

	k8sValidatorMock.AssertCalled(t, "ValidateResources", mock.Anything, 100)
	k8sValidatorMock.AssertCalled(t, "GetK8sFiles", mock.Anything, 100)
}

func TestTestCommandNoInternetConnection(t *testing.T) {
	setup()
	_ = Test(ctx, []string{"valid/path"}, &TestCommandData{Policy: testingPolicy})

	policyCheckData := evaluation.PolicyCheckData{
		FilesConfigurations: filesConfigurations,
		IsInteractiveMode:   true,
		PolicyName:          testingPolicy.Name,
		Policy:              testingPolicy,
	}

	path := "valid/path"
	filesConfigurationsChan := newFilesConfigurationsChan(path)
	invalidK8sFilesChan := newInvalidK8sFilesChan()
	K8sValidationWarnings := validation.K8sValidationWarningPerValidFile{"valid/path": "Validation warning message - no internet"}

	k8sValidatorMock.On("ValidateResources", mock.Anything, mock.Anything).Return(filesConfigurationsChan, invalidK8sFilesChan, K8sValidationWarnings, newErrorsChan())

	k8sValidatorMock.AssertCalled(t, "ValidateResources", mock.Anything, 100)
	mockedEvaluator.AssertCalled(t, "Evaluate", policyCheckData)
}

func executeTestCommand(ctx *TestCommandContext, args []string) error {
	cmd := New(ctx)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return err
}

func test_testCommand_output_flags_validation(t *testing.T, ctx *TestCommandContext) {

	validOutputValues := [5]string{"simple", "json", "yaml", "xml", "JUnit"}

	for _, value := range validOutputValues {
		flags := TestCommandFlags{Output: value}
		err := flags.Validate()
		assert.NoError(t, err)
	}

	values := []string{"Simple", "Json", "Yaml", "Xml", "JUnit", "invalid", "113", "true"}

	for _, value := range values {
		err := executeTestCommand(ctx, []string{"8/*", "--output=" + value})
		expectedErrorStr := "Invalid --output option - \"" + value + "\"\n" +
			"Valid output values are - simple, yaml, json, xml, JUnit\n"
		assert.EqualError(t, err, expectedErrorStr)
	}
}

func test_testCommand_version_flags_validation(t *testing.T, ctx *TestCommandContext) {
	getExpectedErrorStr := func(value string) string {
		expectedStr := "The specified schema-version \"" + value + "\" is not in the correct format.\n" +
			"Make sure you are following the semantic versioning format <MAJOR>.<MINOR>.<PATCH>\n" +
			"Read more about kubernetes versioning: https://kubernetes.io/releases/version-skew-policy/#supported-versions"
		return expectedStr
	}

	values := []string{"1", "1.15", "1.15.", "1.15.0.", "1.15.0.1", "1..15.0", "str.12.bool"}
	for _, value := range values {
		err := executeTestCommand(ctx, []string{"8/*", "--schema-version=" + value})
		assert.EqualError(t, err, getExpectedErrorStr(value))
	}

	flags := TestCommandFlags{K8sVersion: "1.21.0"}
	err := flags.Validate()
	assert.NoError(t, err)
}

func test_testCommand_no_record_flag(t *testing.T, ctx *TestCommandContext) {
	err := executeTestCommand(ctx, []string{"8/*", "--no-record"})
	mockedEvaluator.AssertNotCalled(t, "SendEvaluationResult")
	assert.Equal(t, ViolationsFoundError, err)
}

func newFilesConfigurationsChan(path string) chan *extractor.FileConfigurations {
	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)

	go func() {
		filesConfigurationsChan <- &extractor.FileConfigurations{
			FileName: path,
		}
		close(filesConfigurationsChan)
	}()

	return filesConfigurationsChan
}

func newFilesConfigurations(path string) []*extractor.FileConfigurations {
	var filesConfigurations []*extractor.FileConfigurations
	filesConfigurations = append(filesConfigurations, &extractor.FileConfigurations{
		FileName: path,
	})
	return filesConfigurations
}

func newInvalidK8sFilesChan() chan *extractor.InvalidFile {
	invalidFilesChan := make(chan *extractor.InvalidFile, 1)

	invalidFile := &extractor.InvalidFile{
		Path:             "invalid/path",
		ValidationErrors: []error{},
	}

	go func() {
		invalidFilesChan <- invalidFile
		close(invalidFilesChan)
	}()

	return invalidFilesChan
}

func newIgnoredYamlFilesChan() chan *extractor.FileConfigurations {
	ignoredFilesChan := make(chan *extractor.FileConfigurations)
	ignoredFile := &extractor.FileConfigurations{
		FileName: "path/to/ignored/file",
	}

	go func() {
		ignoredFilesChan <- ignoredFile
		close(ignoredFilesChan)
	}()

	return ignoredFilesChan
}

func newK8sValidationWarningsChan() chan *validation.FileWithWarning {
	k8sValidationWarningsChan := make(chan *validation.FileWithWarning, 1)
	go func() {
		close(k8sValidationWarningsChan)
	}()

	return k8sValidationWarningsChan
}

func newErrorsChan() chan error {
	invalidFilesChan := make(chan error, 1)

	close(invalidFilesChan)
	return invalidFilesChan
}

func mockGetPreRunData() *cliClient.EvaluationPrerunDataResponse {
	const policiesJsonPath = "../../internal/fixtures/policyAsCode/prerun.json"

	fileReader := fileReader.CreateFileReader(nil)
	policiesJsonStr, err := fileReader.ReadFileContent(policiesJsonPath)

	if err != nil {
		panic(err)
	}

	policiesJsonRawData := []byte(policiesJsonStr)

	var policiesJson *cliClient.EvaluationPrerunDataResponse
	err = json.Unmarshal(policiesJsonRawData, &policiesJson)

	if err != nil {
		panic(err)
	}
	return policiesJson
}
