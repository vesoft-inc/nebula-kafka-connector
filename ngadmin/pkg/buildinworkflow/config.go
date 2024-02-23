package buildinworkflow

import (
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/utils"
)

func Config(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	configTask, err := ConfigCluster(args, spec)
	if err != nil {
		return nil, err
	}
	if configTask != nil {
		workflow.Tasks = append(workflow.Tasks, configTask)
	}
	configTask, err = ConfigUtils(args, spec)
	if err != nil {
		return nil, err
	}
	if configTask != nil {
		workflow.Tasks = append(workflow.Tasks, configTask)
	}

	return workflow, nil
}

func ConfigCluster(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
	if spec.Spec.Metad == nil {
		return nil, nil
	}
	connectTasks := []*types.TaskSpec{}
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
			},
		})
	}
	configTask := []*types.TaskSpec{}
	for _, agent := range spec.Spec.Metad.Hosts {
		configTask = append(configTask, &types.TaskSpec{
			Type:        "serial",
			Description: "config graphd",
			SubTasks: []*types.TaskSpec{
				{
					Type: "config",
					Params: &tasks.ConfigParams{
						Host:   agent.Host,
						Config: spec.Spec.Metad.Config,
						Type:   "nebulagraph",
						Dst:    path.Join(utils.GetClusterPath(spec.InstallPath), "etc/nebula-metad.conf"),
					},
				},
				{
					Type: "nebula_operation",
					Params: &tasks.NebulaOperationParams{
						Host:      agent.Host,
						Operation: "restart",
						Component: "metad",
						Path:      utils.GetClusterPath(spec.InstallPath),
					},
				},
			},
		})
	}
	for _, cluster := range spec.Spec.Metad.Clusters {
		for _, agent := range cluster.Graphd.Hosts {
			configTask = append(configTask, &types.TaskSpec{
				Type:        "serial",
				Description: "config graphd",
				SubTasks: []*types.TaskSpec{
					{
						Type: "config",
						Params: &tasks.ConfigParams{
							Host:   agent.Host,
							Config: cluster.Graphd.Config,
							Type:   "nebulagraph",
							Dst:    path.Join(utils.GetClusterPath(spec.InstallPath), "etc/nebula-graphd.conf"),
						},
					},
					{
						Type: "nebula_operation",
						Params: &tasks.NebulaOperationParams{
							Host:      agent.Host,
							Operation: "restart",
							Component: "graphd",
							Path:      utils.GetClusterPath(spec.InstallPath),
						},
					},
				},
			})
		}
		// 4.2 start storaged
		for _, agent := range cluster.Storaged.Hosts {
			configTask = append(configTask, &types.TaskSpec{
				Type:        "serial",
				Description: "config storaged",
				SubTasks: []*types.TaskSpec{
					{
						Type: "config",
						Params: &tasks.ConfigParams{
							Host:   agent.Host,
							Config: cluster.Storaged.Config,
							Dst:    path.Join(utils.GetClusterPath(spec.InstallPath), "etc/nebula-storaged.conf"),
							Type:   "nebulagraph",
						},
					},
					{
						Type: "nebula_operation",
						Params: &tasks.NebulaOperationParams{
							Host:         agent.Host,
							Operation:    "restart",
							Component:    "storaged",
							NeedRollback: true,
							Path:         utils.GetClusterPath(spec.InstallPath),
						},
					},
				},
			})
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
				Type:     "serial",
				SubTasks: configTask,
			},
		},
	}, nil
}

func ConfigUtils(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
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
	//2.1 todo: apply utils config
	configTasks := []*types.TaskSpec{}
	for _, process := range allUtils {
		for _, agent := range process.Hosts {
			task := &types.TaskSpec{
				Type: "serial",
				SubTasks: []*types.TaskSpec{
					{
						Type:        "config",
						Description: "config " + process.Name,
						Params: &tasks.ConfigParams{
							Host:   agent.Host,
							Config: process.Config,
							Dst:    path.Join(utils.GetUtilPath(spec.InstallPath, process.Name), process.ConfigPath),
							Type:   "yaml",
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
						Operation: "restart",
					},
				})
			} else if process.StartType == "systemd" {
				task.SubTasks = append(task.SubTasks, &types.TaskSpec{
					Type: "systemd",
					Params: &tasks.SystemdParams{
						Host:    agent.Host,
						Name:    process.Name,
						Operate: "restart",
					},
				})
			}
			configTasks = append(configTasks, task)
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
				Type:     "serial",
				SubTasks: configTasks,
			},
		},
	}, nil
}
