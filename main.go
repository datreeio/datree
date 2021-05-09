package main

import (
	"os"

	"github.com/datreeio/datree/cmd"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/deploymentConfig"
	"github.com/datreeio/datree/pkg/printer"
)

func main() {
	c := cliClient.NewCliClient(deploymentConfig.URL)
	p := printer.CreateNewPrinter()

	messageChannel := make(chan cliClient.VersionMessage, 1)
	go c.GetVersionMessage(messageChannel, cmd.CliVersion)
	err := cmd.Execute()
	msg := <-messageChannel
	if msg.MessageText != "" {
		p.PrintVersionMessage(msg.MessageText+"\n", msg.MessageColor)
	}
	if err != nil {
		os.Exit(1)
	}
}
