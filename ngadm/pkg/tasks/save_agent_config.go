package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type SaveAgentParams struct {
	Config    map[string]any
	Component string
}

func NewSaveAgent(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*SaveAgentParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &SaveAgent{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		config:     params.Config,
		component:  params.Component,
	}, nil
}

type SaveAgent struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	host       string
	config     map[string]any
	component  string
}

func (d *SaveAgent) Execute() error {
	exec := d.JobContext.GetExecuter(d.host)
	if exec == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	agentExec := exec.(*executor.AgentExecutor)
	agentExec.SaveAgentConfig(d.component, d.config)
	return nil
}

func (d *SaveAgent) Rollback() error {
	return nil
}

func (d *SaveAgent) String() string {
	return "SaveAgent-" + d.component
}
