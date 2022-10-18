package docs

import (
	"github.com/datreeio/datree/pkg/utils"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var docsCmd = &cobra.Command{
		Use:     "docs",
		Short:   "Datree documentation",
		Long:    "It will open a browser with datree documentation",
		Aliases: []string{"documentation"},
		Example: (`
		# Open documentation with 'docs'
		datree docs

		# Open documentation with alias 'documentation'
		datree documentation
		`),
		Run: func(cmd *cobra.Command, args []string) {
			url := "https://hub.datree.io"
			utils.OpenBrowser(url)
		},
	}

	return docsCmd
}
