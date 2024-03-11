package buildinworkflow_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/cmd"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
	"gopkg.in/yaml.v3"
)

func TestStatus(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	job := runner.NewJob("test stop")
	err = job.Run("status", map[string]any{
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
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	spec.Spec.Metad = nil
	job := runner.NewJob("test stop")
	err = job.Run("status", map[string]any{
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
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	lm, ok := spec.UtilsProcesses["license-manager"]
	if !ok {
		log.Fatal("LicenseManager not found")
	}
	lm.StartType = "systemd"
	job := runner.NewJob("test stop")
	err = job.Run("status", map[string]any{
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
