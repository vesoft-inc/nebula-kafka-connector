package buildinworkflow

import (
	"fmt"
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
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
		operationTask, err := OperationServiceGroup(spec, operation, component, host, killWait)
		if err != nil {
			return nil, err
		}
		if operationTask != nil {
			workflow.Tasks = append(workflow.Tasks, operationTask)
		}
	}

	// append utils operation
	operationTask, err := OperateUtils(spec, component, operation, host)
	if err != nil {
		return nil, err
	}
	if operationTask != nil {
		workflow.Tasks = append(workflow.Tasks, operationTask)
	}

	return workflow, nil
}

func OperateUtils(spec *types.JobSpec, component, operation, host string) (*types.TaskSpec, error) {
	allUitils := utils.GetAllUtilsProcess(spec)
	operateTasks := []*types.TaskSpec{}
	connectTasks := []*types.TaskSpec{}
	for _, process := range allUitils {
		for _, _h := range process.Hosts {
			if host != "" && host != _h.Agent.Host {
				continue
			}
			connectTasks = append(connectTasks, &types.TaskSpec{
				Type: "connect",
				Params: &tasks.ConnectParams{
					Host:      _h.Agent.Host,
					SSHConfig: _h.Agent.SSHConfig,
				},
			})

			if process.Name == component || component == "all" {
				if process.StartType == "shell" {
					scriptPath := path.Join(utils.GetUtilPath(spec.InstallPath, process.Name), process.ExecShellPath)
					operateTasks = append(operateTasks, &types.TaskSpec{
						Type: "operate",
						Params: &tasks.OperateParams{
							Operation: operation,
							ExecPath:  scriptPath,
							Host:      _h.Agent.Host,
						},
					})
				} else {
					operateTasks = append(operateTasks, &types.TaskSpec{
						Type: "systemd",
						Params: &tasks.SystemdParams{
							Operate: operation,
							Name:    process.Name,
							Host:    _h.Agent.Host,
						},
					})
				}
			}
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
				SubTasks: operateTasks,
			},
		},
	}, nil
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
	if flag && (newHost == "" || utils.GetHostIP(nowHost) == utils.GetHostIP(newHost)) {
		flag = true
	} else {
		flag = false
	}
	if flag {
		if (*allNeedOperations)[nowHost] == nil {
			(*allNeedOperations)[nowHost] = make(map[types.NebulaServiceComponent]bool)
		}
		(*allNeedOperations)[nowHost][nowComponent] = true
	}
	return flag
}

func OperationServiceGroup(spec *types.JobSpec, operation, component, host, KillWait string) (*types.TaskSpec, error) {
	componentType := types.NebulaServiceComponentMap[component]
	if spec.Spec.Metad == nil {
		return nil, fmt.Errorf("metad spec is nil")
	}
	metaHosts := spec.Spec.Metad.Hosts
	allNeedOperations := make(map[string]map[types.NebulaServiceComponent]bool, 0)
	allAgents := utils.GetAllAgents(spec)
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
	operationTasks := []*types.TaskSpec{}
	//1. connect
	agentsMap := make(map[string]*types.Agent)
	for _, agent := range allAgents {
		connectTasks = append(connectTasks, &types.TaskSpec{
			Type:   "connect",
			Params: &tasks.ConnectParams{Host: agent.Host, SSHConfig: agent.SSHConfig},
		})
		agentsMap[agent.Host] = agent
	}
	//2. operation
	for host, components := range allNeedOperations {
		installPath := utils.GetUserCluster(spec.InstallPath, agentsMap[host].InstallPath)
		// aggregate operation task
		if (components[types.Metad] && components[types.Graphd] && components[types.Storaged]) || components[types.AllNebulaSerivce] {
			operationTasks = append(operationTasks, &types.TaskSpec{
				Type: "nebula_operation",
				Params: &tasks.NebulaOperationParams{Operation: operation, Component: types.AllNebulaSerivce, Host: host,
					KillWait: KillWait,
					Path:     installPath,
				},
			})
			continue
		}
		for component := range components {
			operationTasks = append(operationTasks, &types.TaskSpec{
				Type: "nebula_operation",
				Params: &tasks.NebulaOperationParams{Operation: operation, Component: component, Host: host,
					KillWait: KillWait,
					Path:     installPath,
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
				Type:     "parallel",
				SubTasks: operationTasks,
			},
		},
	}, nil
}
