package yamlparser_test

import (
	"testing"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/yamlparser"
	"gopkg.in/yaml.v3"
)

func TestParseYamlByPath(t *testing.T) {
	jobSpec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	if jobSpec == nil {
		t.Error("jobSpec is nil")
		return
	}
	if jobSpec.InstallPath == "" {
		t.Error("install path is nil")
	}
	if jobSpec.Spec.Metad == nil {
		t.Error("metad spec is nil")
	}
	if jobSpec.Spec.Metad.Hosts == nil {
		t.Error("metad hosts is nil")
	}
	if jobSpec.Spec.Metad.Clusters == nil {
		t.Error("metad clusters is nil")
	}
	if jobSpec.Spec.Metad.Config == nil {
		t.Error("metad config is nil")
	}
	if _, ok := jobSpec.Spec.Metad.Config["port"]; !ok {
		t.Error("port is nil")
	}
	jsonString, err := yaml.Marshal(jobSpec)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonString))
}
