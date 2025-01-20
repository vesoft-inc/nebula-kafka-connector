package buildinworkflow

import (
	"fmt"
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

const defaultMetadResetTimeout = 20
const defaultMetadUser = "root"

func Install(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	installTask, err := InstallMetad(args, spec)
	if err != nil {
		return nil, err
	}
	if installTask != nil {
		workflow.Tasks = append(workflow.Tasks, installTask)
	}

	installTask, err = BuildInstallServiceGroupTask(args, spec)
	if err != nil {
		return nil, err
	}
	if installTask != nil {
		workflow.Tasks = append(workflow.Tasks, installTask)
	}

	installUtilsTask, err := InstallUtils(args, spec)
	if err != nil {
		return nil, err
	}
	if installUtilsTask != nil {
		workflow.Tasks = append(workflow.Tasks, installUtilsTask)
	}

	return workflow, nil
}

func InstallMetad(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
	if args == nil {
		args = map[string]any{}
	}
	if spec.Spec.Metad == nil {
		return nil, nil
	}
	metaServiceGroup := spec.Spec.Metad
	metaHosts := spec.Spec.Metad.Hosts

	uploadTasks := []*types.TaskSpec{}
	connectTasks := []*types.TaskSpec{}
	force, ok := args["force"].(bool)
	if !ok {
		force = false
	}
	allNeedHosts := GetMetadAllNeedHosts(spec)
	for _, agent := range allNeedHosts {
		installPath := utils.GetUserCluster(spec.InstallPath, agent.InstallPath)
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
		dstPath := path.Join(utils.GetUserDownloadPath(spec.InstallPath, agent.InstallPath), path.Base(metaServiceGroup.PackagePath))
		uploadTasks = append(uploadTasks, &types.TaskSpec{
			Type: "serial",
			SubTasks: []*types.TaskSpec{
				{
					Type: "upload",
					Params: &tasks.UploadParams{
						SrcPath: utils.GetUserPackagePath(metaServiceGroup.PackagePath, agent.PackagePath),
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
				{
					Type: "shell",
					Params: &tasks.ShellParams{
						Host: agent.Host,
						Command: fmt.Sprintf("mkdir -p %s && cd %s && echo %s > %s",
							INSTALL_PATH_DIR, INSTALL_PATH_DIR, installPath, INSTALL_PATH_FILE),
					},
				},
			},
		})
	}
	//3. config and start needed processes
	startNeededProcessesTask := []*types.TaskSpec{}
	for _, host := range metaHosts {
		installPath := utils.GetUserCluster(spec.InstallPath, host.Agent.PackagePath)
		startNeededProcessesTask = append(startNeededProcessesTask, &types.TaskSpec{
			Type:        "serial",
			Description: "start metad",
			SubTasks: []*types.TaskSpec{
				{
					Type: "init_config",
					Params: &tasks.InitConfigParams{
						Host: host.Agent.Host,
						ChangeMap: utils.MergeNebulaConfigMap(metaServiceGroup.Config, map[string]string{
							"local_ip":          utils.GetHostIP(host.Agent.Host),
							"meta_server_addrs": utils.GetMetaAddressListString(metaHosts, utils.GetConfigPort(metaServiceGroup.Config)),
						}),
						Dst: path.Join(installPath, "etc/nebula-metad.conf"),
					},
				},
				{
					Type: "nebula_operation",
					Params: &tasks.NebulaOperationParams{
						Host:         host.Agent.Host,
						Operation:    "start",
						Component:    "metad",
						NeedRollback: true,
						Path:         installPath,
					},
				},
				//{
				//	Type: "save-agent-config",
				//	Params: &tasks.SaveAgentParams{
				//		Component: "metad",
				//		Config: map[string]any{
				//			"installPath": installPath,
				//			"host":        utils.GetHostIP(agent.Host),
				//			"port":        utils.GetConfigPort(metaServiceGroup.Config),
				//		},
				//	},
				//},
			},
		})
	}
	//4. reset metad password
	password, ok := args["password"].(string)
	if !ok {
		return nil, fmt.Errorf("password is required")
	}
	metaPassword, ok := args["metaPassword"].(string)
	if !ok {
		return nil, fmt.Errorf("metaPassword is required")
	}
	timeout, ok := args["timeout"].(int)
	if !ok {
		timeout = defaultMetadResetTimeout
	}
	resetPasswordTask := &types.TaskSpec{
		Type: "reset_meta_password",
		Params: &tasks.ResetMetaPasswordParams{
			MetaServerAddress: utils.GetMetaAddressListString(metaHosts, utils.GetConfigPort(metaServiceGroup.Config)),
			Username:          defaultMetadUser,
			Password:          password,
			NewPassword:       metaPassword,
			TimeoutSec:        timeout,
		},
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
			{
				Type:        "parallel",
				Description: "start needed processes",
				SubTasks:    startNeededProcessesTask,
			},
			resetPasswordTask,
		},
	}
	return mainTask, nil
}

func GetMetadAllNeedHosts(spec *types.JobSpec) map[string]*types.Agent {
	allNeedHosts := make(map[string]*types.Agent, 0)
	for _, host := range spec.Spec.Metad.Hosts {
		agentCopy := host.Agent
		allNeedHosts[host.Agent.Host] = &agentCopy
	}
	return allNeedHosts
}

func GetMetadSelectedHosts(selectedHost []string, spec *types.JobSpec) map[string]*types.Agent {
	selectedHosts := make(map[string]*types.Agent, 0)
	for _, host := range spec.Spec.Metad.Hosts {
		for _, h := range selectedHost {
			if host.IP == h {
				agentCopy := host.Agent
				selectedHosts[host.Agent.Host] = &agentCopy
			}
		}
	}
	for _, cluster := range spec.Spec.Metad.ServiceGroups {
		for _, host := range cluster.Graphd.Hosts {
			for _, h := range selectedHost {
				if host.IP == h {
					agentCopy := host.Agent
					selectedHosts[host.Agent.Host] = &agentCopy
				}
			}
		}
		for _, host := range cluster.Storaged.Hosts {
			for _, h := range selectedHost {
				if host.IP == h {
					agentCopy := host.Agent
					selectedHosts[host.Agent.Host] = &agentCopy
				}
			}
		}
	}
	return selectedHosts
}
