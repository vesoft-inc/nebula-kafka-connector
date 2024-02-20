package buildinworkflow

import (
	"path"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/utils"
)

func Install(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	installTask, err := InstallCluster(args, spec)
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

func InstallCluster(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
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
						Path:  utils.GetClusterPath(spec.InstallPath),
						Force: force,
					},
				},
			},
		})
		//2.upload
		dstPath := path.Join(utils.GetDownloadPath(spec.InstallPath), path.Base(metaCluster.PackagePath))
		uploadTasks = append(uploadTasks, &types.TaskSpec{
			Type: "serial",
			SubTasks: []*types.TaskSpec{
				{
					Type: "upload",
					Params: &tasks.UploadParams{
						SrcPath: metaCluster.PackagePath,
						DstPath: dstPath,
						Host:    agent.Host,
					},
				},
				{
					Type: "extract",
					Params: &tasks.ExtractParams{
						Host:        agent.Host,
						ExtractPath: utils.GetClusterPath(spec.InstallPath),
						PkgPath:     dstPath,
					},
				},
			},
		})
	}
	//3. config and start needed processes
	startNeededProcessesTask := []*types.TaskSpec{}
	for _, agent := range metaHosts {
		startNeededProcessesTask = append(startNeededProcessesTask, &types.TaskSpec{
			Type:        "serial",
			Description: "start metad",
			SubTasks: []*types.TaskSpec{
				{
					Type: "init_config",
					Params: &tasks.InitConfigParams{
						Host: agent.Host,
						ChangeMap: utils.MergeConfigMap(metaCluster.Config, map[string]string{
							"local_ip":          utils.GetHostIP(agent.Host),
							"meta_server_addrs": utils.GetMetaAddressListString(metaHosts, metaCluster.Config["port"]),
						}),
						Dst: path.Join(utils.GetClusterPath(spec.InstallPath), "etc/nebula-metad.conf"),
					},
				},
				{
					Type: "nebula_operation",
					Params: &tasks.NebulaOperationParams{
						Host:         agent.Host,
						Operation:    "start",
						Component:    "metad",
						NeedRollback: true,
						Path:         utils.GetClusterPath(spec.InstallPath),
					},
				},
			},
		})
	}
	for _, cluster := range spec.Spec.Metad.Clusters {
		// 4.1 start graphd
		for _, agent := range cluster.Graphd.Hosts {
			startNeededProcessesTask = append(startNeededProcessesTask, &types.TaskSpec{
				Type:        "serial",
				Description: "start graphd",
				SubTasks: []*types.TaskSpec{
					{
						Type: "init_config",
						Params: &tasks.InitConfigParams{
							Host: agent.Host,
							ChangeMap: utils.MergeConfigMap(cluster.Graphd.Config, map[string]string{
								"local_ip":          utils.GetHostIP(agent.Host),
								"meta_server_addrs": utils.GetMetaAddressListString(metaHosts, metaCluster.Config["port"]),
							}),
							Dst: path.Join(utils.GetClusterPath(spec.InstallPath), "etc/nebula-graphd.conf"),
						},
					},
					{
						Type: "nebula_operation",
						Params: &tasks.NebulaOperationParams{
							Host:      agent.Host,
							Operation: "start",
							Component: "graphd",
							Path:      utils.GetClusterPath(spec.InstallPath),
						},
					},
				},
			})
		}
		// 4.2 start storaged
		for _, agent := range cluster.Storaged.Hosts {
			startNeededProcessesTask = append(startNeededProcessesTask, &types.TaskSpec{
				Type:        "serial",
				Description: "start storaged",
				SubTasks: []*types.TaskSpec{
					{
						Type: "init_config",
						Params: &tasks.InitConfigParams{
							Host: agent.Host,
							ChangeMap: utils.MergeConfigMap(cluster.Storaged.Config, map[string]string{
								"local_ip":          utils.GetHostIP(agent.Host),
								"meta_server_addrs": utils.GetMetaAddressListString(metaHosts, metaCluster.Config["port"]),
							}),
							Dst: path.Join(utils.GetClusterPath(spec.InstallPath), "etc/nebula-storaged.conf"),
						},
					},
					{
						Type: "nebula_operation",
						Params: &tasks.NebulaOperationParams{
							Host:         agent.Host,
							Operation:    "start",
							Component:    "storaged",
							NeedRollback: true,
							Path:         utils.GetClusterPath(spec.InstallPath),
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
				ClusterSpec: &cluster,
				MetaSpec:    metaCluster,
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

func GetMetadAllNeedHosts(spec *types.JobSpec) map[string]*types.Agent {
	allNeedHosts := make(map[string]*types.Agent, 0)
	for _, agent := range spec.Spec.Metad.Hosts {
		allNeedHosts[agent.Host] = &agent
	}
	for _, cluster := range spec.Spec.Metad.Clusters {
		for _, agent := range cluster.Graphd.Hosts {
			allNeedHosts[agent.Host] = &agent
		}
		for _, agent := range cluster.Storaged.Hosts {
			allNeedHosts[agent.Host] = &agent
		}
	}
	return allNeedHosts
}
