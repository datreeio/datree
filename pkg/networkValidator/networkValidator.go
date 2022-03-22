package networkValidator

import (
	"errors"
	"strings"
)

type NetworkValidator struct {
	isBackendAvailable bool
	offlineMode        string
}

func NewNetworkValidator() *NetworkValidator {
	return &NetworkValidator{
		isBackendAvailable: true,
		offlineMode:        "fail",
	}
}

func (nv *NetworkValidator) SetOfflineMode(offlineMode string) {
	nv.offlineMode = offlineMode
}

func (nv *NetworkValidator) SetIsBackendAvailable(errStr string) error {
	if strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "ECONNREFUSED") {
		if nv.offlineMode == "fail" {
			return errors.New("connection refused and offline mode is on fail")
		}
		nv.isBackendAvailable = false
	}
	return nil
}

func (nv *NetworkValidator) IsLocalMode() bool {
	return !nv.isBackendAvailable && nv.offlineMode == "local"
}
