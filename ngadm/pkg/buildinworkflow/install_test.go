package buildinworkflow_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/buildinworkflow"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/yamlparser"
	"gopkg.in/yaml.v3"
)

func TestInstall(t *testing.T) {
	tasks.Init()
	args := map[string]interface{}{
		"force": true,
	}
	spec := GetNebulaYaml(t)
	spec.Rollback = false
	spec.Spec.Metad.PackagePath = "../../bin/nebula-graph-5.0-x86_64-glibc-2.17.sh"
	lm, ok := spec.UtilsProcesses["LicenseManager"]
	if ok {
		lm.PackagePath = "../../bin/lm.tar.gz"
	}
	workflow, err := buildinworkflow.Install(args, spec)
	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	yamls, err := yaml.Marshal(workflow)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(yamls))
	job := runner.NewJob("test")
	err = job.Run("install", args, spec)
	assert.NoError(t, err)
	// Add more assertions for the workflow.Tasks if needed
}

func TestClusterInstall(t *testing.T) {
	tasks.Init()

	args := map[string]interface{}{
		"force": true,
	}
	spec := GetNebulaYaml(t)
	delete(spec.UtilsProcesses, "license-manager")
	spec.Rollback = true
	spec.Spec.Metad.PackagePath = "../../bin/nebula-graph-5.0-x86_64-glibc-2.17.sh"
	// spec.Spec.Metad.Clusters = []types.Cluster{}
	workflow, err := buildinworkflow.Install(args, spec)
	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	yamls, err := yaml.Marshal(workflow)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(yamls))
	job := runner.NewJob("test")
	err = job.Run("install", args, spec)
	assert.NoError(t, err)
	// Add more assertions for the workflow.Tasks if needed
}

func TestUtilsInstall(t *testing.T) {
	tasks.Init()

	args := map[string]interface{}{
		"force": true,
	}
	spec := GetNebulaYaml(t)
	spec.Rollback = false
	spec.Spec.Metad = nil
	lm := spec.UtilsProcesses["license-manager"]
	lm.PackagePath = "../../bin/lm.tar.gz"
	workflow, err := buildinworkflow.Install(args, spec)
	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	yamls, err := yaml.Marshal(workflow)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(yamls))
	job := runner.NewJob("test")
	err = job.Run("install", args, spec)
	assert.NoError(t, err)
	// Add more assertions for the workflow.Tasks if needed
}

func TestUtilsInstallSystemd(t *testing.T) {
	tasks.Init()

	args := map[string]interface{}{
		"force": true,
	}
	spec := GetNebulaYaml(t)
	spec.Rollback = false
	spec.Spec.Metad = nil
	lm := spec.UtilsProcesses["LicenseManager"]
	lm.StartType = "systemd"
	lm.PackagePath = "../../bin/lm.tar.gz"
	workflow, err := buildinworkflow.Install(args, spec)
	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	yamls, err := yaml.Marshal(workflow)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(yamls))
	job := runner.NewJob("test")
	err = job.Run("install", args, spec)
	assert.NoError(t, err)
	// Add more assertions for the workflow.Tasks if needed
}

func TestAgentInstallSystemd(t *testing.T) {
	tasks.Init()
	args := map[string]interface{}{
		"force": true,
	}
	spec, err := yamlparser.ParseYamlByPath("../../examples/agent.yaml")
	if err != nil {
		t.Error(err)
	}
	spec.CAFile = "../../certs/ca.crt"
	spec.KeyFile = "../../certs/client.key"
	spec.CertFile = "../../certs/client.crt"
	spec.ServerCertFile = "../../certs/server.crt"
	spec.ServerKeyFile = "../../certs/server.key"
	spec.Rollback = false
	spec.Spec.Metad = nil

	agent := spec.UtilsProcesses["agent"]
	agent.StartType = "shell"
	agent.PackagePath = "../../bin/agent.tar.gz"
	job := runner.NewJob("test")

	err = job.Run("install-agent", args, spec)
	assert.NoError(t, err)
	// Add more assertions for the workflow.Tasks if needed
}
