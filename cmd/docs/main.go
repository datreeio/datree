package docs

import (
	"fmt"
	"os"

	"github.com/datreeio/datree/pkg/utils"
	"github.com/spf13/cobra"
)

const DEFAULT_ERR_EXIT_CODE = 1

func New() *cobra.Command {
	var docsCmd = &cobra.Command{
		Use:     "docs",
		Short:   "Open datree documentation page",
		Long:    "Open the default browser with datree documentation page",
		Aliases: []string{"documentation"},
		Example: (`
		# Open documentation with 'docs'
		datree docs

		# Open documentation with alias 'documentation'
		datree documentation
		`),
		Run: func(cmd *cobra.Command, args []string) {
			url := "https://hub.datree.io"
			err := utils.OpenBrowser(url)
			if err != nil {
				fmt.Println(err)
				os.Exit(DEFAULT_ERR_EXIT_CODE)
			}
		},
	}

	return docsCmd
}
