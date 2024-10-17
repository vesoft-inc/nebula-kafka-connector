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
	for _, host := range metaHosts {
		allNeedHosts[host.Agent.Host] = &host.Agent
	}
	for _, cluster := range spec.Spec.Metad.ServiceGroups {
		for _, host := range cluster.Graphd.Hosts {
			allNeedHosts[host.Agent.Host] = &host.Agent
		}
		for _, host := range cluster.Storaged.Hosts {
			allNeedHosts[host.Agent.Host] = &host.Agent
		}
	}

	connectTasks := []*types.TaskSpec{}
	for _, agent := range allNeedHosts {
		//1. connect
		connectTasks = append(connectTasks, &types.TaskSpec{
			Type:   "connect",
			Params: &tasks.ConnectParams{Host: agent.Host, SSHConfig: agent.SSHConfig},
		})
	}
	return &types.TaskSpec{
		Type:     "parallel",
		SubTasks: connectTasks,
	}, nil
}
