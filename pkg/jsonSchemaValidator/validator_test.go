package jsonSchemaValidator

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateCustomKeysFail(t *testing.T) {
	resourceJson := `{
  "apiVersion": "v1",
  "kind": "Pod",
  "metadata": {
    "name": "frontend"
  },
  "spec": {
      "containers": [
      {
        "name": "cpu-demo",
        "image": "mages.my-company.example/app:v4",
        "resources": {
          "requests": {
            "memory": "64Mi",
            "cpu": "250m"
          },
          "limits": {
            "memory": "128Mi",
            "cpu": "1G"
          }
        }
      }
    ]
  }
}`
	customRuleSchemaJson := `{
  "properties": {
    "spec": {
      "properties": {
        "containers": {
          "items": {
            "properties": {
              "resources": {
                "properties": {
                  "limits": {
                    "properties": {
                      "cpu": {
                        "resourceMinimum": "250m",
                        "resourceMaximum": "500m"
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`

	jsonSchemaValidator := New()

	resourceYaml, _ := yaml.JSONToYAML([]byte(resourceJson))
	customRuleYaml, _ := yaml.JSONToYAML([]byte(customRuleSchemaJson))

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(string(customRuleYaml), string(resourceYaml))

	fmt.Println(len(errorsResult))
	assert.GreaterOrEqual(t, len(errorsResult), 1)

}

func TestValidateCustomKeysPass(t *testing.T) {
	resourceJson := `{
  "apiVersion": "v1",
  "kind": "Pod",
  "metadata": {
    "name": "frontend"
  },
  "spec": {
      "containers": [
      {
        "name": "cpu-demo",
        "image": "mages.my-company.example/app:v4",
        "resources": {
          "requests": {
            "memory": "64Mi",
            "cpu": "250m"
          },
          "limits": {
            "memory": "128Mi",
            "cpu": "350m"
          }
        }
      }
    ]
  }
}`
	customRuleSchemaJson := `{
  "properties": {
    "spec": {
      "properties": {
        "containers": {
          "items": {
            "properties": {
              "resources": {
                "properties": {
                  "limits": {
                    "properties": {
                      "cpu": {
                        "resourceMinimum": "250m",
                        "resourceMaximum": "500m"
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`

	jsonSchemaValidator := New()

	resourceYaml, _ := yaml.JSONToYAML([]byte(resourceJson))
	customRuleYaml, _ := yaml.JSONToYAML([]byte(customRuleSchemaJson))

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(string(customRuleYaml), string(resourceYaml))

	fmt.Println(len(errorsResult))
	assert.Empty(t, errorsResult)

}
