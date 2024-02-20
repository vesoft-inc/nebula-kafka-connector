package buildinworkflow

import (
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/utils"
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
				// TODO: start systemd
				// task.SubTasks = append(task.SubTasks, &types.TaskSpec{
				// 	Type: "systemd",
				// 	Params: &tasks.SystemdParams{
				// 		Host: agent.Host,
				// 		Name: process.Name,
				// 		Cmd:  "start",
				// 	},
				// })
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

func GetAllUtilsProcess(spec *types.JobSpec) []*types.Process {
	allUtils := []*types.Process{}
	if spec.Spec.LicenseManager != nil {
		spec.Spec.LicenseManager.Name = "license-manager"
		allUtils = append(allUtils, spec.Spec.LicenseManager)
	}
	return allUtils
}
