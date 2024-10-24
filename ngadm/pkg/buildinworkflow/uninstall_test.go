package buildinworkflow_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/yamlparser"
	"gopkg.in/yaml.v3"
)

func TestUninstall(t *testing.T) {
	spec := GetNebulaYaml(t)
	job := runner.NewJob("test uninstall")
	err := job.Run("uninstall", map[string]any{
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
func TestUninstallServiceGroup(t *testing.T) {
	spec := GetNebulaYaml(t)
	// spec.Spec.Metad.ServiceGroups = []types.ServiceGroup{} // remove all clusters for only uninstall metad
	delete(spec.UtilsProcesses, "license-manager")
	job := runner.NewJob("test uninstall")
	err := job.Run("uninstall", map[string]any{
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
	spec := GetNebulaYaml(t)
	spec.Spec.Metad = nil
	job := runner.NewJob("test uninstall")
	err := job.Run("uninstall", map[string]any{
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
	spec := GetNebulaYaml(t)
	spec.Spec.Metad = nil
	lm, ok := spec.UtilsProcesses["LicenseManager"]
	if !ok {
		log.Fatal("LicenseManager not found")
	}
	lm.StartType = "systemd"
	job := runner.NewJob("test uninstall")
	err := job.Run("uninstall", map[string]any{
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

func TestUninstallAgent(t *testing.T) {
	spec, err := yamlparser.ParseYamlByPath("../../examples/agent.yaml")
	if err != nil {
		t.Error(err)
	}
	spec.CAFile = "../../certs/ca.crt"
	spec.KeyFile = "../../certs/client.key"
	spec.CertFile = "../../certs/client.crt"
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
}
