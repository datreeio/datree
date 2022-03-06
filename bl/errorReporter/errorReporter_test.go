package errorReporter

import (
	"errors"
	"testing"

	"github.com/datreeio/datree/cmd"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

// client mock
type mockCliClient struct {
	mock.Mock
}

func (m *mockCliClient) ReportCliError(reportCliErrorRequest cliClient.ReportCliErrorRequest) (StatusCode int, Error error) {
	m.Called(reportCliErrorRequest)
	return 201, nil
}

// config mock
type mockConfig struct {
	mock.Mock
}

func (lc *mockConfig) GetLocalConfiguration() (*localConfig.ConfigContent, error) {
	lc.Called()
	return &localConfig.ConfigContent{
		CliId:         "2qRg9jzJGcA73ftqEcXuBp",
		SchemaVersion: "111",
	}, nil
}

// test type
type ErrorReporterTestCaseArgs struct {
	panicErr interface{}
}
type ErrorReporterTestCase struct {
	name     string
	args     *ErrorReporterTestCaseArgs
	expected cliClient.ReportCliErrorRequest
}

func TestErrorReporter(t *testing.T) {
	tests := []*ErrorReporterTestCase{
		reportErrorWithError(),
		reportErrorWithStringError(),
	}
	mockedCliClient := &mockCliClient{}
	mockedConfig := &mockConfig{}
	mockedCliClient.On("ReportCliError", mock.Anything).Return(nil)
	mockedConfig.On("GetLocalConfiguration").Return(nil)
	errorReporter := &ErrorReporter{
		client: mockedCliClient,
		config: mockedConfig,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorReporter.ReportCliError(tt.args.panicErr)
			reportCliErrorCalledArgs := (mockedCliClient.Calls[0].Arguments.Get(0)).(cliClient.ReportCliErrorRequest)
			assert.Equal(t, tt.expected.ErrorMessage, reportCliErrorCalledArgs.ErrorMessage)
			assert.Equal(t, tt.expected.Token, reportCliErrorCalledArgs.Token)
			assert.Equal(t, tt.expected.CliVersion, reportCliErrorCalledArgs.CliVersion)
		})
	}
}

func reportErrorWithError() *ErrorReporterTestCase {
	return &ErrorReporterTestCase{
		name: "should call cli client with correct args, when error is of type error",
		args: &ErrorReporterTestCaseArgs{
			panicErr: errors.New("this is the error message"),
		},
		expected: cliClient.ReportCliErrorRequest{
			Token:        "2qRg9jzJGcA73ftqEcXuBp",
			CliVersion:   cmd.CliVersion,
			ErrorMessage: "this is the error message",
			StackTrace:   mock.Anything,
		},
	}
}

func reportErrorWithStringError() *ErrorReporterTestCase {
	return &ErrorReporterTestCase{
		name: "should call cli client with correct args, when error is string",
		args: &ErrorReporterTestCaseArgs{
			panicErr: "this is the error message",
		},
		expected: cliClient.ReportCliErrorRequest{
			Token:        "2qRg9jzJGcA73ftqEcXuBp",
			CliVersion:   cmd.CliVersion,
			ErrorMessage: "this is the error message",
			StackTrace:   mock.Anything,
		},
	}
}
