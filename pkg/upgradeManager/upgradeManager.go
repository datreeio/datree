package upgrademanager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/datreeio/datree/pkg/httpClient"
)

type UpgradeManager struct {
	platform   func() string
	newCommand func(string, ...string) *exec.Cmd
	client     *httpClient.Client
}

type Asset struct {
	BrowserDownloadUrl string `json:"browser_download_url"`
}

type LatestReleaseAssetList struct {
	Assets []Asset `json:"assets"`
}

func NewUpgradeManager() *UpgradeManager {
	return &UpgradeManager{
		platform: func() string {
			var arch string
			if runtime.GOARCH == "amd64" {
				arch = "x86_64"
			}
			return fmt.Sprintf("%s_%s", strings.Title(runtime.GOOS), arch)
		},
		newCommand: exec.Command,
		client:     httpClient.NewClient("https://api.github.com", nil),
	}
}

func (m *UpgradeManager) CheckIfDatreeInstalledUsingBrew() bool {
	output, _ := m.newCommand("brew", "list", "datree").CombinedOutput()
	if strings.Contains(string(output), "Error: No such keg") {
		return false
	}
	return true
}

func (m *UpgradeManager) Upgrade() error {

	client := m.client
	var body LatestReleaseAssetList

	response, err := client.Request(http.MethodGet, "/repos/datreeio/datree/releases/latest", nil, map[string]string{})
	if err != nil {
		return err
	}

	json.Unmarshal(response.Body, &body)

	for i := range body.Assets {
		// fmt.Println(i, body.Assets[i].BrowserDownloadUrl, m.platform())
		if strings.Contains(body.Assets[i].BrowserDownloadUrl, m.platform()) {
			// download the binary in the temp dir
			// we're doing slicing from 18 bytes because length of  "https://github.com" is 18
			err = m.downloadAsset(body.Assets[i].BrowserDownloadUrl[18:], "/tmp/datree-latest.zip")
			break
		}
	}

	return nil
}

func (m *UpgradeManager) downloadAsset(url string, destPath string) error {
	client := httpClient.NewClient("https://github.com", nil)
	header := map[string]string{
		"Accept": "application/octect-stream",
	}
	response, err := client.Request(http.MethodGet, url, nil, header)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := f.Write(response.Body)
	if err != nil {
		return err
	}

	if n > 0 {
		fmt.Println("Datree file downlaoded successfully, file location is /tmp/datree-latest.zip")
	}

	return nil
}
