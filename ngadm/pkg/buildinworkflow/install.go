package buildinworkflow

import (
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func Install(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	installTask, err := InstallMetaCluster(args, spec)
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

func InstallMetaCluster(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
	if args == nil {
		args = map[string]any{}
	}
	if spec.Spec.Metad == nil {
		return nil, nil
	}
	metaCluster := spec.Spec.Metad
	metaHosts := spec.Spec.Metad.Hosts

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
	//3. config and start needed processes
	startNeededProcessesTask := []*types.TaskSpec{}
	for _, agent := range metaHosts {
		installPath := utils.GetUserClusterPath(spec.InstallPath, agent.PackagePath)
		startNeededProcessesTask = append(startNeededProcessesTask, &types.TaskSpec{
			Type:        "serial",
			Description: "start metad",
			SubTasks: []*types.TaskSpec{
				{
					Type: "init_config",
					Params: &tasks.InitConfigParams{
						Host: agent.Host,
						ChangeMap: utils.MergeNebulaConfigMap(metaCluster.Config, map[string]string{
							"local_ip":          utils.GetHostIP(agent.Host),
							"meta_server_addrs": utils.GetMetaAddressListString(metaHosts, utils.GetConfigPort(metaCluster.Config)),
						}),
						Dst: path.Join(installPath, "etc/nebula-metad.conf"),
					},
				},
				{
					Type: "nebula_operation",
					Params: &tasks.NebulaOperationParams{
						Host:         agent.Host,
						Operation:    "start",
						Component:    "metad",
						NeedRollback: true,
						Path:         installPath,
					},
				},
				{
					Type: "save-agent-config",
					Params: &tasks.SaveAgentParams{
						Component: "metad",
						Config: map[string]any{
							"installPath": installPath,
							"host":        utils.GetHostIP(agent.Host),
							"port":        utils.GetConfigPort(metaCluster.Config),
						},
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
			{
				Type:        "parallel",
				Description: "start needed processes",
				SubTasks:    startNeededProcessesTask,
			},
		},
	}
	return mainTask, nil
}

func GetMetadAllNeedHosts(spec *types.JobSpec) map[string]*types.Agent {
	allNeedHosts := make(map[string]*types.Agent, 0)
	for _, agent := range spec.Spec.Metad.Hosts {
		agentCopy := agent
		allNeedHosts[agent.Host] = &agentCopy
	}
	for _, cluster := range spec.Spec.Metad.Clusters {
		for _, agent := range cluster.Graphd.Hosts {
			agentCopy := agent
			allNeedHosts[agent.Host] = &agentCopy
		}
		for _, agent := range cluster.Storaged.Hosts {
			agentCopy := agent
			allNeedHosts[agent.Host] = &agentCopy
		}
	}
	return allNeedHosts
}
