package buildinworkflow

import (
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

func UploadFile(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: true,
		Tasks:    []*types.TaskSpec{},
	}
	if args == nil {
		args = map[string]any{}
	}
	filePath := args["file_path"].([]string)
	dstPath := args["dst_path"].([]string)
	hosts := args["host"].([]string)
	connectTasks := []*types.TaskSpec{}
	uploadTasks := []*types.TaskSpec{}
	allNeedHosts := GetMetadAllNeedHosts(spec)
	// for i, agent := range allNeedHosts {
	for i, file := range filePath {
		agent := allNeedHosts[hosts[i]]
		if agent == nil {
			agent = &types.Agent{
				Host: hosts[i],
			}
		}
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
					Type: "check_dir",
					Params: &tasks.CheckDirParams{
						Host:  agent.Host,
						Path:  agent.InstallPath,
						Force: true,
					},
				},
			},
		})
		uploadTasks = append(uploadTasks, &types.TaskSpec{
			Type: "upload",
			Params: &tasks.UploadParams{
				SrcPath:  file,
				Host: agent.Host,
				DstPath:  dstPath[i],
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
			{
				Type:        "parallel",
				Description: "upload packages",
				SubTasks:    uploadTasks,
			},
		},
	}
	workflow.Tasks = append(workflow.Tasks, mainTask)
	return workflow, nil
}