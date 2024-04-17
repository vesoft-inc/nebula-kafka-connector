package buildinworkflow

import (
	"path"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func InstallCluster(args map[string]any, spec *types.JobSpec, username, password string) (*types.TaskSpec, error) {
	if args == nil {
		args = map[string]any{}
	}
	if spec.Spec.Metad == nil {
		return nil, nil
	}
	metaCluster := spec.Spec.Metad
	metaHosts := spec.Spec.Metad.Hosts
	metaServerAddress := utils.GetMetaAddressListString(metaHosts, utils.GetConfigPort(metaCluster.Config))

	uploadTasks := []*types.TaskSpec{}
	connectTasks := []*types.TaskSpec{}
	force, ok := args["force"].(bool)
	if !ok {
		force = false
	}
	needHosts := GetClusterNeedHosts(spec)
	for _, agent := range needHosts {
		if agent.SSHConfig == nil {
			connectTasks = append(connectTasks, &types.TaskSpec{
				Type: "serial",
				SubTasks: []*types.TaskSpec{
					{
						Type: "connect",
						Params: &tasks.ConnectParams{
							Host: agent.Host,
						},
					},
				},
			})
			continue
		}

		// TODO: check if the host already has nebula installed
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
	for _, cluster := range spec.Spec.Metad.Clusters {
		// 3.1 start graphd
		for _, agent := range cluster.Graphd.Hosts {
			startNeededProcessesTask = append(startNeededProcessesTask, &types.TaskSpec{
				Type:        "serial",
				Description: "start graphd",
				SubTasks: []*types.TaskSpec{
					{
						Type: "init_config",
						Params: &tasks.InitConfigParams{
							Host: agent.Host,
							ChangeMap: utils.MergeNebulaConfigMap(cluster.Graphd.Config, map[string]string{
								"local_ip":          utils.GetHostIP(agent.Host),
								"meta_server_addrs": metaServerAddress,
							}),
							Dst: path.Join(utils.GetUserClusterPath(spec.InstallPath, agent.InstallPath), "etc/nebula-graphd.conf"),
						},
					},
					{
						Type: "nebula_operation",
						Params: &tasks.NebulaOperationParams{
							Host:      agent.Host,
							Operation: "start",
							Component: "graphd",
							Path:      utils.GetUserClusterPath(spec.InstallPath, agent.InstallPath),
						},
					},
				},
			})
		}
		// .2 start storaged
		for _, agent := range cluster.Storaged.Hosts {
			startNeededProcessesTask = append(startNeededProcessesTask, &types.TaskSpec{
				Type:        "serial",
				Description: "start storaged",
				SubTasks: []*types.TaskSpec{
					{
						Type: "init_config",
						Params: &tasks.InitConfigParams{
							Host: agent.Host,
							ChangeMap: utils.MergeNebulaConfigMap(cluster.Storaged.Config, map[string]string{
								"local_ip":          utils.GetHostIP(agent.Host),
								"meta_server_addrs": metaServerAddress,
							}),
							Dst: path.Join(utils.GetUserClusterPath(spec.InstallPath, agent.InstallPath), "etc/nebula-storaged.conf"),
						},
					},
					{
						Type: "nebula_operation",
						Params: &tasks.NebulaOperationParams{
							Host:         agent.Host,
							Operation:    "start",
							Component:    "storaged",
							NeedRollback: true,
							Path:         utils.GetUserClusterPath(spec.InstallPath, agent.InstallPath),
						},
					},
				},
			})
		}
	}

	// 4.init cluster
	createClusterTasks := []*types.TaskSpec{}
	for _, cluster := range spec.Spec.Metad.Clusters {
		createClusterTasks = append(createClusterTasks, &types.TaskSpec{
			Type: "create_cluster",
			Params: &tasks.CreateClusterParams{
				ClusterSpec:       &cluster,
				MetaSpec:          metaCluster,
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
				SubTasks: createClusterTasks,
			},
		},
	}, nil
}

func GetClusterNeedHosts(spec *types.JobSpec) map[string]*types.Agent {
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
