package config

import (
	"testing"

	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/stretchr/testify/mock"
)

func (lc *LocalConfigMock) Get(key string) string {
	lc.Called(key)
	return ""
}

func TestGetCommand(t *testing.T) {
	messager := &mockMessager{}
	messager.On("LoadVersionMessages", mock.Anything)

	printerMock := &PrinterMock{}
	printerMock.On("PrintWarnings", mock.Anything)
	printerMock.On("PrintSummaryTable", mock.Anything)
	printerMock.On("PrintMessage", mock.Anything, mock.Anything)
	printerMock.On("PrintEvaluationSummary", mock.Anything, mock.Anything)

	localConfigMock := &LocalConfigMock{}
	localConfigMock.On("GetLocalConfiguration").Return(&localConfig.LocalConfig{Token: "previousToken"}, nil)
	localConfigMock.On("Get", mock.Anything)

	ctx := &ConfigCommandContext{
		Messager:    messager,
		CliVersion:  "1.2.3",
		Printer:     printerMock,
		LocalConfig: localConfigMock,
	}

	get(ctx, "testkey")
	localConfigMock.AssertCalled(t, "Get", "testkey")
}
