package evaluation

import "github.com/shirou/gopsutil/host"

type OSInfo struct {
	OS              string
	PlatformVersion string
	KernelVersion   string
}

func NewOSInfo() *OSInfo {
	infoStat, _ := host.Info()
	return &OSInfo{
		OS:              infoStat.OS,
		KernelVersion:   infoStat.KernelVersion,
		PlatformVersion: infoStat.PlatformVersion,
	}
}
