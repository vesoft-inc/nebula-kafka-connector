package buildinworkflow_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"gopkg.in/yaml.v3"
)

func TestUninstall(t *testing.T) {
	tasks.Init()
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
func TestUninstallCluster(t *testing.T) {
	tasks.Init()
	spec := GetNebulaYaml(t)
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
	tasks.Init()
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
	tasks.Init()
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
