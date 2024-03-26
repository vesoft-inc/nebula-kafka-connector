package buildinworkflow

import (
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func StatusUtils(spec *types.JobSpec, component string) (*types.TaskSpec, error) {
	allUtils := utils.GetAllUtilsProcess(spec)
	connectTasks := []*types.TaskSpec{}
	statusTasks := []*types.TaskSpec{}
	for _, process := range allUtils {
		if component == "all" || component == process.Name {
			for _, host := range process.Hosts {
				connectTasks = append(connectTasks, &types.TaskSpec{
					Type:   "connect",
					Params: &tasks.ConnectParams{Host: host.Host, SSHConfig: host.SSHConfig},
				})
			}

			if process.StartType == "shell" {
				scriptPath := path.Join(utils.GetUtilPath(spec.InstallPath, process.Name), process.ExecShellPath)
				statusTasks = append(statusTasks, &types.TaskSpec{
					Type: "status",
					Params: &tasks.StatusParams{
						Component:     types.NebulaServiceComponent(component),
						Name:          process.Name,
						Host:          process.Hosts[0].Host,
						ExecShellPath: scriptPath,
						Port:          utils.GetConfigPort(process.Config),
					},
				})
			} else if process.StartType == "systemd" {
				statusTasks = append(statusTasks, &types.TaskSpec{
					Type: "systemd",
					Params: &tasks.SystemdParams{
						Host:    process.Hosts[0].Host,
						Name:    process.Name,
						Operate: "status",
						Port:    utils.GetConfigPort(process.Config),
					},
				})
			}
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
				Type:     "parallel",
				SubTasks: statusTasks,
			},
		},
	}, nil
}
