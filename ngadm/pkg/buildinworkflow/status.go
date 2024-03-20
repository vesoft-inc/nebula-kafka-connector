package buildinworkflow

import (
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func Status(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "parallel",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	component, ok := args["component"].(string) //default "" means all
	if !ok || component == "" {
		component = "all"
	}
	// append nebula component status
	if isNebulaComponent(component) {
		statusTask, err := StatusCluster(spec, component, "")
		if err != nil {
			return nil, err
		}
		if statusTask != nil {
			workflow.Tasks = append(workflow.Tasks, statusTask)
		}
	}
	utilStatusTask, err := StatusUtils(spec, component)
	if err != nil {
		return nil, err
	}
	if utilStatusTask != nil {
		workflow.Tasks = append(workflow.Tasks, utilStatusTask)
	}

	return workflow, nil
}

func StatusCluster(spec *types.JobSpec, component, host string) (*types.TaskSpec, error) {
	componentType := types.NebulaServiceComponentMap[component]
	if spec.Spec.Metad == nil {
		return nil, nil
	}
	metaHosts := spec.Spec.Metad.Hosts
	allNeedOperations := make(map[string]map[types.NebulaServiceComponent]bool, 0)
	for _, agent := range metaHosts {
		addNeedOperation(&allNeedOperations, types.Metad, componentType, agent.Host, host)
	}
	for _, cluster := range spec.Spec.Metad.Clusters {
		for _, agent := range cluster.Graphd.Hosts {
			addNeedOperation(&allNeedOperations, types.Graphd, componentType, agent.Host, host)
		}
		for _, agent := range cluster.Storaged.Hosts {
			addNeedOperation(&allNeedOperations, types.Storaged, componentType, agent.Host, host)
		}
	}

	connectTasks := []*types.TaskSpec{}
	statusTasks := []*types.TaskSpec{}
	//1. connect
	//2. status
	for host, components := range allNeedOperations {
		connectTasks = append(connectTasks, &types.TaskSpec{
			Type:   "connect",
			Params: &tasks.ConnectParams{Host: host},
		})
		// aggregate status task
		if (components[types.Metad] && components[types.Graphd] && components[types.Storaged]) || components[types.AllNebulaSerivce] {
			statusTasks = append(statusTasks, &types.TaskSpec{
				Type: "nebula_status",
				Params: &tasks.NebulaStatusParams{
					Component: types.AllNebulaSerivce,
					Host:      host,
					Path:      utils.GetClusterPath(spec.InstallPath)},
			})
			continue
		}
		for component := range components {
			statusTasks = append(statusTasks, &types.TaskSpec{
				Type: "nebula_status",
				Params: &tasks.NebulaStatusParams{Component: component, Host: host,
					Path: utils.GetClusterPath(spec.InstallPath)},
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
				Type:     "parallel",
				SubTasks: statusTasks,
			},
		},
	}, nil
}
