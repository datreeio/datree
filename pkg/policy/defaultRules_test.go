package policy

import (
    "testing"
    "fmt"
)

type Service struct {
    APIVersion string `yaml:"apiVersion"`
    Kind       string `yaml:"kind"`
    Metadata   struct {
        Name      string `yaml:"name"`
        Namespace string `yaml:"namespace"`
        Labels    struct {
            RouterDeisIoRoutable string `yaml:"router.deis.io/routable"`
        } `yaml:"labels"`
        Annotations struct {
            RouterDeisIoDomains string `yaml:"router.deis.io/domains"`
        } `yaml:"annotations"`
    } `yaml:"metadata"`
    Spec struct {
        Type     string `yaml:"type"`
        Selector struct {
            App string `yaml:"app"`
        } `yaml:"selector"`
        Ports []struct {
            Name       string `yaml:"name"`
            Port       int    `yaml:"port"`
            TargetPort int    `yaml:"targetPort"`
            NodePort   int    `yaml:"nodePort,omitempty"`
        } `yaml:"ports"`
    } `yaml:"spec"`
}


func TestDefaultRulesSyntax(t *testing.T) {
    var service Service
    file = yaml.Unmarshal('./defaultRules.yaml', &service)
 if err != nil {
    panic(err)
}

fmt.Print(service.Metadata.Name)
}
