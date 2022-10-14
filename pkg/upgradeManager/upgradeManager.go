package upgrademanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/datreeio/datree/pkg/httpClient"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
			return fmt.Sprintf("%s_%s", cases.Title(language.English, cases.Compact).String(runtime.GOOS), arch)
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

	err = json.Unmarshal(response.Body, &body)
	if err != nil {
		return err
	}

	for i := range body.Assets {
		if strings.Contains(body.Assets[i].BrowserDownloadUrl, m.platform()) {
			// download the binary in the temp dir
			destPath, err := os.CreateTemp(os.TempDir(), "datree-latest-*.zip")
			defer os.RemoveAll(destPath.Name())
			if err != nil {
				return err
			}
			// we're doing slicing from 18 bytes because length of  "https://github.com" is 18
			err = m.downloadAsset(body.Assets[i].BrowserDownloadUrl[18:], destPath)
			if err != nil {
				return err
			}

			// unzip the asset file and extract the zip content
			output, err := m.newCommand("unzip", "-d", "/usr/local/bin", "-o", destPath.Name(), "datree").CombinedOutput()
			if err != nil {
				return errors.New(string(output))
			}

			return nil
		}
	}

	return errors.New("Looks like nothing happened, weird! Reach out to datree for support")
}

func (m *UpgradeManager) downloadAsset(url string, destPath *os.File) error {
	client := httpClient.NewClient("https://github.com", nil)
	header := map[string]string{
		"Accept": "application/octect-stream",
	}

	response, err := client.Request(http.MethodGet, url, nil, header)
	if err != nil {
		return err
	}

	defer destPath.Close()

	_, err = destPath.Write(response.Body)
	if err != nil {
		return err
	}

	return nil
}
