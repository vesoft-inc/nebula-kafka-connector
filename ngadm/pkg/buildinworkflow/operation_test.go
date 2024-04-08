package buildinworkflow_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"gopkg.in/yaml.v3"
)

func TestStop(t *testing.T) {
	tasks.Init()
	spec := GetNebulaYaml(t)
	job := runner.NewJob("test stop")
	err := job.Run("operation", map[string]any{
		"operation": "stop",
		"component": "nebulagraph",
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
	spec := GetNebulaYaml(t)
	job := runner.NewJob("test stop")
	err := job.Run("operation", map[string]any{
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
	spec := GetNebulaYaml(t)
	spec.Rollback = false
	job := runner.NewJob("test start")
	err := job.Run("operation", map[string]any{
		"operation": "start",
		"component": "nebulagraph",
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
	spec := GetNebulaYaml(t)
	job := runner.NewJob("test stop")
	err := job.Run("operation", map[string]any{
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
