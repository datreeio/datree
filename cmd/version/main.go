package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCommand(cliVersion string) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("datree version: ", cliVersion)
		},
	}
	return versionCmd
}
