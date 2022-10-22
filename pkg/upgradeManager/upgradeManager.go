package upgrademanager

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/datreeio/datree/pkg/httpClient"
)

type UpgradeManager struct {
	newCommand func(string, ...string) *exec.Cmd
	client     *httpClient.Client
}

func NewUpgradeManager() *UpgradeManager {
	return &UpgradeManager{
		newCommand: exec.Command,
		client:     httpClient.NewClient("https://raw.githubusercontent.com", nil),
	}
}

func (m *UpgradeManager) CheckIfDatreeInstalledUsingBrew() bool {
	_, err := m.newCommand("brew", "list", "datree").CombinedOutput()
	return err == nil
}

func (m *UpgradeManager) Upgrade() error {

	client := m.client

	// fetch the content of the shell script i.e. `install.sh`
	response, err := client.Request(http.MethodGet, "/datreeio/datree/main/install.sh", nil, map[string]string{})
	if err != nil {
		return err
	}

	destPath, err := os.CreateTemp(os.TempDir(), "datree-install-*.sh")
	// resource cleanup if job complete successfully
	defer func() {
		os.RemoveAll(destPath.Name())
	}()
	if err != nil {
		return err
	}

	_, err = destPath.Write(response.Body)
	if err != nil {
		return err
	}
	// closing the file after writing content to it
	destPath.Close()

	// executing the shell script
	_, err = m.newCommand("bash", destPath.Name()).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
