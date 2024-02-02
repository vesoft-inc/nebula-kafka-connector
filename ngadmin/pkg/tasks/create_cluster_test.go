package tasks_test

import (
	"testing"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
)

func TestCreateCluster(t *testing.T) {
	tasks.Init()
	executor.SetCertPath("../../certs")

	spec, err := yamlparser.ParseYamlByPath("../../examples/nebula.yaml")
	if err != nil {
		t.Error(err)
	}
	spec.Rollback = false
	spec.Spec.Metad.PackagePath = "../../bin/nebula-graph-5.0-x86_64-glibc-2.17.sh"

	job := runner.NewJob("test")
	task, err := tasks.NewCreateCluster(&types.TaskSpec{
		Type: "create_cluster",
		Params: &tasks.CreateClusterParams{
			ClusterSpec: &spec.Spec.Metad.Clusters[0],
			MetaSpec:    spec.Spec.Metad,
		},
	}, job.Context)
	if err != nil {
		t.Error(err)
	}
	err = task.Execute()
	if err != nil {
		t.Error(err)
	}
	// err = task.Rollback()
	// if err != nil {
	// 	t.Error(err)
	// }
}
