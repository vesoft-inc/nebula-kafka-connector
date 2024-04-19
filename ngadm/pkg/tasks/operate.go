package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type OperateParams struct {
	Host      string
	ExecPath  string
	Operation string
	Component string
}

func NewOperate(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*OperateParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	if params.Component == "" {
		params.Component = "all"
	}
	return &Operate{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		host:       params.Host,
		execPath:   params.ExecPath,
		operation:  params.Operation,
		component:  params.Component,
	}, nil
}

type Operate struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	host       string
	execPath   string
	operation  string
	component  string
}

func (d *Operate) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	cmd := fmt.Sprintf("%s %s %s", d.execPath, d.operation, d.component)
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("execute command failed: %s, %s, %s", cmd, stdout, stderr)
	}
	d.JobContext.Logger.Info(fmt.Sprintf("%s: %s", cmd, stdout))
	return nil
}

func (d *Operate) Rollback() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	rollbackOperation := ""
	if d.operation == "start" { //now only install will be rollback
		rollbackOperation = "stop"
	}
	cmd := fmt.Sprintf("%s %s %s", d.execPath, rollbackOperation, d.component)
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("execute command failed: %s, %s, %s", cmd, stdout, stderr)
	}
	d.JobContext.Logger.Info(fmt.Sprintf("%s:\n %s", cmd, stdout))
	return nil
}

func (d *Operate) String() string {
	return fmt.Sprintf("Operate: %s-%s-%s", d.host, d.operation, d.component)
}
