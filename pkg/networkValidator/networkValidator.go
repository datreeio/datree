package networkValidator

import (
	"strings"
)

type NetworkValidator struct {
	isBackendAvailable bool
}

func NewNetworkValidator() *NetworkValidator {
	return &NetworkValidator{}
}

func (nv *NetworkValidator) SetIsBackendAvailable(errStr string) {
	if strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "ECONNREFUSED") {
		nv.isBackendAvailable = false
	}
}

func (nv *NetworkValidator) IsBackendAvailable() bool {
	return nv.isBackendAvailable
}
