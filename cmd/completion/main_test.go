package completion

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func Test_ExecuteCommand(t *testing.T) {

	outBuff := new(bytes.Buffer)
	var tt string
	cmd := New()

	b := bytes.NewBufferString("bash") // repalce bash with zsh, fish, powershell for different shell
	cmd.SetOut(b)
	cmd.SetArgs([]string{"bash"}) // repalce bash with zsh, fish, powershell for different shell
	cmd.Run = func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(outBuff)
			tt = outBuff.String()
			if !strings.Contains(tt, "# bash completion for completion") {
				t.Fatalf("Failed: expected \"%s\" got \"%s\"", "# bash completion for completion", tt[:35])
			} else {
				t.Logf("Success: expected \"%s\" got \"%s\"", "# bash completion for completion", tt[:35])
			}
		case "zsh":
			cmd.Root().GenBashCompletion(outBuff)
			tt = outBuff.String()
			if !strings.Contains(tt, "# zsh completion for completion") {
				t.Fatalf("Failed: expected \"%s\" got \"%s\"", "# zsh completion for completion", tt[:35])
			} else {
				t.Logf("Success: expected \"%s\" got \"%s\"", "# zsh completion for completion", tt[:35])
			}
		case "fish":
			cmd.Root().GenBashCompletion(outBuff)
			tt = outBuff.String()
			if !strings.Contains(tt, "# fish completion for completion") {
				t.Fatalf("Failed: expected \"%s\" got \"%s\"", "# fish completion for completion", tt[:35])
			} else {
				t.Logf("Success: expected \"%s\" got \"%s\"", "# fish completion for completion", tt[:35])
			}
		case "powershell":
			cmd.Root().GenBashCompletion(outBuff)
			tt = outBuff.String()
			if !strings.Contains(tt, "# powershell completion for completion") {
				t.Fatalf("Failed: expected \"%s\" got \"%s\"", "# powershell completion for completion", tt[:35])
			} else {
				t.Logf("Success: expected \"%s\" got \"%s\"", "# powershell completion for completion", tt[:35])
			}
		default:
			t.Fatalf("Failed: expected \"%s\" got \"%s\"", "[bash|zsh|fish|powershell]", b.String())
		}

	}
	cmd.Execute()
}
