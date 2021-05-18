package messager

import (
	"testing"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockMessagesClient struct {
	mock.Mock
}

func (m *mockMessagesClient) GetVersionMessage(cliVersion string, timeout int) (*cliClient.VersionMessage, error) {
	args := m.Called(cliVersion, timeout)
	return args.Get(0).(*cliClient.VersionMessage), args.Error(1)
}

func TestLoadVersionMessages(t *testing.T) {
	mockedMessagesClient := &mockMessagesClient{}
	mockedMessage := &cliClient.VersionMessage{
		CliVersion:   "0.1.1",
		MessageText:  "message",
		MessageColor: "yellow",
	}
	mockedMessagesClient.On("GetVersionMessage", mock.Anything, mock.Anything).Return(mockedMessage, nil)
	messager := &Messager{
		defaultTimeout: 900,
		messagesClient: mockedMessagesClient,
	}

	messagesChan := make(chan *VersionMessage)
	messager.LoadVersionMessages(messagesChan, "0.1.1")

	msg := <-messagesChan
	assert.Equal(t, mockedMessage.CliVersion, msg.CliVersion)
	assert.Equal(t, mockedMessage.MessageColor, msg.MessageColor)
	assert.Equal(t, mockedMessage.MessageText, msg.MessageText)
	mockedMessagesClient.AssertCalled(t, "GetVersionMessage", "0.1.1", 900)
}
