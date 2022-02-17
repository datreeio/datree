package validation

import (
	"fmt"
	"io"
	"testing"

	"github.com/datreeio/datree/pkg/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	kubeconformValidator "github.com/yannh/kubeconform/pkg/validator"
)
const FIXTURES_PATH string = "../../fixtures/kube";
type mockValidationClient struct {
	mock.Mock
}

func (m *mockValidationClient) Validate(filename string, r io.ReadCloser) []kubeconformValidator.Result {
	args := m.Called(filename, r)
	return args.Get(0).([]kubeconformValidator.Result)
}

func TestValidateResources(t *testing.T) {
	test_valid_multiple_configurations(t)
	test_valid_multiple_configurations_only_k8s_files(t)
	test_invalid_file(t)
	test_default_schema_location(t)
	test_get_datree_crd_schema_by_name(t)
	t.Run("test empty file", test_empty_file)
}

func test_valid_multiple_configurations(t *testing.T) {
	validationClient := &mockValidationClient{}
	validationClient.On("Validate", mock.Anything, mock.Anything).Return([]kubeconformValidator.Result{
		{Status: kubeconformValidator.Valid},
	})
	k8sValidator := K8sValidator{
		validationClient: validationClient,
	}

	path :=  fmt.Sprintf("%s/pass-all.yaml", FIXTURES_PATH)

	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)
	filesConfigurationsChan <- &extractor.FileConfigurations{
		FileName:       path,
		Configurations: []extractor.Configuration{},
	}
	close(filesConfigurationsChan)
	validConfigurationsChan, _ := k8sValidator.ValidateResources(filesConfigurationsChan, 1)

	for p := range validConfigurationsChan {
		assert.Equal(t, path, p.FileName)
	}
}

func test_valid_multiple_configurations_only_k8s_files(t *testing.T) {
	validationClient := &mockValidationClient{}
	validationClient.On("Validate", mock.Anything, mock.Anything).Return([]kubeconformValidator.Result{
		{Status: kubeconformValidator.Valid},
	})
	k8sValidator := K8sValidator{
		validationClient: validationClient,
	}

	 path  := fmt.Sprintf("%s/Chart.yaml", FIXTURES_PATH)

	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)
	filesConfigurationsChan <- &extractor.FileConfigurations{
		FileName:       path,
		Configurations: []extractor.Configuration{},
	}
	close(filesConfigurationsChan)
	validK8sFilesChan, _ := k8sValidator.GetK8sFiles(filesConfigurationsChan, 1)

	for p := range validK8sFilesChan {
		assert.Equal(t, path, p.FileName)
	}
}

func test_invalid_file(t *testing.T) {
	validationClient := &mockValidationClient{}
	validationClient.On("Validate", mock.Anything, mock.Anything).Return([]kubeconformValidator.Result{
		{Status: kubeconformValidator.Invalid, Err: fmt.Errorf("missing 'apiVersion' key")},
	})
	k8sValidator := K8sValidator{
		validationClient: validationClient,
	}

	path  := fmt.Sprintf("%s/invalidK8sSchema.yaml", FIXTURES_PATH)

	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)
	filesConfigurationsChan <- &extractor.FileConfigurations{
		FileName:       path,
		Configurations: []extractor.Configuration{},
	}
	close(filesConfigurationsChan)
	_, invalidFilesChan := k8sValidator.ValidateResources(filesConfigurationsChan, 1)

	for p := range invalidFilesChan {
		assert.Equal(t, path, p.Path)
	}
}

func test_empty_file(t *testing.T) {
	validationClient := &mockValidationClient{}
	validationClient.On("Validate", mock.Anything, mock.Anything).Return([]kubeconformValidator.Result{
		{Status: kubeconformValidator.Invalid, Err: fmt.Errorf("empty file")},
	})
	k8sValidator := K8sValidator{
		validationClient: validationClient,
	}

	path  := fmt.Sprintf("%s/empty.yaml", FIXTURES_PATH)

	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)
	filesConfigurationsChan <- &extractor.FileConfigurations{
		FileName:       path,
		Configurations: []extractor.Configuration{},
	}
	close(filesConfigurationsChan)
	_, invalidFilesChan := k8sValidator.ValidateResources(filesConfigurationsChan, 1)

	for p := range invalidFilesChan {
		assert.Equal(t, path, p.Path)
	}
}

func test_default_schema_location(t *testing.T) {
	expectedOutput := []string{"default", "https://raw.githubusercontent.com/datreeio/CRDs-catalog/master/argo/{{ .ResourceKind }}_{{ .ResourceAPIVersion }}.json"}
	actual := getDefaultSchemaLocation()
	assert.Equal(t, expectedOutput, actual)
}

func test_get_datree_crd_schema_by_name(t *testing.T) {
	input := "argo"
	expectedOutput := "https://raw.githubusercontent.com/datreeio/CRDs-catalog/master/argo/{{ .ResourceKind }}_{{ .ResourceAPIVersion }}.json"
	actual := getDatreeCRDSchemaByName(input)

	if actual != expectedOutput {
		t.Errorf("Expected: %s, Actual: %s", expectedOutput, actual)
	}
}
