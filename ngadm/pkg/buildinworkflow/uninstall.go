package buildinworkflow

import (
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func Uninstall(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	uninstallTask, err := UninstallServiceGroup(args, spec)
	if err != nil {
		return nil, err
	}
	if uninstallTask != nil {
		workflow.Tasks = append(workflow.Tasks, uninstallTask)
	}
	uninstallTask, err = UninstallUtils(args, spec)
	if err != nil {
		return nil, err
	}
	if uninstallTask != nil {
		workflow.Tasks = append(workflow.Tasks, uninstallTask)
	}
	return workflow, nil
}

func UninstallServiceGroup(args map[string]any, spec *types.JobSpec) (*types.TaskSpec, error) {
	if spec.Spec.Metad == nil {
		return nil, nil
	}
	killWait, ok := args["kill-wait"].(string)
	if !ok {
		killWait = ""
	}
	uninstallTask := &types.TaskSpec{
		Type:     "serial", //all task serial for safe delete data
		SubTasks: []*types.TaskSpec{},
	}
	// stop first
	stopTasks, err := OperationServiceGroup(spec, "stop", "all", "", killWait)
	if err != nil {
		return nil, err
	}
	if stopTasks != nil {
		uninstallTask.SubTasks = append(uninstallTask.SubTasks, stopTasks)
	}
	// delete data
	drain, _ := args["drain"].(bool)
	clusters := spec.Spec.Metad.ServiceGroups
	deletedMap := make(map[string]*types.TaskSpec)
	for _, cluster := range clusters {
		storaged := cluster.Storaged
		for _, storage := range storaged.Hosts {
			task := &types.TaskSpec{
				Type: "delete_nebula_data",
				Params: &tasks.DeleteNebulaDataParams{
					Host:  storage.Agent.Host,
					Path:  utils.GetCluster(spec.InstallPath),
					Drain: drain,
				},
			}
			deletedMap[storage.Agent.Host] = task
		}
		for _, graphd := range cluster.Graphd.Hosts {
			if _, ok := deletedMap[graphd.Agent.Host]; ok {
				continue
			}
			task := &types.TaskSpec{
				Type: "delete_nebula_data",
				Params: &tasks.DeleteNebulaDataParams{
					Host:  graphd.Agent.Host,
					Path:  utils.GetCluster(spec.InstallPath),
					Drain: true, //for graph delete all data
				},
			}
			deletedMap[graphd.Agent.Host] = task
		}
	}
	for _, metad := range spec.Spec.Metad.Hosts {
		if _, ok := deletedMap[metad.Agent.Host]; ok {
			continue
		}
		task := &types.TaskSpec{
			Type: "delete_nebula_data",
			Params: &tasks.DeleteNebulaDataParams{
				Host:  metad.Agent.Host,
				Path:  utils.GetCluster(spec.InstallPath),
				Drain: drain,
			},
		}
		deletedMap[metad.Agent.Host] = task
	}
	for _, task := range deletedMap {
		uninstallTask.SubTasks = append(uninstallTask.SubTasks, task)
	}
	allHosts := GetMetadAllNeedHosts(spec)
	// delete download path
	for _, host := range allHosts {
		task := &types.TaskSpec{
			Type: "rm",
			Params: &tasks.RMParams{
				Host: host.Host,
				Path: utils.GetDownloadPath(spec.InstallPath),
			},
		}
		uninstallTask.SubTasks = append(uninstallTask.SubTasks, task)
	}
	return uninstallTask, nil
}
