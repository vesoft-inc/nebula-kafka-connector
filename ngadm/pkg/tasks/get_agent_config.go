package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type GetAgentParams struct {
	Config    map[string]any
	Component string
	Name      string
}

func NewGetAgent(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*GetAgentParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	if params.Name == "" {
		params.Name = params.Component + "-agent-config"
	}
	return &GetAgent{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		config:     params.Config,
		component:  params.Component,
		name:       params.Name,
	}, nil
}

type GetAgent struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	host       string
	config     map[string]any
	component  string
	name       string
}

func (d *GetAgent) Execute() error {
	exec := d.JobContext.GetExecuter(d.host)
	if exec == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	agentExec := exec.(*executor.AgentExecutor)
	res, err := agentExec.GetAgentConfigg(d.component)
	if err != nil {
		return err
	}
	d.JobContext.SetValue(d.name, res)
	return nil
}

func (d *GetAgent) Rollback() error {
	return nil
}

func (d *GetAgent) String() string {
	return "GetAgent-" + d.component
}
