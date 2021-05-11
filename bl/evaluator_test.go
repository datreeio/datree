package bl

import (
	"fmt"
	"testing"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/datreeio/datree/pkg/propertiesExtractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPropertiesExtractor struct {
	mock.Mock
}

func (m *mockPropertiesExtractor) ReadFilesFromPaths(paths []string, conc int) ([]*propertiesExtractor.FileProperties, []propertiesExtractor.FileError, []error) {
	args := m.Called(paths, conc)
	return args.Get(0).([]*propertiesExtractor.FileProperties), args.Get(1).([]propertiesExtractor.FileError), args.Get(2).([]error)
}

type mockCliClient struct {
	mock.Mock
}

func (m *mockCliClient) CreateEvaluation(createEvaluationRequest cliClient.CreateEvaluationRequest) (int, error) {
	args := m.Called(createEvaluationRequest)
	return args.Get(0).(int), args.Error(1)
}

func (m *mockCliClient) RequestEvaluation(evaluationRequest cliClient.EvaluationRequest) (cliClient.EvaluationResponse, error) {
	args := m.Called(evaluationRequest)
	return args.Get(0).(cliClient.EvaluationResponse), args.Error(1)
}

type mockPrinter struct {
	mock.Mock
}

func (m *mockPrinter) PrintWarnings(warnings []printer.Warning) {
	m.Called(warnings)
}

func (c *mockPrinter) PrintSummaryTable(summary printer.Summary) {
	c.Called(summary)
}

type propertiesExtractorMockTestCase struct {
	readFilesFromPaths struct {
		properties  []*propertiesExtractor.FileProperties
		filesErrors []propertiesExtractor.FileError
		errors      []error
	}
}

type cliClientMockTestCase struct {
	createEvaluation struct {
		EvaluationId int
		errors       error
	}
	requestEvaluation struct {
		response cliClient.EvaluationResponse
		errors   error
	}
}

type evaluateTestCase struct {
	name string
	args struct {
		paths                   []string
		cliId                   string
		evaluationConc          int
		evaluationRequest       cliClient.EvaluationRequest
		createEvaluationRequest cliClient.CreateEvaluationRequest
	}
	mock struct {
		propertiesExtractor propertiesExtractorMockTestCase
		cliClient           cliClientMockTestCase
	}
	expected struct {
		response   *EvaluationResults
		fileErrors []propertiesExtractor.FileError
		err        error
	}
}

func TestEvaluate(t *testing.T) {
	tests := []*evaluateTestCase{
		test_create_evaluation_failedRequest(),
		//test_evaluate_failedRequest(),
		//test_evaluate_success(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			propertiesExtractor := &mockPropertiesExtractor{}
			cliClient := &mockCliClient{}
			printer := &mockPrinter{}

			propertiesExtractor.On("ReadFilesFromPaths", mock.Anything, mock.Anything).Return(tt.mock.propertiesExtractor.readFilesFromPaths.properties, tt.mock.propertiesExtractor.readFilesFromPaths.filesErrors, tt.mock.propertiesExtractor.readFilesFromPaths.errors)
			cliClient.On("RequestEvaluation", mock.Anything).Return(tt.mock.cliClient.requestEvaluation.response, tt.mock.cliClient.requestEvaluation.errors)
			cliClient.On("CreateEvaluation", mock.Anything).Return(tt.mock.cliClient.createEvaluation.EvaluationId, tt.mock.cliClient.createEvaluation.errors)

			evaluator := &Evaluator{
				propertiesExtractor: propertiesExtractor,
				cliClient:           cliClient,
				printer:             printer,
				osInfo: &OSInfo{
					OS:              "darwin",
					PlatformVersion: "1.2.3",
					KernelVersion:   "4.5.6",
				},
			}

			actualResponse, actualFilesErrs, actualErr := evaluator.Evaluate(tt.args.paths, tt.args.cliId, tt.args.evaluationConc, "0.0.1")

			propertiesExtractor.AssertCalled(t, "ReadFilesFromPaths", tt.args.paths, tt.args.evaluationConc)
			cliClient.AssertCalled(t, "CreateEvaluation", tt.args.createEvaluationRequest)
			if tt.mock.cliClient.createEvaluation.errors == nil {
				cliClient.AssertCalled(t, "RequestEvaluation", tt.args.evaluationRequest)
			}

			assert.Equal(t, tt.expected.response, actualResponse)
			assert.Equal(t, tt.expected.fileErrors, actualFilesErrs)
			assert.Equal(t, tt.expected.err, actualErr)
		})
	}
}

func test_evaluate_success() *evaluateTestCase {
	evaluationResponse := cliClient.EvaluationResponse{
		Results: []cliClient.EvaluationResult{
			{
				Passed: true,
				Results: struct {
					Matches    []cliClient.Match "json:\"matches\""
					Mismatches []cliClient.Match "json:\"mismatches\""
				}{
					Matches:    []cliClient.Match{},
					Mismatches: []cliClient.Match{},
				},
			},
		},
	}
	createEvaluationResponse := cliClient.CreateEvaluationResponse{
		EvaluationId: 432,
	}
	return &evaluateTestCase{
		name: "success",
		args: struct {
			paths                   []string
			cliId                   string
			evaluationConc          int
			evaluationRequest       cliClient.EvaluationRequest
			createEvaluationRequest cliClient.CreateEvaluationRequest
		}{
			paths:          []string{"path1/path2/file.yaml"},
			cliId:          "cliId-test",
			evaluationConc: 1,
			evaluationRequest: cliClient.EvaluationRequest{
				EvaluationId: createEvaluationResponse.EvaluationId,
				Files: []propertiesExtractor.FileProperties{{
					FileName:       "path1/path2/file.yaml",
					Configurations: []propertiesExtractor.K8sConfiguration{{"apiVersion": "extensions/v1beta1"}},
				}},
			},
			createEvaluationRequest: cliClient.CreateEvaluationRequest{
				CliId: "cliId-test",
				Metadata: cliClient.Metadata{
					CliVersion:      "0.0.1",
					Os:              "darwin",
					PlatformVersion: "1.2.3",
					KernelVersion:   "4.5.6",
				},
			},
		},
		mock: struct {
			propertiesExtractor propertiesExtractorMockTestCase
			cliClient           cliClientMockTestCase
		}{
			propertiesExtractor: propertiesExtractorMockTestCase{
				readFilesFromPaths: struct {
					properties  []*propertiesExtractor.FileProperties
					filesErrors []propertiesExtractor.FileError
					errors      []error
				}{
					properties: []*propertiesExtractor.FileProperties{{
						FileName:       "path1/path2/file.yaml",
						Configurations: []propertiesExtractor.K8sConfiguration{{"apiVersion": "extensions/v1beta1"}},
					}},
					filesErrors: []propertiesExtractor.FileError{},
					errors:      []error{},
				},
			},
			cliClient: cliClientMockTestCase{
				requestEvaluation: struct {
					response cliClient.EvaluationResponse
					errors   error
				}{
					response: evaluationResponse,
					errors:   nil,
				},
				createEvaluation: struct {
					EvaluationId int
					errors       error
				}{
					EvaluationId: createEvaluationResponse.EvaluationId,
				},
			},
		},
		expected: struct {
			response   *EvaluationResults
			fileErrors []propertiesExtractor.FileError
			err        error
		}{
			response: &EvaluationResults{
				FileNameRuleMapper: map[string]map[int]*Rule{}, Summary: struct {
					RulesCount       int
					TotalFailedRules int
					FilesCount       int
				}{RulesCount: 1, TotalFailedRules: 0, FilesCount: 1},
			},
			fileErrors: []propertiesExtractor.FileError{},
			err:        nil,
		},
	}
}

func test_evaluate_failedRequest() *evaluateTestCase {
	createEvaluationResponse := cliClient.CreateEvaluationResponse{
		EvaluationId: 432,
	}
	evaluationResponse := cliClient.EvaluationResponse{}
	return &evaluateTestCase{
		name: "fail",
		args: struct {
			paths                   []string
			cliId                   string
			evaluationConc          int
			evaluationRequest       cliClient.EvaluationRequest
			createEvaluationRequest cliClient.CreateEvaluationRequest
		}{
			paths:          []string{"path1/path2/file.yaml"},
			cliId:          "cliId-test",
			evaluationConc: 1,
			evaluationRequest: cliClient.EvaluationRequest{
				EvaluationId: createEvaluationResponse.EvaluationId,
				Files: []propertiesExtractor.FileProperties{{
					FileName:       "path1/path2/file.yaml",
					Configurations: []propertiesExtractor.K8sConfiguration{{"apiVersion": "extensions/v1beta1"}},
				}},
			},
			createEvaluationRequest: cliClient.CreateEvaluationRequest{
				CliId: "cliId-test",
				Metadata: cliClient.Metadata{
					CliVersion:      "0.0.1",
					Os:              "darwin",
					PlatformVersion: "1.2.3",
					KernelVersion:   "4.5.6",
				},
			},
		},
		mock: struct {
			propertiesExtractor propertiesExtractorMockTestCase
			cliClient           cliClientMockTestCase
		}{
			propertiesExtractor: propertiesExtractorMockTestCase{
				readFilesFromPaths: struct {
					properties  []*propertiesExtractor.FileProperties
					filesErrors []propertiesExtractor.FileError
					errors      []error
				}{
					properties: []*propertiesExtractor.FileProperties{{
						FileName:       "path1/path2/file.yaml",
						Configurations: []propertiesExtractor.K8sConfiguration{{"apiVersion": "extensions/v1beta1"}}}},
					filesErrors: []propertiesExtractor.FileError{},
					errors:      []error{},
				},
			},
			cliClient: cliClientMockTestCase{
				requestEvaluation: struct {
					response cliClient.EvaluationResponse
					errors   error
				}{
					response: evaluationResponse,
					errors:   fmt.Errorf("error"),
				},
				createEvaluation: struct {
					EvaluationId int
					errors       error
				}{
					EvaluationId: createEvaluationResponse.EvaluationId,
					errors:       nil,
				},
			},
		},
		expected: struct {
			response   *EvaluationResults
			fileErrors []propertiesExtractor.FileError
			err        error
		}{
			response:   nil,
			fileErrors: []propertiesExtractor.FileError{},
			err:        fmt.Errorf("error"),
		},
	}
}

func test_create_evaluation_failedRequest() *evaluateTestCase {
	createEvaluationResponse := cliClient.CreateEvaluationResponse{}
	evaluationResponse := cliClient.EvaluationResponse{}
	return &evaluateTestCase{
		name: "fail",
		args: struct {
			paths                   []string
			cliId                   string
			evaluationConc          int
			evaluationRequest       cliClient.EvaluationRequest
			createEvaluationRequest cliClient.CreateEvaluationRequest
		}{
			paths:          []string{"path1/path2/file.yaml"},
			cliId:          "cliId-test",
			evaluationConc: 1,
			evaluationRequest: cliClient.EvaluationRequest{
				EvaluationId: createEvaluationResponse.EvaluationId,
				Files: []propertiesExtractor.FileProperties{{
					FileName:       "path1/path2/file.yaml",
					Configurations: []propertiesExtractor.K8sConfiguration{{"apiVersion": "extensions/v1beta1"}},
				}},
			},
			createEvaluationRequest: cliClient.CreateEvaluationRequest{
				CliId: "cliId-test",
				Metadata: cliClient.Metadata{
					CliVersion:      "0.0.1",
					Os:              "darwin",
					PlatformVersion: "1.2.3",
					KernelVersion:   "4.5.6",
				},
			},
		},
		mock: struct {
			propertiesExtractor propertiesExtractorMockTestCase
			cliClient           cliClientMockTestCase
		}{
			propertiesExtractor: propertiesExtractorMockTestCase{
				readFilesFromPaths: struct {
					properties  []*propertiesExtractor.FileProperties
					filesErrors []propertiesExtractor.FileError
					errors      []error
				}{
					properties: []*propertiesExtractor.FileProperties{{
						FileName:       "path1/path2/file.yaml",
						Configurations: []propertiesExtractor.K8sConfiguration{{"apiVersion": "extensions/v1beta1"}}}},
					filesErrors: []propertiesExtractor.FileError{},
					errors:      []error{},
				},
			},
			cliClient: cliClientMockTestCase{
				requestEvaluation: struct {
					response cliClient.EvaluationResponse
					errors   error
				}{
					response: evaluationResponse,
					errors:   fmt.Errorf("error"),
				},
				createEvaluation: struct {
					EvaluationId int
					errors       error
				}{
					errors: fmt.Errorf("create evaluation error"),
				},
			},
		},
		expected: struct {
			response   *EvaluationResults
			fileErrors []propertiesExtractor.FileError
			err        error
		}{
			response:   nil,
			fileErrors: []propertiesExtractor.FileError{},
			err:        fmt.Errorf("create evaluation error"),
		},
	}
}

func TestPrintResults(t *testing.T) {
	printerSpy := &mockPrinter{}

	printerSpy.On("PrintWarnings", mock.Anything).Return()
	printerSpy.On("PrintSummaryTable", mock.Anything).Return()

	evaluator := &Evaluator{
		propertiesExtractor: &mockPropertiesExtractor{},
		cliClient:           &mockCliClient{},
		printer:             printerSpy,
	}

	results := EvaluationResults{
		FileNameRuleMapper: map[string]map[int]*Rule{},
		Summary: struct {
			RulesCount       int
			TotalFailedRules int
			FilesCount       int
		}{},
	}

	evaluator.PrintResults(&results, "cli_id", "")

	expectedPrinterWarnings := []printer.Warning{}
	printerSpy.AssertCalled(t, "PrintWarnings", expectedPrinterWarnings)

	plainRows := []printer.SummaryItem{
		{
			RightCol: "0",
			LeftCol:  "Enabled rules in policy “default”",
			RowIndex: 0,
		},
		{
			RightCol: "0",
			LeftCol:  "Configs tested against policy",
			RowIndex: 1,
		},
		{
			RightCol: "0",
			LeftCol:  "Total rules evaluated",
			RowIndex: 2,
		},
		{
			RightCol: "https://app.datree.io/login?cliId=cli_id",
			LeftCol:  "See all rules in policy",
			RowIndex: 5,
		}}

	expectedPrinterSummary := printer.Summary{
		PlainRows: plainRows,
		ErrorRow: printer.SummaryItem{
			LeftCol:  "Total rules failed",
			RightCol: "0",
			RowIndex: 3,
		},
		SuccessRow: printer.SummaryItem{
			LeftCol:  "Total rules passed",
			RightCol: "0",
			RowIndex: 4,
		},
	}

	printerSpy.AssertCalled(t, "PrintSummaryTable", expectedPrinterSummary)

}
