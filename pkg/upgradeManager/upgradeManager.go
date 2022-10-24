package upgrademanager

import (
	"os"
	"os/exec"
	"strings"
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

func (m *UpgradeManager) CheckIfOsIsWindows() bool {
	return strings.Contains("windowssdsd", "windows")
}

func (m *UpgradeManager) Upgrade() error {
	oneLineInstallationCommand := exec.Command("bash", "-c", "curl https://get.datree.io | /bin/bash")
	oneLineInstallationCommand.Stdout = os.Stdout
	return oneLineInstallationCommand.Run()
}
