package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

type ConnectParams struct {
	Host string
}

func NewConnect(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*ConnectParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &Connect{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		host:       params.Host,
	}, nil // Change Debug{} to &Debug{}
}

type Connect struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	host       string
}

func (d *Connect) Execute() error { // Change (d *Debug) to (d *Debug)
	executor, err := executor.NewAgentExecuter(d.host, 60)
	if err != nil {
		return err
	}
	d.JobContext.SetExecuter(d.host, executor)
	return nil
}

func (d *Connect) Rollback() error {
	return nil
}

func (d *Connect) String() string {
	return "Connect " + d.host
}
