package tasks

import (
	"fmt"
	"os/exec"
	"path"
	"syscall"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

type NebulaOperationParams struct {
	Host         string
	Path         string
	Operation    string
	Component    types.NebulaServiceComponent // graphd, metad, storaged
	NeedRollback bool                         //for install, we need rollback to stop process
	KillWait     string
}

func NewNebulaOperation(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*NebulaOperationParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &NebulaOperation{
		JobContext:   taskContext,
		taskSpec:     taskSpec,
		host:         params.Host,
		path:         params.Path,
		operation:    params.Operation,
		component:    params.Component,
		needRollback: params.NeedRollback,
		killWait:     params.KillWait,
	}, nil
}

type NebulaOperation struct {
	JobContext   *JobContext
	taskSpec     *types.TaskSpec
	host         string
	path         string
	operation    string
	component    types.NebulaServiceComponent
	needRollback bool
	killWait     string
}

func (d *NebulaOperation) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	cmd := fmt.Sprintf("%s %s %s", path.Join(d.path, "scripts/nebula.service"), string(d.operation), d.component)
	// add timeout for kill
	if d.killWait != "" && d.operation == "stop" {
		cmd = fmt.Sprintf("timeout %s %s", d.killWait, cmd)
	}
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 124 && d.operation == "stop" {
					return d.Kill() // kill the process if timeout for stop
				}
			}
		}
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}
	return nil
}

func (d *NebulaOperation) Kill() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	cmd := fmt.Sprintf("%s %s %s", path.Join(d.path, "scripts/nebula.service"), "kill", d.component)
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}
	return nil
}

func (d *NebulaOperation) Rollback() error {
	if !d.needRollback {
		return nil
	}
	operation := ""
	if d.operation == "start" {
		operation = "stop"
	} else if d.operation == "stop" {
		operation = "start"
	}
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}

	cmd := fmt.Sprintf("%s %s %s", path.Join(d.path, "scripts/nebula.service"), string(operation), d.component)
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil {
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}

	return nil
}

func (d *NebulaOperation) String() string {
	return "NebulaOperation-" + d.operation + "-" + d.component.String()
}
