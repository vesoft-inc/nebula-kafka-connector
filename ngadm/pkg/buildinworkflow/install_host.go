package buildinworkflow

import (
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func InstallHost(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	if args == nil {
		args = map[string]any{}
	}
	if spec.Spec.Metad == nil {
		return nil, nil
	}
	metaCluster := spec.Spec.Metad
	uploadTasks := []*types.TaskSpec{}
	connectTasks := []*types.TaskSpec{}
	force, ok := args["force"].(bool)
	if !ok {
		force = false
	}
	allNeedHosts := GetMetadAllNeedHosts(spec)
	for _, agent := range allNeedHosts {
		installPath := utils.GetUserClusterPath(spec.InstallPath, agent.InstallPath)
		//1. connect
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
						Path:  installPath,
						Force: force,
					},
				},
			},
		})
		//2.upload
		dstPath := path.Join(utils.GetUserDownloadPath(spec.InstallPath, agent.InstallPath), path.Base(metaCluster.PackagePath))
		uploadTasks = append(uploadTasks, &types.TaskSpec{
			Type: "serial",
			SubTasks: []*types.TaskSpec{
				{
					Type: "upload",
					Params: &tasks.UploadParams{
						SrcPath: utils.GetUserPackagePath(metaCluster.PackagePath, agent.PackagePath),
						DstPath: dstPath,
						Host:    agent.Host,
					},
				},
				{
					Type: "extract",
					Params: &tasks.ExtractParams{
						Host:        agent.Host,
						ExtractPath: installPath,
						PkgPath:     dstPath,
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
