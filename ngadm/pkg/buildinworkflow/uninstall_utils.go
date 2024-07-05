package buildinworkflow

import (
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

// todo: support remove specific util
func UninstallUtils(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
	allUtils := utils.GetAllUtilsProcess(spec)
	//1. connect all need hosts
	connectTasks := []*types.TaskSpec{}
	for _, process := range allUtils {
		for _, agent := range process.Hosts {
			connectTasks = append(connectTasks, &types.TaskSpec{
				Type: "connect",
				Params: &tasks.ConnectParams{
					Host:      agent.Agent.Host,
					SSHConfig: agent.Agent.SSHConfig,
				},
			})
		}
	}
	//2. stop and uninstall
	stopTasks := []*types.TaskSpec{}
	for _, process := range allUtils {
		for _, host := range process.Hosts {
			task := &types.TaskSpec{
				Type:     "serial",
				SubTasks: []*types.TaskSpec{},
			}
			installPath := utils.GetUserUtilPath(spec.InstallPath, host.Agent.InstallPath, process.Name)
			if process.StartType == "shell" {
				scriptPath := path.Join(installPath, process.ExecShellPath)
				task.SubTasks = append(task.SubTasks, &types.TaskSpec{
					Type: "operate",
					Params: &tasks.OperateParams{
						Host:      host.Agent.Host,
						ExecPath:  scriptPath,
						Operation: "stop",
					},
				})
			} else if process.StartType == "systemd" {
				task.SubTasks = append(task.SubTasks, &types.TaskSpec{
					Type: "systemd",
					Params: &tasks.SystemdParams{
						Host:    host.Agent.Host,
						Name:    process.Name,
						Operate: "uninstall",
					},
				})
			}
			task.SubTasks = append(task.SubTasks, &types.TaskSpec{
				Type: "rm",
				Params: &tasks.RMParams{
					Host: host.Agent.Host,
					Path: installPath,
				},
			})
			task.SubTasks = append(task.SubTasks, &types.TaskSpec{
				Type: "rm",
				Params: &tasks.RMParams{
					Host: host.Agent.Host,
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
