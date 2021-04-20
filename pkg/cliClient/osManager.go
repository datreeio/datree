package cliClient

import "github.com/shirou/gopsutil/host"

type UserAgent struct {
	OS              string `json:"os"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
}

func getUserAgent() (*UserAgent, error) {
	osInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	return &UserAgent{
		OS:              osInfo.OS,
		PlatformVersion: osInfo.PlatformVersion,
		KernelVersion:   osInfo.KernelVersion,
	}, nil
}
