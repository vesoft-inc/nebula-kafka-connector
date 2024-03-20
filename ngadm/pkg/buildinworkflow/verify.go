package buildinworkflow

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

func Verify(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {

	verifyTask, err := VerifyAgents(args, spec)
	if err != nil {
		return nil, err
	}
	workflow := &types.WorkflowSpec{
		Type:        "parallel",
		Description: "verify nebula cluster",
		Rollback:    spec.Rollback,
		Tasks: []*types.TaskSpec{
			verifyTask,
		},
	}
	return workflow, nil
}

func VerifyAgents(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
	if spec.Spec.Metad == nil {
		return nil, fmt.Errorf("metad spec is nil")
	}
	metaHosts := spec.Spec.Metad.Hosts
	allNeedHosts := make(map[string]*types.Agent, 0)
	for _, agent := range metaHosts {
		allNeedHosts[agent.Host] = &agent
	}
	for _, cluster := range spec.Spec.Metad.Clusters {
		for _, agent := range cluster.Graphd.Hosts {
			allNeedHosts[agent.Host] = &agent
		}
		for _, agent := range cluster.Storaged.Hosts {
			allNeedHosts[agent.Host] = &agent
		}
	}

	connectTasks := []*types.TaskSpec{}
	for _, agent := range allNeedHosts {
		//1. connect
		connectTasks = append(connectTasks, &types.TaskSpec{
			Type:   "connect",
			Params: &tasks.ConnectParams{Host: agent.Host},
		})
	}
	return &types.TaskSpec{
		Type:     "parallel",
		SubTasks: connectTasks,
	}, nil
}
