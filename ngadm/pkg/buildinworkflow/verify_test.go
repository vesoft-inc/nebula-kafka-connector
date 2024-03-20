package buildinworkflow_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/buildinworkflow"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"gopkg.in/yaml.v3"
)

func TestVerify(t *testing.T) {
	tasks.Init()
	spec := GetNebulaYaml(t)
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
