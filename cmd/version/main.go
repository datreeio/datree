package version

import (
	"fmt"

	"github.com/datreeio/datree/bl"
	"github.com/spf13/cobra"
)

type VersionCommandContext struct {
	CliVersion           string
	VersionMessageClient bl.VersionMessageClient
}

func version(ctx *VersionCommandContext) {
	messageChannel := bl.PopulateVersionMessageChan(ctx.VersionMessageClient, ctx.CliVersion)
	fmt.Println(ctx.CliVersion)
	bl.HandleVersionMessage(messageChannel)
}

func NewVersionCommand(ctx *VersionCommandContext) *cobra.Command {
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
