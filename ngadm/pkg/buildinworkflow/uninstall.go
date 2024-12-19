package buildinworkflow

import (
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

func Uninstall(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	workflow := &types.WorkflowSpec{
		Type:     "serial",
		Rollback: spec.Rollback,
		Tasks:    []*types.TaskSpec{},
	}
	uninstallTask, err := BuildInstallServiceGroupTask(args, spec)
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
