package buildinworkflow

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

func Uninstall(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "parallel",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	uninstallTask, err := UninstallCluster(args, spec)
	if err != nil {
		return nil, err
	}
	if uninstallTask != nil {
		workflow.Tasks = append(workflow.Tasks, uninstallTask)
	}
	return workflow, nil
}

func UninstallCluster(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
	if args == nil {
		args = map[string]any{}
	}
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
			Type: "serial",
			SubTasks: []*types.TaskSpec{
				{
					Type:   "connect",
					Params: &tasks.ConnectParams{Host: agent.Host},
				},
			},
		})
	}
	//2. stop todo: add stop task
	stopTasks := []*types.TaskSpec{}

	return &types.TaskSpec{
		Type: "serial",
		SubTasks: []*types.TaskSpec{
			{
				Type:     "parallel",
				SubTasks: connectTasks,
			},
			{
				Type:     "parallel",
				SubTasks: stopTasks,
			},
		},
	}, nil
}
