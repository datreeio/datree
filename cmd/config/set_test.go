package config

import (
	"testing"

	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/stretchr/testify/mock"
)

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

	m.Called(messages, cliVersion)
	return messages
}

func (m *mockMessager) HandleVersionMessage(messageChannel <-chan *messager.VersionMessage) {
	m.Called(messageChannel)
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

func (p *PrinterMock) PrintMessage(messageText string, messageColor string) {
	p.Called(messageText, messageColor)
}

func (p *PrinterMock) PrintEvaluationSummary(evaluationSummary printer.EvaluationSummary, k8sVersion string) {
	p.Called(evaluationSummary)
}

type LocalConfigMock struct {
	mock.Mock
}

func (lc *LocalConfigMock) GetLocalConfiguration() (*localConfig.LocalConfig, error) {
	lc.Called()
	return &localConfig.LocalConfig{Token: "previousToken"}, nil
}

func (lc *LocalConfigMock) Set(key string, value string) error {
	lc.Called(key, value)
	return nil
}

func TestSetCommand(t *testing.T) {
	messager := &mockMessager{}
	messager.On("LoadVersionMessages", mock.Anything)

	printerMock := &PrinterMock{}
	printerMock.On("PrintWarnings", mock.Anything)
	printerMock.On("PrintSummaryTable", mock.Anything)
	printerMock.On("PrintMessage", mock.Anything, mock.Anything)
	printerMock.On("PrintEvaluationSummary", mock.Anything, mock.Anything)

	localConfigMock := &LocalConfigMock{}
	localConfigMock.On("GetLocalConfiguration").Return(&localConfig.LocalConfig{Token: "previousToken"}, nil)
	localConfigMock.On("Set", mock.Anything, mock.Anything)

	ctx := &ConfigCommandContext{
		Messager:    messager,
		CliVersion:  "1.2.3",
		Printer:     printerMock,
		LocalConfig: localConfigMock,
	}

	set(ctx, "testkey", "testvalue")
	localConfigMock.AssertCalled(t, "Set", "testkey", "testvalue")
}
