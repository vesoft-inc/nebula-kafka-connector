package buildinworkflow

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/utils"
)

func Operation(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "parallel",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	component, ok := args["component"].(string) //default "" means all
	if !ok || component == "" {
		component = "all"
	}
	host, ok := args["host"].(string)
	if !ok {
		host = ""
	}
	operation, ok := args["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation is not set")
	}
	killWait, ok := args["kill-wait"].(string)
	if !ok {
		killWait = ""
	}

	// append nebula component operation
	if isNebulaComponent(component) {
		operationTask, err := OperationCluster(spec, operation, component, host, killWait)
		if err != nil {
			return nil, err
		}
		if operationTask != nil {
			workflow.Tasks = append(workflow.Tasks, operationTask)
		}
	}
	//todo: add other component operation. explorer, studio, etc
	return workflow, nil
}

func isNebulaComponent(component string) bool {
	if component == "all" {
		return true
	}
	_, ok := types.NebulaServiceComponentMap[component]
	return ok
}

func addNeedOperation(allNeedOperations *map[string]map[types.NebulaServiceComponent]bool, nowComponent types.NebulaServiceComponent, newComponent types.NebulaServiceComponent, nowHost, newHost string) bool {
	flag := false
	if nowComponent == newComponent || newComponent == types.AllNebulaSerivce {
		flag = true
	} else {
		flag = false
	}
	if newHost == "" || utils.GetHostIP(nowHost) == utils.GetHostIP(newHost) {
		flag = true
	} else {
		flag = false
	}
	if flag {
		if (*allNeedOperations)[nowHost] == nil {
			(*allNeedOperations)[nowHost] = make(map[types.NebulaServiceComponent]bool)
		}
		(*allNeedOperations)[nowHost][newComponent] = true
	}
	return flag
}

func OperationCluster(spec *types.JobSpec, operation, component, host, KillWait string) (*types.TaskSpec, error) {
	componentType := types.NebulaServiceComponentMap[component]
	if spec.Spec.Metad == nil {
		return nil, fmt.Errorf("metad spec is nil")
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
	operationTasks := []*types.TaskSpec{}
	//1. connect
	//2. operation
	for host, components := range allNeedOperations {
		connectTasks = append(connectTasks, &types.TaskSpec{
			Type:   "connect",
			Params: &tasks.ConnectParams{Host: host},
		})
		// aggregate operation task
		if (components[types.Metad] && components[types.Graphd] && components[types.Storaged]) || components[types.AllNebulaSerivce] {
			operationTasks = append(operationTasks, &types.TaskSpec{
				Type: "nebula_operation",
				Params: &tasks.NebulaOperationParams{Operation: operation, Component: types.AllNebulaSerivce, Host: host,
					KillWait: KillWait,
					Path:     utils.GetClusterPath(spec.InstallPath)},
			})
			continue
		}
		for component := range components {
			operationTasks = append(operationTasks, &types.TaskSpec{
				Type: "nebula_operation",
				Params: &tasks.NebulaOperationParams{Operation: operation, Component: component, Host: host,
					KillWait: KillWait,
					Path:     utils.GetClusterPath(spec.InstallPath)},
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
				SubTasks: operationTasks,
			},
		},
	}, nil
}
