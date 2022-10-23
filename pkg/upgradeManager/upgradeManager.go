package upgrademanager

import (
	"os"
	"os/exec"
)

type UpgradeManager struct {
}

func NewUpgradeManager() *UpgradeManager {
	return &UpgradeManager{}
}

func (m *UpgradeManager) CheckIfDatreeInstalledUsingBrew() bool {
	_, err := exec.Command("brew", "list", "datree").Output()
	return err == nil
}

func (m *UpgradeManager) Upgrade() error {

	shellScript := exec.Command("curl", "https://get.datree.io")
	execShellScript := exec.Command("bash")
	execShellScript.Stdin, _ = shellScript.StdoutPipe()
	execShellScript.Stdout = os.Stdout

	err := execShellScript.Start()
	if err != nil {
		return err
	}

	err = shellScript.Run()
	if err != nil {
		return err
	}

	err = execShellScript.Wait()
	if err != nil {
		return err
	}

	return nil
}
