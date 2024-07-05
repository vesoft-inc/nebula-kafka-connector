package buildinworkflow

import (
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func UninstallHost(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	if args == nil {
		args = map[string]any{}
	}
	connectTasks := []*types.TaskSpec{}
	allNeedHosts := GetMetadAllNeedHosts(spec)
	for _, agent := range allNeedHosts {
		installPath := utils.GetUserClusterPath(spec.InstallPath, agent.InstallPath)
		// connect and uninstall
		connectTasks = append(connectTasks, &types.TaskSpec{
			Type: "serial",
			SubTasks: []*types.TaskSpec{
				{
					Type: "connect",
					Params: &tasks.ConnectParams{
						Host:      agent.Host,
						SSHConfig: agent.SSHConfig,
					},
				},
				{
					Type: "rm",
					Params: &tasks.RMParams{
						Host:  agent.Host,
						Path:  installPath,
					},
				},
			},
		})
	}
	mainTask := &types.TaskSpec{
		Type: "serial",
		SubTasks: []*types.TaskSpec{
			{
				Type:     "parallel",
				SubTasks: connectTasks,
			},
		},
	}
	workflow.Tasks = append(workflow.Tasks, mainTask)
	return workflow, nil
}
