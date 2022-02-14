package publish

import (
	"errors"
	"strings"
	"testing"

	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type LocalConfigMock struct {
	mock.Mock
}

func (lc *LocalConfigMock) GetLocalConfiguration() (*localConfig.LocalConfig, error) {
	lc.Called()
	return &localConfig.LocalConfig{Token: "token"}, nil
}

type MessagerMock struct {
	mock.Mock
}

func (m *MessagerMock) LoadVersionMessages(cliVersion string) chan *messager.VersionMessage {
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

type PrinterMock struct {
	mock.Mock
}

func (p *PrinterMock) PrintMessage(messageText string, messageColor string) {
	p.Called(messageText, messageColor)
}

type PublishClientMock struct {
	mock.Mock
}

func (p *PublishClientMock) PublishPolicies(policiesConfiguration files.UnknownStruct, cliId string) (*cliClient.PublishFailedResponse, error) {
	args := p.Called(policiesConfiguration, cliId)
	return args.Get(0).(*cliClient.PublishFailedResponse), args.Error(1)
}

type TokenClientMock struct {
	mock.Mock
}

func (t *TokenClientMock) CreateToken() (*cliClient.CreateTokenResponse, error) {
	args := t.Called()
	return args.Get(0).(*cliClient.CreateTokenResponse), args.Error(1)
}

func TestPublishCommand(t *testing.T) {
	localConfigMock := &LocalConfigMock{}
	localConfigMock.On("GetLocalConfiguration")

	messagerMock := &MessagerMock{}
	messagerMock.On("LoadVersionMessages", mock.Anything)

	printerMock := &PrinterMock{}
	printerMock.On("PrintMessage", mock.Anything, mock.Anything)

	publishClientMock := &PublishClientMock{}

	ctx := &PublishCommandContext{
		CliVersion:       "0.0.1",
		LocalConfig:      localConfigMock,
		Messager:         messagerMock,
		Printer:          printerMock,
		PublishCliClient: publishClientMock,
	}

	localConfigContent, _ := ctx.LocalConfig.GetLocalConfiguration()

	testPublishCommandSuccess(t, ctx, publishClientMock, localConfigContent)
	testPublishCommandFailedYaml(t, ctx, localConfigContent)
	testPublishCommandFailedSchema(t, ctx, publishClientMock, localConfigContent)
}

func testPublishCommandSuccess(t *testing.T, ctx *PublishCommandContext, publishClientMock *PublishClientMock, localConfigContent *localConfig.LocalConfig) {
	publishClientMock.On("PublishPolicies", mock.Anything, mock.Anything).Return(&cliClient.PublishFailedResponse{}, nil).Once()
	_, err := publish(ctx, "../../fixtures/policyAsCode/valid-schema.yaml", localConfigContent)
	assert.Equal(t, nil, err)
}

func testPublishCommandFailedYaml(t *testing.T, ctx *PublishCommandContext, localConfigContent *localConfig.LocalConfig) {
	_, err := publish(ctx, "../../fixtures/policyAsCode/invalid-yaml.yaml", localConfigContent)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "yaml: line 2: did not find expected key", err.Error())
}

func testPublishCommandFailedSchema(t *testing.T, ctx *PublishCommandContext, publishClientMock *PublishClientMock, localConfigContent *localConfig.LocalConfig) {
	publishFailedPayloadMock := []string{"first error", "second error"}
	errMessage := strings.Join(publishFailedPayloadMock, ",")
	publishFailedResponseMock := &cliClient.PublishFailedResponse{
		Code:    "mocked code",
		Message: errMessage,
		Payload: publishFailedPayloadMock,
	}
	const fileNamePath = "../../fixtures/policyAsCode/invalid-schemas/duplicate-rule-id.yaml"
	publishClientMock.On("PublishPolicies", mock.Anything, mock.Anything).Return(publishFailedResponseMock, errors.New(errMessage)).Once()
	publishFailedRes, err := publish(ctx, fileNamePath, localConfigContent)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, errMessage, err.Error())
	assert.Equal(t, publishFailedResponseMock, publishFailedRes)
}
