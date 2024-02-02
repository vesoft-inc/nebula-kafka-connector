package buildinworkflow_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/buildinworkflow"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
	"gopkg.in/yaml.v3"
)

func TestVerify(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	workflow, err := buildinworkflow.Verify(nil, spec)
	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	assert.Equal(t, "parallel", workflow.Type)
	yamls, err := yaml.Marshal(workflow)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(yamls))
	job := runner.NewJob("test")
	err = job.Run("verify", nil, spec)
	assert.NoError(t, err)
	// Add more assertions for the workflow.Tasks if needed
}
