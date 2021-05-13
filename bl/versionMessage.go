package bl

import (
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/printer"
)

type VersionMessageClient interface {
	GetVersionMessage(cliVersion string) (*cliClient.VersionMessage, error)
}

func PopulateVersionMessageChan(c VersionMessageClient, cliVersion string) <-chan *cliClient.VersionMessage {
	messageChannel := make(chan *cliClient.VersionMessage)
	go func() {
		msg, _ := c.GetVersionMessage(cliVersion)
		if msg != nil {
			messageChannel <- msg
		}
		close(messageChannel)
	}()

	return messageChannel
}

func HandleVersionMessage(messageChannel <-chan *cliClient.VersionMessage) {
	msg, ok := <-messageChannel
	if ok {
		p := printer.CreateNewPrinter()
		p.PrintVersionMessage(msg.MessageText+"\n", msg.MessageColor)
	}
}
