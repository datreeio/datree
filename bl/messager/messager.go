package messager

import (
	"github.com/datreeio/datree/pkg/cliClient"
)

type MessagesClient interface {
	GetVersionMessage(cliVersion string, timeout int) (*cliClient.VersionMessage, error)
}

type Printer interface {
	PrintVersionMessage(messageText string, messageColor string)
}

type Messager struct {
	defaultTimeout int
	messagesClient MessagesClient
	printer        Printer
}

func New(c MessagesClient, p Printer) *Messager {
	return &Messager{
		defaultTimeout: 900,
		messagesClient: c,
		printer:        p,
	}
}

type VersionMessage struct {
	CliVersion   string
	MessageText  string
	MessageColor string
}

func (m *Messager) PopulateVersionMessageChan(cliVersion string) <-chan *VersionMessage {
	messageChannel := make(chan *VersionMessage)
	go func() {
		msg, _ := m.messagesClient.GetVersionMessage(cliVersion, 900)
		if msg != nil {
			messageChannel <- m.toVersionMessage(msg)
		}
		close(messageChannel)
	}()

	return messageChannel
}

func (m *Messager) HandleVersionMessage(messageChannel <-chan *VersionMessage) {
	msg, ok := <-messageChannel
	if ok {
		m.printer.PrintVersionMessage(msg.MessageText+"\n", msg.MessageColor)
	}
}

func (m *Messager) toVersionMessage(msg *cliClient.VersionMessage) *VersionMessage {
	return &VersionMessage{
		CliVersion:   msg.CliVersion,
		MessageText:  msg.MessageText,
		MessageColor: msg.MessageColor,
	}
}
