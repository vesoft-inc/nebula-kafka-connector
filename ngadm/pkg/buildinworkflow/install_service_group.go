package buildinworkflow

import (
	"fmt"
	"path"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func InstallServiceGroup(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	installTask, err := BuildInstallServiceGroupTask(args, spec)
	if err != nil {
		return nil, err
	}
	if installTask != nil {
		workflow.Tasks = append(workflow.Tasks, installTask)
	}
	return workflow, nil
}

func BuildInstallServiceGroupTask(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
	if args == nil {
		args = map[string]any{}
	}
	// username, password string
	if spec.Spec.Metad == nil {
		return nil, nil
	}
	username, ok := args["username"].(string)
	if !ok {
		return nil, fmt.Errorf("username is required")
	}
	password, ok := args["metaPassword"].(string)
	if !ok {
		return nil, fmt.Errorf("metaPassword is required")
	}

	metaServiceGroup := spec.Spec.Metad
	metaHosts := spec.Spec.Metad.Hosts
	metaServerAddress := utils.GetMetaAddressListString(metaHosts, utils.GetConfigPort(metaServiceGroup.Config))

	uploadTasks := []*types.TaskSpec{}
	connectTasks := []*types.TaskSpec{}
	force, ok := args["force"].(bool)
	if !ok {
		force = false
	}
	metadHosts := GetMetadAllNeedHosts(spec)
	needHosts := GetServiceGroupNeedHosts(spec)
	for _, agent := range needHosts {
		installPath := utils.GetUserCluster(spec.InstallPath, agent.InstallPath)
		isSameMetaHost := true
		if _, ok := metadHosts[agent.Host]; !ok {
			isSameMetaHost = false
		}
		// if the host is the same as the meta host, we default think use the meta package, we only need to connect
		if isSameMetaHost {
			connectTasks = append(connectTasks, &types.TaskSpec{
				Type: "serial",
				SubTasks: []*types.TaskSpec{
					{
						Type: "connect",
						Params: &tasks.ConnectParams{
							Host: agent.Host,
							// SSHConfig: agent.SSHConfig,
						},
					},
				},
			})
			continue
		}
		//1. connect
		connectTasks = append(connectTasks, &types.TaskSpec{
			Type: "serial",
			SubTasks: []*types.TaskSpec{
				{
					Type: "connect",
					Params: &tasks.ConnectParams{
						Host: agent.Host,
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
		//2.upload, if the package is already uploaded, it will be skipped
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
			},
		})
	}
	//3. config and start needed processes
	startNeededProcessesTask := []*types.TaskSpec{}
	for _, cluster := range spec.Spec.Metad.ServiceGroups {
		// 3.1 start graphd
		for _, host := range cluster.Graphd.Hosts {
			installPath := utils.GetUserCluster(spec.InstallPath, host.Agent.InstallPath)
			startNeededProcessesTask = append(startNeededProcessesTask, &types.TaskSpec{
				Type:        "serial",
				Description: "start graphd",
				SubTasks: []*types.TaskSpec{
					{
						Type: "init_config",
						Params: &tasks.InitConfigParams{
							Host: host.Agent.Host,
							ChangeMap: utils.MergeNebulaConfigMap(cluster.Graphd.Config, map[string]string{
								"local_ip":          utils.GetHostIP(host.Agent.Host),
								"meta_server_addrs": metaServerAddress,
							}),
							Dst: path.Join(installPath, "etc/nebula-graphd.conf"),
						},
					},
					{
						Type: "nebula_operation",
						Params: &tasks.NebulaOperationParams{
							Host:      host.Agent.Host,
							Operation: "start",
							Component: "graphd",
							Path:      installPath,
						},
					},
				},
			})
		}
		// 3.2 start storaged
		for _, host := range cluster.Storaged.Hosts {
			installPath := utils.GetUserCluster(spec.InstallPath, host.Agent.InstallPath)
			startNeededProcessesTask = append(startNeededProcessesTask, &types.TaskSpec{
				Type:        "serial",
				Description: "start storaged",
				SubTasks: []*types.TaskSpec{
					{
						Type: "init_config",
						Params: &tasks.InitConfigParams{
							Host: host.Agent.Host,
							ChangeMap: utils.MergeNebulaConfigMap(cluster.Storaged.Config, map[string]string{
								"local_ip":          utils.GetHostIP(host.Agent.Host),
								"meta_server_addrs": metaServerAddress,
							}),
							Dst: path.Join(installPath, "etc/nebula-storaged.conf"),
						},
					},
					{
						Type: "nebula_operation",
						Params: &tasks.NebulaOperationParams{
							Host:         host.Agent.Host,
							Operation:    "start",
							Component:    "storaged",
							NeedRollback: true,
							Path:         installPath,
						},
					},
				},
			})
		}
	}

	// 4.init cluster
	createServiceGroupTasks := []*types.TaskSpec{}
	for _, cluster := range spec.Spec.Metad.ServiceGroups {
		r := cluster
		createServiceGroupTasks = append(createServiceGroupTasks, &types.TaskSpec{
			Type: "create_cluster",
			Params: &tasks.CreateServiceGroupParams{
				ServiceGroupSpec:  &r,
				MetaSpec:          metaServiceGroup,
				MetaServerAddress: metaServerAddress,
				Username:          username,
				Password:          password,
			},
		})
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
				Description: "upload packages",
				SubTasks:    uploadTasks,
			},
			{
				Type:        "parallel",
				Description: "start needed processes",
				SubTasks:    startNeededProcessesTask,
			},
			{
				Type: "delay",
				Params: &tasks.DelayParams{
					Duration: 5 * time.Second,
				},
				Description: "wait for nebula service start",
			},
			{
				Type:     "serial",
				SubTasks: createServiceGroupTasks,
			},
		},
	}, nil
}

func GetServiceGroupNeedHosts(spec *types.JobSpec) map[string]*types.Agent {
	allNeedHosts := make(map[string]*types.Agent, 0)
	// for _, host := range spec.Spec.Metad.Hosts {
	// 	agentCopy := host.Agent
	// 	allNeedHosts[host.Agent.Host] = &agentCopy
	// }
	for _, cluster := range spec.Spec.Metad.ServiceGroups {
		for _, host := range cluster.Graphd.Hosts {
			agentCopy := host.Agent
			allNeedHosts[host.Agent.Host] = &agentCopy
		}
		for _, host := range cluster.Storaged.Hosts {
			agentCopy := host.Agent
			allNeedHosts[host.Agent.Host] = &agentCopy
		}
	}
	return allNeedHosts
}
