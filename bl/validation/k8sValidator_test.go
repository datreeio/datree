package validation

import (
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/datreeio/datree/pkg/extractor"
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
	test_valid_multiple_configurations(t)
	test_valid_multiple_configurations_only_k8s_files(t)
	test_invalid_file(t)
	test_get_all_schema_locations_online(t)
	test_get_all_schema_locations_offline(t)
	test_get_datree_crd_schema_by_name(t)
	t.Run("test empty file", test_empty_file)
	t.Run("test_offline_with_remote_custom_schema_location", test_offline_with_remote_custom_schema_location)
	t.Run("test missing schema skipped", test_missing_schema_skipped)
	t.Run("test_validateResource_offline_with_local_schema", test_validateResource_offline_with_local_schema)
	t.Run("test_validateResource_offline_without_custom_schema_location", test_validateResource_offline_without_custom_schema_location)
}

func test_valid_multiple_configurations(t *testing.T) {
	validationClient := &mockValidationClient{}
	validationClient.On("Validate", mock.Anything, mock.Anything).Return([]kubeconformValidator.Result{
		{Status: kubeconformValidator.Valid},
	})
	k8sValidator := K8sValidator{
		validationClient: validationClient,
	}

	path := "../../internal/fixtures/kube/pass-all.yaml"

	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)
	filesConfigurationsChan <- &extractor.FileConfigurations{
		FileName:       path,
		Configurations: []extractor.Configuration{},
	}
	close(filesConfigurationsChan)
	validConfigurationsChan, _, _ := k8sValidator.ValidateResources(filesConfigurationsChan, 1)

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

	path := "../../internal/fixtures/kube/Chart.yaml"

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

	path := "../../internal/fixtures/kube/invalidK8sSchema.yaml"

	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)
	filesConfigurationsChan <- &extractor.FileConfigurations{
		FileName:       path,
		Configurations: []extractor.Configuration{},
	}
	close(filesConfigurationsChan)
	_, invalidFilesChan, _ := k8sValidator.ValidateResources(filesConfigurationsChan, 1)

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

	path := "../../internal/fixtures/kube/empty.yaml"

	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)
	filesConfigurationsChan <- &extractor.FileConfigurations{
		FileName:       path,
		Configurations: []extractor.Configuration{},
	}
	close(filesConfigurationsChan)
	_, invalidFilesChan, _ := k8sValidator.ValidateResources(filesConfigurationsChan, 1)

	for p := range invalidFilesChan {
		assert.Equal(t, path, p.Path)
	}
}

func test_offline_with_remote_custom_schema_location(t *testing.T) {
	validationClient := &mockValidationClient{}
	validationClient.On("Validate", mock.Anything, mock.Anything).Return([]kubeconformValidator.Result{
		{Status: kubeconformValidator.Error, Err: fmt.Errorf("no such host")},
	})
	k8sValidator := K8sValidator{
		validationClient:              validationClient,
		areThereCustomSchemaLocations: true,
		isOffline:                     true,
	}

	path := "../../internal/fixtures/kube/pass-all.yaml"

	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)
	filesConfigurationsChan <- &extractor.FileConfigurations{
		FileName:       path,
		Configurations: []extractor.Configuration{},
	}
	close(filesConfigurationsChan)

	_, invalidFilesChan, filesWithWarningsChan := k8sValidator.ValidateResources(filesConfigurationsChan, 1)
	for p := range invalidFilesChan {
		assert.Equal(t, 1, len(p.ValidationErrors))
		assert.Equal(t, "k8s schema validation error: no such host\n", p.ValidationErrors[0].Error())
	}
	for p := range filesWithWarningsChan {
		panic("expected 0 warnings when custom --schema-location provided, instead got warning: " + p.Warning)
	}
}

func test_missing_schema_skipped(t *testing.T) {
	validationClient := &mockValidationClient{}
	validationClient.On("Validate", mock.Anything, mock.Anything).Return([]kubeconformValidator.Result{
		{Status: kubeconformValidator.Skipped, Err: nil},
	})
	k8sValidator := K8sValidator{
		validationClient: validationClient,
	}

	path := "../../internal/fixtures/kube/invalid-kind.yaml"

	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)
	filesConfigurationsChan <- &extractor.FileConfigurations{
		FileName:       path,
		Configurations: []extractor.Configuration{},
	}
	close(filesConfigurationsChan)
	k8sValidationWarningPerValidFile := make(K8sValidationWarningPerValidFile)

	var wg sync.WaitGroup
	filesConfigurationsChanRes, invalidFilesChan, filesWithWarningsChan := k8sValidator.ValidateResources(filesConfigurationsChan, 1)
	wg.Add(1)
	go func() {
		for p := range filesConfigurationsChanRes {
			_ = p
		}
		for p := range invalidFilesChan {
			_ = p
		}
		for p := range filesWithWarningsChan {
			k8sValidationWarningPerValidFile[p.Filename] = *p
		}
		wg.Done()
	}()
	wg.Wait()

	assert.Equal(t, 1, len(k8sValidationWarningPerValidFile))
	assert.Equal(t, "k8s schema validation skipped: --ignore-missing-schemas flag was used", k8sValidationWarningPerValidFile[path].Warning)
}

func test_get_all_schema_locations_online(t *testing.T) {
	expectedOutput := []string{
		"/my-local-schema-location",
		"default",
		"https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{ .NormalizedKubernetesVersion }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
		"https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/argo/{{ .ResourceKind }}_{{ .ResourceAPIVersion }}.json",
	}
	actual := getAllSchemaLocations([]string{"/my-local-schema-location"}, false)
	assert.Equal(t, expectedOutput, actual)
}

func test_get_all_schema_locations_offline(t *testing.T) {
	expectedOutput := []string{
		"/my-local-schema-location",
	}
	actual := getAllSchemaLocations([]string{"/my-local-schema-location"}, true)
	assert.Equal(t, expectedOutput, actual)
}

func test_get_datree_crd_schema_by_name(t *testing.T) {
	input := "argo"
	expectedOutput := "https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/argo/{{ .ResourceKind }}_{{ .ResourceAPIVersion }}.json"
	actual := getDatreeCRDSchemaByName(input)

	if actual != expectedOutput {
		t.Errorf("Expected: %s, Actual: %s", expectedOutput, actual)
	}
}

func test_validateResource_offline_with_local_schema(t *testing.T) {
	k8sValidator := &K8sValidator{
		validationClient: newKubeconformValidator("1.21.0", false, getAllSchemaLocations([]string{
			"some-path-to-non-existing-file-to-get-404.yaml",
		}, true)),
		isOffline:                     true,
		areThereCustomSchemaLocations: true,
	}

	isValid, validationErrors, validationWarningResult, err := k8sValidator.validateResource("../../internal/fixtures/kube/pass-all.yaml")
	var nilValidationWarning *validationWarning
	assert.Equal(t, nil, err)
	assert.Equal(t, false, isValid)
	assert.Equal(t, "k8s schema validation error: could not find schema for Deployment\nYou can skip files with missing schemas instead of failing by using the `--ignore-missing-schemas` flag\n", validationErrors[0].Error())
	assert.Equal(t, nilValidationWarning, validationWarningResult)
}

func test_validateResource_offline_without_custom_schema_location(t *testing.T) {
	k8sValidator := &K8sValidator{
		validationClient:              newKubeconformValidator("1.21.0", false, getAllSchemaLocations([]string{}, true)),
		isOffline:                     true,
		areThereCustomSchemaLocations: false,
	}

	isValid, validationErrors, validationWarningResult, err := k8sValidator.validateResource("../../internal/fixtures/kube/pass-all.yaml")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, isValid)
	assert.Equal(t, 0, len(validationErrors))
	assert.Equal(t, &validationWarning{
		WarningKind:    NetworkError,
		WarningMessage: "k8s schema validation skipped: no internet connection",
	}, validationWarningResult)
}
