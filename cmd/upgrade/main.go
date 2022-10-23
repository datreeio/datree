package upgrade

import (
	"fmt"

	"github.com/datreeio/datree/pkg/cliClient"
	upgrademanager "github.com/datreeio/datree/pkg/upgradeManager"
	"github.com/spf13/cobra"
)

type Printer interface {
	PrintMessage(messageText string, messageColor string)
}
type UpgradeCommandContext struct {
	CliVersion       string
	Printer          Printer
	UpgradeCliClient *cliClient.CliClient
}

func New(ctx *UpgradeCommandContext) *cobra.Command {
	m := upgrademanager.NewUpgradeManager()
	var upgradeCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade datree to the latest version",
		Long:  "Upgrade datree to the latest version",
		Run: func(cmd *cobra.Command, args []string) {
			if m.CheckIfDatreeInstalledUsingBrew() {
				ctx.Printer.PrintMessage("Looks like you installed Datree via brew, to upgrade datree run - brew upgrade datree\n", "error")
			} else {
				if m.CheckIfDatreeInstalledUsingBrew() {
					ctx.Printer.PrintMessage("Looks like you installed Datree via brew, to upgrade datree run - brew upgrade datree\n", "error")
				} else {
					err := m.Upgrade()
					if err != nil {
						errMsg := fmt.Sprintf("Failed to upgrade datree to the latest version, error encountered %s\n", err)
						ctx.Printer.PrintMessage(errMsg, "error")
						return
					}
					ctx.Printer.PrintMessage("Datree upgraded successfully\n", "green")
				}
			}
		},
	}
	return upgradeCmd
}