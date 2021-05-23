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

type LoadVersionMessagesTestCaseArgs struct {
	messagesChan chan *VersionMessage
	cliVersion   string
}

type LoadVersionMessagesTestCaseMock struct {
	message *cliClient.VersionMessage
	timeout int
}

type LoadVersionMessagesTestCaseExpected struct {
	message *VersionMessage
}
type LoadVersionMessagesTestCase struct {
	name     string
	mock     *LoadVersionMessagesTestCaseMock
	args     *LoadVersionMessagesTestCaseArgs
	expected *LoadVersionMessagesTestCaseExpected
}

func TestLoadVersionMessages(t *testing.T) {
	tests := []*LoadVersionMessagesTestCase{
		load_message_success(),
		load_message_failed(),
	}
	for _, tt := range tests {
		mockedMessagesClient := &mockMessagesClient{}
		t.Run(tt.name, func(t *testing.T) {
			mockedMessagesClient.On("GetVersionMessage", mock.Anything, mock.Anything).Return(tt.mock.message, nil)
			messager := &Messager{
				defaultTimeout: tt.mock.timeout,
				messagesClient: mockedMessagesClient,
			}
			messager.LoadVersionMessages(tt.args.messagesChan, tt.args.cliVersion)

			msg := <-tt.args.messagesChan

			if tt.mock.message != nil {
				assert.Equal(t, tt.mock.message.CliVersion, msg.CliVersion)
				assert.Equal(t, tt.mock.message.MessageColor, msg.MessageColor)
				assert.Equal(t, tt.mock.message.MessageText, msg.MessageText)
				mockedMessagesClient.AssertCalled(t, "GetVersionMessage", tt.args.cliVersion, tt.mock.timeout)
			}
		})
	}

}

func load_message_success() *LoadVersionMessagesTestCase {
	return &LoadVersionMessagesTestCase{
		name: "should load mock message into channel",
		mock: &LoadVersionMessagesTestCaseMock{
			message: &cliClient.VersionMessage{
				CliVersion:   "0.1.1",
				MessageText:  "message",
				MessageColor: "yellow",
			},
			timeout: 900,
		},
		args: &LoadVersionMessagesTestCaseArgs{
			messagesChan: make(chan *VersionMessage),
			cliVersion:   "0.1.2",
		},
		expected: &LoadVersionMessagesTestCaseExpected{
			message: &VersionMessage{
				CliVersion:   "0.1.1",
				MessageText:  "message",
				MessageColor: "yellow",
			},
		},
	}
}

func load_message_failed() *LoadVersionMessagesTestCase {
	return &LoadVersionMessagesTestCase{
		name: "should load mock message into channel",
		mock: &LoadVersionMessagesTestCaseMock{
			message: nil,
			timeout: 900,
		},
		args: &LoadVersionMessagesTestCaseArgs{
			messagesChan: make(chan *VersionMessage),
			cliVersion:   "0.1.2",
		},
		expected: &LoadVersionMessagesTestCaseExpected{
			message: nil,
		},
	}
}
