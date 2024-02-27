package buildinworkflow_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
	"gopkg.in/yaml.v3"
)

func TestUninstall(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	job := runner.NewJob("test uninstall")
	err = job.Run("uninstall", map[string]any{
		"drain":     true,
		"kill-wait": "10s",
	}, spec)
	if err != nil {
		yamls, err := yaml.Marshal(job.WorkflowSpec)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(yamls))
	}
	assert.NoError(t, err)
}
func TestUninstallCluster(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/cluster.yaml")
	if err != nil {
		t.Error(err)
	}
	job := runner.NewJob("test uninstall")
	err = job.Run("uninstall", map[string]any{
		"drain":     true,
		"kill-wait": "10s",
	}, spec)
	if err != nil {
		yamls, err := yaml.Marshal(job.WorkflowSpec)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(yamls))
	}
	assert.NoError(t, err)
}

func TestUninstallUtils(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	spec.Spec.Metad = nil
	job := runner.NewJob("test uninstall")
	err = job.Run("uninstall", map[string]any{
		"drain":     true,
		"kill-wait": "10s",
	}, spec)
	if err != nil {
		yamls, err := yaml.Marshal(job.WorkflowSpec)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(yamls))
	}
	assert.NoError(t, err)
}

func TestUninstallUtilsSystemd(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	spec.Spec.Metad = nil
	spec.Spec.LicenseManager.StartType = "systemd"
	job := runner.NewJob("test uninstall")
	err = job.Run("uninstall", map[string]any{
		"drain":     true,
		"kill-wait": "10s",
	}, spec)
	if err != nil {
		yamls, err := yaml.Marshal(job.WorkflowSpec)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(yamls))
	}
	assert.NoError(t, err)
}
