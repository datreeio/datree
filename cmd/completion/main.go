package completion

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var out io.Writer = os.Stdout

func New() *cobra.Command {
	var completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script for bash,zsh,fish,powershell",
		Long: `To load completions:

		Bash:
		
		  $ source <(datree completion bash)
		
		  # To load completions for each session, execute once:
		  # Linux:
		  $ datree completion bash > /etc/bash_completion.d/datree
		  # macOS:
		  $ datree completion bash > /usr/local/etc/bash_completion.d/datree
		
		Zsh:
		
		  # If shell completion is not already enabled in your environment,
		  # you will need to enable it.  You can execute the following once:
		
		  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
		
		  # To load completions for each session, execute once:
		  $ datree completion zsh > "${fpath[1]}/_datree"
		
		  # You will need to start a new shell for this setup to take effect.
		
		fish:
		
		  $ datree completion fish | source
		
		  # To load completions for each session, execute once:
		  $ datree completion fish > ~/.config/fish/completions/datree.fish
		
		PowerShell:
		
		  PS> datree completion powershell | Out-String | Invoke-Expression
		
		  # To load completions for every new session, run:
		  PS> datree completion powershell > datree.ps1
		  # and source this file from your PowerShell profile.
		`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				err := cmd.Root().GenBashCompletion(out)
				if err != nil {
					fmt.Printf("Failed to generate bash script. Error : %v", err)
				}
			case "zsh":
				err := cmd.Root().GenZshCompletion(out)
				if err != nil {
					fmt.Printf("Failed to generate zsh script. Error : %v", err)
				}
			case "fish":
				err := cmd.Root().GenFishCompletion(out, true)
				if err != nil {
					fmt.Printf("Failed to generate fish script. Error : %v", err)
				}
			case "powershell":
				err := cmd.Root().GenPowerShellCompletionWithDesc(out)
				if err != nil {
					fmt.Printf("Failed to generate powershell script. Error : %v", err)
				}
			}
		},
	}
	return completionCmd
}
