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

func (m *mockCliClient) ReportCliError(reportCliErrorRequest cliClient.ReportCliErrorRequest, uri string) (StatusCode int, Error error) {
	m.Called(reportCliErrorRequest)
	return 201, nil
}

// config mock
type mockConfig struct {
	mock.Mock
}

func (lc *mockConfig) GetLocalConfiguration() (*localConfig.LocalConfig, error) {
	lc.Called()
	return &localConfig.LocalConfig{
		Token:         "2qRg9jzJGcA73ftqEcXuBp",
		SchemaVersion: "111",
	}, nil
}

// printer mock
type mockPrinter struct {
	mock.Mock
}

func (mp *mockPrinter) PrintMessage(messageText string, messageColor string) {
	mp.Called()
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
	mockedPrinter := &mockPrinter{}
	mockedCliClient.On("ReportCliError", mock.Anything).Return(nil)
	mockedConfig.On("GetLocalConfiguration").Return(nil)
	mockedPrinter.On("PrintMessage", mock.Anything, mock.Anything).Return()
	errorReporter := &ErrorReporter{
		client:  mockedCliClient,
		config:  mockedConfig,
		printer: mockedPrinter,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorReporter.ReportError(tt.args.panicErr, "/report-cli-panic-error")
			reportCliErrorCalledArgs := (mockedCliClient.Calls[0].Arguments.Get(0)).(cliClient.ReportCliErrorRequest)
			assert.Equal(t, tt.expected.ErrorMessage, reportCliErrorCalledArgs.ErrorMessage)
			assert.Equal(t, tt.expected.ClientId, reportCliErrorCalledArgs.ClientId)
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
			ClientId:     "2qRg9jzJGcA73ftqEcXuBp",
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
			ClientId:     "2qRg9jzJGcA73ftqEcXuBp",
			Token:        "2qRg9jzJGcA73ftqEcXuBp",
			CliVersion:   cmd.CliVersion,
			ErrorMessage: "this is the error message",
			StackTrace:   mock.Anything,
		},
	}
}
