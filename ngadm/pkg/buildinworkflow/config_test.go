package buildinworkflow_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/buildinworkflow"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/yamlparser"
	"gopkg.in/yaml.v3"
)

func GetNebulaYaml(t *testing.T) *types.JobSpec {
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	spec.CAFile = "../../certs/ca.crt"
	spec.KeyFile = "../../certs/client.key"
	spec.CertFile = "../../certs/client.crt"
	return spec
}

func TestConfig(t *testing.T) {
	tasks.Init()
	args := map[string]interface{}{
		"force": true,
	}
	spec := GetNebulaYaml(t)
	spec.Spec.Metad.Config["test_metad"] = time.Now().String()
	spec.Spec.Metad.ServiceGroups[0].Graphd.Config["test_graphd"] = time.Now().String()
	spec.Spec.Metad.ServiceGroups[0].Storaged.Config["test_storaged"] = time.Now().String()
	lm := spec.UtilsProcesses["license-manager"]
	lm.Config["test_lm"] = time.Now().String()
	workflow, err := buildinworkflow.Config(args, spec)
	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	yamls, err := yaml.Marshal(workflow)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(yamls))
	job := runner.NewJob("test")
	err = job.Run("config", args, spec)
	assert.NoError(t, err)
	// Add more assertions for the workflow.Tasks if needed
}
