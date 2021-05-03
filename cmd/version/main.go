package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCommand(cliVersion string) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cliVersion)
		},
	}
	return versionCmd
}
