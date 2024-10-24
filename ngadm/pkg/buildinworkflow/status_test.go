package buildinworkflow_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/cmd"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"gopkg.in/yaml.v3"
)

func TestStatus(t *testing.T) {
	spec := GetNebulaYaml(t)
	job := runner.NewJob("test stop")
	err := job.Run("status", map[string]any{
		"component": "all",
	}, spec)
	assert.NoError(t, err)
	yamls, err := yaml.Marshal(job.WorkflowSpec)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(yamls))
	assert.NoError(t, err)
	cmd.RenderStatusTableByJob(job)
}

func TestUtilsStatus(t *testing.T) {
	spec := GetNebulaYaml(t)
	spec.Spec.Metad = nil
	job := runner.NewJob("test stop")
	err := job.Run("status", map[string]any{
		"component": "all",
	}, spec)
	assert.NoError(t, err)
	yamls, err := yaml.Marshal(job.WorkflowSpec)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(yamls))
	assert.NoError(t, err)
	cmd.RenderStatusTableByJob(job)
}

func TestSystemdStatus(t *testing.T) {
	spec := GetNebulaYaml(t)
	lm, ok := spec.UtilsProcesses["license-manager"]
	if !ok {
		log.Fatal("LicenseManager not found")
	}
	lm.StartType = "systemd"
	job := runner.NewJob("test stop")
	err := job.Run("status", map[string]any{
		"component": "all",
	}, spec)
	assert.NoError(t, err)
	yamls, err := yaml.Marshal(job.WorkflowSpec)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(yamls))
	assert.NoError(t, err)
	cmd.RenderStatusTableByJob(job)
}
