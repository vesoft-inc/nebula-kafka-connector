package buildinworkflow

import (
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func Config(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	configTask, err := ConfigServiceGroup(args, spec)
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

func ConfigServiceGroup(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
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
						Host:      agent.Host,
						SSHConfig: agent.SSHConfig,
					},
				},
			},
		})
	}
	configTask := []*types.TaskSpec{}
	for _, host := range spec.Spec.Metad.Hosts {
		configTask = append(configTask, &types.TaskSpec{
			Type:        "serial",
			Description: "config metad",
			SubTasks: []*types.TaskSpec{
				{
					Type: "config",
					Params: &tasks.ConfigParams{
						Host:   host.Agent.Host,
						Config: spec.Spec.Metad.Config,
						Type:   "nebulagraph",
						Dst:    path.Join(utils.GetServiceGroupPath(spec.InstallPath), "etc/nebula-metad.conf"),
					},
				},
			},
		})
	}
	for _, cluster := range spec.Spec.Metad.ServiceGroups {
		for _, host := range cluster.Graphd.Hosts {
			configTask = append(configTask, &types.TaskSpec{
				Type:        "serial",
				Description: "config graphd",
				SubTasks: []*types.TaskSpec{
					{
						Type: "config",
						Params: &tasks.ConfigParams{
							Host:   host.Agent.Host,
							Config: cluster.Graphd.Config,
							Type:   "nebulagraph",
							Dst:    path.Join(utils.GetServiceGroupPath(spec.InstallPath), "etc/nebula-graphd.conf"),
						},
					},
				},
			})
		}
		// 4.2 start storaged
		for _, host := range cluster.Storaged.Hosts {
			configTask = append(configTask, &types.TaskSpec{
				Type:        "serial",
				Description: "config storaged",
				SubTasks: []*types.TaskSpec{
					{
						Type: "config",
						Params: &tasks.ConfigParams{
							Host:   host.Agent.Host,
							Config: cluster.Storaged.Config,
							Dst:    path.Join(utils.GetServiceGroupPath(spec.InstallPath), "etc/nebula-storaged.conf"),
							Type:   "nebulagraph",
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
	allUtils := utils.GetAllUtilsProcess(spec)
	//1. connect all need hosts
	connectTasks := []*types.TaskSpec{}
	for _, process := range allUtils {
		for _, host := range process.Hosts {
			connectTasks = append(connectTasks, &types.TaskSpec{
				Type: "connect",
				Params: &tasks.ConnectParams{
					Host:      host.Agent.Host,
					SSHConfig: host.Agent.SSHConfig,
				},
			})
		}
	}
	//2. upload and start
	//2.1 todo: apply utils config
	configTasks := []*types.TaskSpec{}
	for _, process := range allUtils {
		for _, host := range process.Hosts {
			task := &types.TaskSpec{
				Type: "serial",
				SubTasks: []*types.TaskSpec{
					{
						Type: "config",
						Params: &tasks.ConfigParams{
							Host:   host.Agent.Host,
							Config: process.Config,
							Dst:    path.Join(utils.GetUtilPath(spec.InstallPath, process.Name), process.ConfigPath),
							Type:   "yaml",
						},
					},
				},
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
