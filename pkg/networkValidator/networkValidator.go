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

func (nv *NetworkValidator) IdentifyNetworkError(errStr string) error {
	if errStr == "network error" || strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "no such host") {
		if nv.offlineMode == "fail" {
			return errors.New("Failed since internet connection refused, you can use the following command to set your config to run offline:\ndatree config set offline local")
		}
		nv.isBackendAvailable = false
	}
	return nil
}

func (nv *NetworkValidator) IsLocalMode() bool {
	return !nv.isBackendAvailable && nv.offlineMode == "local"
}
