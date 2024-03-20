package buildinworkflow

import (
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func InstallUtils(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
	allUtils := GetAllUtilsProcess(spec)
	//1. connect all need hosts
	connectTasks := []*types.TaskSpec{}
	for _, process := range allUtils {
		for _, agent := range process.Hosts {
			connectTasks = append(connectTasks, &types.TaskSpec{
				Type: "connect",
				Params: &tasks.ConnectParams{
					Host: agent.Host,
				},
			})
		}
	}
	//2. upload and start
	//2.1 todo: apply utils config
	uploadAndStartTasks := []*types.TaskSpec{}
	for _, process := range allUtils {
		for _, agent := range process.Hosts {
			dstPath := path.Join(utils.GetDownloadPath(spec.InstallPath), path.Base(process.PackagePath))
			task := &types.TaskSpec{
				Type: "serial",
				SubTasks: []*types.TaskSpec{
					{
						Type: "upload",
						Params: &tasks.UploadParams{
							SrcPath: process.PackagePath,
							DstPath: dstPath,
							Host:    agent.Host,
						},
					},
					{
						Type: "extract",
						Params: &tasks.ExtractParams{
							Host:        agent.Host,
							ExtractPath: utils.GetUtilPath(spec.InstallPath, process.Name),
							PkgPath:     dstPath,
						},
					},
					{
						Type: "config",
						Params: &tasks.ConfigParams{
							Host:   agent.Host,
							Config: process.Config,
							Dst:    path.Join(utils.GetUtilPath(spec.InstallPath, process.Name), process.ConfigPath),
							Type:   "yaml",
						},
					},
				},
			}
			if process.StartType == "shell" {
				scriptPath := path.Join(utils.GetUtilPath(spec.InstallPath, process.Name), process.ExecShellPath)
				task.SubTasks = append(task.SubTasks, &types.TaskSpec{
					Type: "operate",
					Params: &tasks.OperateParams{
						Host:      agent.Host,
						ExecPath:  scriptPath,
						Operation: "start",
					},
				})
			} else if process.StartType == "systemd" {
				task.SubTasks = append(task.SubTasks, &types.TaskSpec{
					Type: "systemd",
					Params: &tasks.SystemdParams{
						Host:             agent.Host,
						Name:             process.Name,
						ExecStartPath:    path.Join(utils.GetUtilPath(spec.InstallPath, process.Name), process.ExecStartPath),
						WorkingDirectory: path.Join(utils.GetUtilPath(spec.InstallPath, process.Name), process.WorkingDir),
						Operate:          "install",
					},
				})
			}
			uploadAndStartTasks = append(uploadAndStartTasks, task)
		}
	}

	return &types.TaskSpec{
		Type: "serial",
		SubTasks: []*types.TaskSpec{
			{
				Type:     "parallel",
				SubTasks: connectTasks,
			},
			{
				Type:        "parallel",
				Description: "upload and start utils",
				SubTasks:    uploadAndStartTasks,
			},
		},
	}, nil
}

// map component name
func GetAllUtilsProcess(spec *types.JobSpec) []*types.Process {
	allUtils := []*types.Process{}

	for _, process := range spec.UtilsProcesses {
		if process.Config == nil {
			process.Config = map[string]any{}
		}
		allUtils = append(allUtils, process)
	}
	return allUtils
}
