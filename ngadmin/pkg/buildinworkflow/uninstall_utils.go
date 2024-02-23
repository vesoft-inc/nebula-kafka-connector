package buildinworkflow

import (
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/utils"
)

func UninstallUtils(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
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
	//2. stop and uninstall
	stopTasks := []*types.TaskSpec{}
	for _, process := range allUtils {
		for _, agent := range process.Hosts {
			task := &types.TaskSpec{
				Type:     "serial",
				SubTasks: []*types.TaskSpec{},
			}
			if process.StartType == "shell" {
				scriptPath := path.Join(utils.GetUtilPath(spec.InstallPath, process.Name), process.ExecShellPath)
				task.SubTasks = append(task.SubTasks, &types.TaskSpec{
					Type: "operate",
					Params: &tasks.OperateParams{
						Host:      agent.Host,
						ExecPath:  scriptPath,
						Operation: "stop",
					},
				})
			} else if process.StartType == "systemd" {
				task.SubTasks = append(task.SubTasks, &types.TaskSpec{
					Type: "systemd",
					Params: &tasks.SystemdParams{
						Host:    agent.Host,
						Name:    process.Name,
						Operate: "uninstall",
					},
				})
			}
			task.SubTasks = append(task.SubTasks, &types.TaskSpec{
				Type: "rm",
				Params: &tasks.RMParams{
					Host: agent.Host,
					Path: utils.GetUtilPath(spec.InstallPath, process.Name),
				},
			})
			task.SubTasks = append(task.SubTasks, &types.TaskSpec{
				Type: "rm",
				Params: &tasks.RMParams{
					Host: agent.Host,
					Path: utils.GetDownloadPath(spec.InstallPath),
				},
			})

			stopTasks = append(stopTasks, task)
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
				Description: "stop and uninstall utils",
				SubTasks:    stopTasks,
			},
		},
	}, nil
}
