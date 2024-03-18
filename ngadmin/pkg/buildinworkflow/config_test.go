package buildinworkflow_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/buildinworkflow"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
	"gopkg.in/yaml.v3"
)

func GetNebulaYaml(t *testing.T) *types.JobSpec {
	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	spec.CAFile = "../../certs/ca.crt"
	spec.KeyFile = "../../certs/ngadmin.key"
	spec.CertFile = "../../certs/ngadmin.crt"
	return spec
}

func TestConfig(t *testing.T) {
	tasks.Init()
	args := map[string]interface{}{
		"force": true,
	}
	spec := GetNebulaYaml(t)
	spec.Spec.Metad.Config["test_metad"] = time.Now().String()
	spec.Spec.Metad.Clusters[0].Graphd.Config["test_graphd"] = time.Now().String()
	spec.Spec.Metad.Clusters[0].Storaged.Config["test_storaged"] = time.Now().String()
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
