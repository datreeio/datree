package messager

import (
	"github.com/datreeio/datree/pkg/cliClient"
)

type MessagesClient interface {
	GetVersionMessage(cliVersion string, timeout int) (*cliClient.VersionMessage, error)
}

type Messager struct {
	defaultTimeout int
	messagesClient MessagesClient
}

func New(c MessagesClient) *Messager {
	return &Messager{
		defaultTimeout: 900,
		messagesClient: c,
	}
}

type VersionMessage struct {
	CliVersion   string
	MessageText  string
	MessageColor string
}

func (m *Messager) LoadVersionMessages(cliVersion string) chan *VersionMessage {
	messages := make(chan *VersionMessage, 1)
	go func() {
		msg, _ := m.messagesClient.GetVersionMessage(cliVersion, 900)
		if msg != nil {
			messages <- m.toVersionMessage(msg)
		}
		close(messages)
	}()
	return messages
}

func (m *Messager) toVersionMessage(msg *cliClient.VersionMessage) *VersionMessage {
	if msg != nil {
		return &VersionMessage{
			CliVersion:   msg.CliVersion,
			MessageText:  msg.MessageText,
			MessageColor: msg.MessageColor,
		}
	}

	return nil
}
