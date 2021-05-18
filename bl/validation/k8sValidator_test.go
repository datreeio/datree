package validation

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	kubeconformValidator "github.com/yannh/kubeconform/pkg/validator"
)

type mockValidationClient struct {
	mock.Mock
}

func (m *mockValidationClient) Validate(filename string, r io.ReadCloser) []kubeconformValidator.Result {
	args := m.Called(filename, r)
	return args.Get(0).([]kubeconformValidator.Result)
}

func TestValidateResources(t *testing.T) {
	validationClient := &mockValidationClient{}
	validationClient.On("Validate", mock.Anything, mock.Anything).Return([]kubeconformValidator.Result{
		{Status: kubeconformValidator.Valid},
	})
	k8sValidator := K8sValidator{
		validationClient: validationClient,
	}

	path := "../../internal/fixtures/kube/pass-all.yaml"
	valid, _, errors := k8sValidator.ValidateResources([]string{path})

	for p := range valid {
		assert.Equal(t, "/Users/noaabarki/Dev/datree/internal/fixtures/kube/pass-all.yaml", p)
	}
	assert.Equal(t, nil, <-errors)
}
