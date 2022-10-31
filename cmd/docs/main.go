package docs

import (
	"github.com/datreeio/datree/pkg/utils"
	"github.com/spf13/cobra"
)

type DocsCommandContext struct {
	BrowserCtx utils.OpenBrowserContext
	URL        string
}

func (do *DocsCommandContext) OpenURL(url string) error {
	err := do.BrowserCtx.OpenBrowser(url)
	if err != nil {
		return err
	}
	return nil
}

func New(ctx *DocsCommandContext) *cobra.Command {
	var docsCmd = &cobra.Command{
		Use:     "docs",
		Short:   "Open datree documentation page",
		Long:    "Open the default browser with datree documentation page",
		Aliases: []string{"documentation"},
		Example: utils.Example(`
		# Open documentation with 'docs'
		datree docs

		# Open documentation with alias 'documentation'
		datree documentation
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if ctx.URL == "" {
				ctx.URL = "https://hub.datree.io"
			}
			err := ctx.BrowserCtx.UrlOpener.OpenURL(ctx.URL)
			if err != nil {
				return err
			}
			return nil
		},
	}

	return docsCmd
}
