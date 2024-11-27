package buildinworkflow

import (
	"fmt"
	"path/filepath"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

func UninstallHost(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	if args == nil {
		args = map[string]any{}
	}
	uninstallAll, ok := args["uninstallAll"].(bool)
	if !ok {
		uninstallAll = true
	}
	drain, ok := args["drain"].(bool)
	if !ok {
		drain = true
	}
	var allNeedHosts map[string]*types.Agent
	if uninstallAll {
		allNeedHosts = GetMetadAllNeedHosts(spec)
	} else {
		selectedHost := args["selectedHost"].([]string)
		allNeedHosts = GetMetadSelectedHosts(selectedHost, spec)
	}
	connectTasks := &types.TaskSpec{
		Type:     "parallel",
		SubTasks: make([]*types.TaskSpec, 0),
	}
	uninstallTasks := &types.TaskSpec{
		Type:     "parallel",
		SubTasks: make([]*types.TaskSpec, 0),
	}
	deleteTasks := &types.TaskSpec{
		Type:     "parallel",
		SubTasks: make([]*types.TaskSpec, 0),
	}
	for _, agent := range allNeedHosts {
		installPath := utils.GetUserCluster(spec.InstallPath, agent.InstallPath)
		// connect and uninstall
		connectTasks.SubTasks = append(connectTasks.SubTasks, &types.TaskSpec{
			Type: "connect",
			Params: &tasks.ConnectParams{
				Host:      agent.Host,
				SSHConfig: agent.SSHConfig,
			},
		})
		var uninstallCmd string
		if drain {
			// if drain, we need to ignore the error of uninstall
			uninstallCmd = fmt.Sprintf("echo Y | %s || echo 1", filepath.Join(installPath, "scripts", "nebula.uninstall"))
		} else {
			uninstallCmd = fmt.Sprintf("echo Y | " + filepath.Join(installPath, "scripts", "nebula.uninstall"))
		}
		uninstallTasks.SubTasks = append(uninstallTasks.SubTasks, &types.TaskSpec{
			Type: "shell",
			Params: &tasks.ShellParams{
				Host:    agent.Host,
				Command: uninstallCmd,
			},
		})
		var deleteCmd string
		if drain {
			if installPath == "" {
				return nil, fmt.Errorf("install path is empty")
			}
			deleteCmd = fmt.Sprintf("rm -rf %s &&", installPath)
		}
		deleteCmd += fmt.Sprintf("rm -f %s", filepath.Join(INSTALL_PATH_DIR, INSTALL_PATH_FILE))

		deleteTasks.SubTasks = append(deleteTasks.SubTasks, &types.TaskSpec{
			Type: "shell",
			Params: &tasks.ShellParams{
				Host:    agent.Host,
				Command: deleteCmd,
			},
		})
	}

	workflow.Tasks = []*types.TaskSpec{connectTasks, uninstallTasks, deleteTasks}

	return workflow, nil
}
