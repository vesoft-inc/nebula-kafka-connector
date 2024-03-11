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

func TestStop(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	job := runner.NewJob("test stop")
	err = job.Run("operation", map[string]any{
		"operation": "stop",
		"component": "all",
		"host":      "",
		"kill-wait": "",
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

func TestStopUtil(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	job := runner.NewJob("test stop")
	err = job.Run("operation", map[string]any{
		"operation": "stop",
		"component": "license-manager",
		"host":      "",
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

func TestStart(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	spec.Rollback = false
	job := runner.NewJob("test start")
	err = job.Run("operation", map[string]any{
		"operation": "start",
		"component": "all",
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

func TestStopWithHostComponent(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	job := runner.NewJob("test stop")
	err = job.Run("operation", map[string]any{
		"operation": "stop",
		"component": "graphd",
		"host":      "192.168.8.240",
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
