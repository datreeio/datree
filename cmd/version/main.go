package version

import (
	"fmt"

	"github.com/datreeio/datree/bl/messager"
	"github.com/spf13/cobra"
)

type Messager interface {
	PopulateVersionMessageChan(cliVersion string) <-chan *messager.VersionMessage
	HandleVersionMessage(messageChannel <-chan *messager.VersionMessage)
}
type VersionCommandContext struct {
	CliVersion string
	Messager   Messager
}

func version(ctx *VersionCommandContext) {
	messageChannel := ctx.Messager.PopulateVersionMessageChan(ctx.CliVersion)
	fmt.Println(ctx.CliVersion)
	ctx.Messager.HandleVersionMessage(messageChannel)
}

func NewCommand(ctx *VersionCommandContext) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			version(ctx)
		},
	}
	return versionCmd
}
