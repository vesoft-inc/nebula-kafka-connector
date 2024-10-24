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
		statusTask, err := StatusServiceGroup(spec, component, "")
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

func StatusServiceGroup(spec *types.JobSpec, component, host string) (*types.TaskSpec, error) {
	componentType := types.NebulaServiceComponentMap[component]
	if spec.Spec.Metad == nil {
		return nil, nil
	}
	allAgents := utils.GetAllAgents(spec)
	metaHosts := spec.Spec.Metad.Hosts
	allNeedOperations := make(map[string]map[types.NebulaServiceComponent]bool, 0)
	for _, _host := range metaHosts {
		addNeedOperation(&allNeedOperations, types.Metad, componentType, _host.Agent.Host, host)
	}
	for _, cluster := range spec.Spec.Metad.ServiceGroups {
		for _, _host := range cluster.Graphd.Hosts {
			addNeedOperation(&allNeedOperations, types.Graphd, componentType, _host.Agent.Host, host)
		}
		for _, _host := range cluster.Storaged.Hosts {
			addNeedOperation(&allNeedOperations, types.Storaged, componentType, _host.Agent.Host, host)
		}
	}

	connectTasks := []*types.TaskSpec{}
	statusTasks := []*types.TaskSpec{}
	//1. connect
	agentsMap := make(map[string]*types.Agent)
	for _, agent := range allAgents {
		connectTasks = append(connectTasks, &types.TaskSpec{
			Type:   "connect",
			Params: &tasks.ConnectParams{Host: agent.Host, SSHConfig: agent.SSHConfig},
		})
		agentsMap[agent.Host] = agent
	}
	//2. status
	for host, components := range allNeedOperations {
		installPath := utils.GetUserCluster(spec.InstallPath, agentsMap[host].InstallPath)
		// aggregate status task
		if (components[types.Metad] && components[types.Graphd] && components[types.Storaged]) || components[types.AllNebulaSerivce] {
			statusTasks = append(statusTasks, &types.TaskSpec{
				Type: "nebula_status",
				Params: &tasks.NebulaStatusParams{
					Component: types.AllNebulaSerivce,
					Host:      host,
					Path:      installPath},
			})
			continue
		}
		for component := range components {
			statusTasks = append(statusTasks, &types.TaskSpec{
				Type: "nebula_status",
				Params: &tasks.NebulaStatusParams{Component: component, Host: host,
					Path: installPath},
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
