package bl

import (
	"time"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/deploymentConfig"
	"github.com/datreeio/datree/pkg/printer"
)

func PopulateVersionMessageChan(cliVersion string) chan *cliClient.VersionMessage {
	messageChannel := make(chan *cliClient.VersionMessage, 1)
	go func() {
		c := cliClient.NewCliClient(deploymentConfig.URL)
		msg, _ := c.GetVersionMessage(cliVersion)
		if msg != nil {
			messageChannel <- msg
		}
		close(messageChannel)
	}()
	return messageChannel
}

func HandleVersionMessage(messageChannel chan *cliClient.VersionMessage) {
	select {
	case msg := <-messageChannel:
		if msg != nil {
			p := printer.CreateNewPrinter()
			p.PrintVersionMessage(msg.MessageText+"\n", msg.MessageColor)
		}
	case <-time.After(600 * time.Millisecond):
		break
	}
}
