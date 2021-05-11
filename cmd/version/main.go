package version

import (
	"fmt"

	"github.com/datreeio/datree/bl"
	"github.com/spf13/cobra"
)

type VersionCommandContext struct {
	CliVersion string
	CliClient  bl.VersionMessageClient
}

func NewVersionCommand(ctx *VersionCommandContext) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			messageChannel := bl.PopulateVersionMessageChan(ctx.CliClient, ctx.CliVersion)
			fmt.Println(ctx.CliVersion)
			bl.HandleVersionMessage(messageChannel)
		},
	}
	return versionCmd
}
