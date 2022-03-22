package networkValidator

import (
	"strings"
)

type NetworkValidator struct {
	isBackendAvailable bool
	offlineMode        string
}

func NewNetworkValidator() *NetworkValidator {
	return &NetworkValidator{}
}

func (nv *NetworkValidator) SetOfflineMode(offlineMode string) {
	nv.offlineMode = offlineMode
}

func (nv *NetworkValidator) SetIsBackendAvailable(errStr string) {
	if strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "ECONNREFUSED") {
		nv.isBackendAvailable = false
	}
}

func (nv *NetworkValidator) IsBackendAvailable() bool {
	return nv.isBackendAvailable
}
