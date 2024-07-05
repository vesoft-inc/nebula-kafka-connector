package buildinworkflow

import (
	"fmt"
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func InstallAgent(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	process := spec.UtilsProcesses["agent"]
	//1. connect all need hosts
	connectTasks := []*types.TaskSpec{}
	for _, host := range process.Hosts {
		if host.Agent.SSHConfig == nil {
			continue
		}
		connectTasks = append(connectTasks, &types.TaskSpec{
			Type: "connect",
			Params: &tasks.ConnectParams{
				Host:      host.Agent.Host,
				SSHConfig: host.Agent.SSHConfig,
			},
		})
	}
	//2. upload and start
	//2.1 todo: apply utils config
	uploadAndStartTasks := []*types.TaskSpec{}
	for _, host := range process.Hosts {
		dstPath := path.Join(utils.GetDownloadPath(spec.InstallPath), path.Base(process.PackagePath))
		installPath := utils.GetUserUtilPath(spec.InstallPath, host.Agent.InstallPath, process.Name)
		task := &types.TaskSpec{
			Type: "serial",
			SubTasks: []*types.TaskSpec{
				{
					Type: "upload",
					Params: &tasks.UploadParams{
						SrcPath: process.PackagePath,
						DstPath: dstPath,
						Host:    host.Agent.Host,
					},
				},
				{
					Type: "extract",
					Params: &tasks.ExtractParams{
						Host:        host.Agent.Host,
						ExtractPath: installPath,
						PkgPath:     dstPath,
					},
				},
				{
					Type: "config",
					Params: &tasks.ConfigParams{
						Host:   host.Agent.Host,
						Config: utils.MergeAgentConfig(host.Agent.Config, process.Config),
						Dst:    path.Join(installPath, process.ConfigPath),
						Type:   "yaml",
					},
				},
			},
		}
		// upload server certs file
		if spec.ServerCertFile == "" || spec.ServerKeyFile == "" {
			return nil, fmt.Errorf("agent server cert file or server key file is empty")
		}
		task.SubTasks = append(task.SubTasks, &types.TaskSpec{
			Type: "upload",
			Params: &tasks.UploadParams{
				SrcPath: spec.ServerCertFile,
				DstPath: path.Join(installPath, "certs/server.crt"),
				Host:    host.Agent.Host,
			},
		}, &types.TaskSpec{
			Type: "upload",
			Params: &tasks.UploadParams{
				SrcPath: spec.ServerKeyFile,
				DstPath: path.Join(installPath, "certs/server.key"),
				Host:    host.Agent.Host,
			},
		}, &types.TaskSpec{
			Type: "upload",
			Params: &tasks.UploadParams{
				SrcPath: spec.CAFile,
				DstPath: path.Join(installPath, "certs/ca.crt"),
				Host:    host.Agent.Host,
			},
		})

		if process.StartType == "shell" {
			scriptPath := path.Join(installPath, process.ExecShellPath)
			task.SubTasks = append(task.SubTasks, &types.TaskSpec{
				Type: "operate",
				Params: &tasks.OperateParams{
					Host:      host.Agent.Host,
					ExecPath:  scriptPath,
					Operation: "start",
				},
			})
		} else if process.StartType == "systemd" {
			task.SubTasks = append(task.SubTasks, &types.TaskSpec{
				Type: "systemd",
				Params: &tasks.SystemdParams{
					Host:             host.Agent.Host,
					Name:             process.Name,
					ExecStartPath:    path.Join(installPath, process.ExecStartPath),
					WorkingDirectory: path.Join(installPath, process.WorkingDir),
					Operate:          "install",
				},
			})
		}
		uploadAndStartTasks = append(uploadAndStartTasks, task)
	}
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks: []*types.TaskSpec{
			{
				Type:     "parallel",
				SubTasks: connectTasks,
			},
			{
				Type:        "parallel",
				Description: "upload and start agent",
				SubTasks:    uploadAndStartTasks,
			},
		},
	}
	return workflow, nil
}
