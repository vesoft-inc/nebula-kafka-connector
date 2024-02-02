package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

type ShellParams struct {
	Host    string
	Command string
	Sudo    bool
	CmdID   string
}

func NewShell(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*ShellParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}

	return &Shell{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		host:       params.Host,
		command:    params.Command,
		sudo:       params.Sudo,
		cmdID:      params.CmdID,
	}, nil
}

type Shell struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	host       string
	command    string
	sudo       bool
	cmdID      string
}

func (d *Shell) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	stdout, stderr, err := executor.Shell(d.command, d.sudo)
	outputID := d.host
	if d.cmdID != "" {
		outputID = d.cmdID
	}
	d.JobContext.SetValue(outputID, map[string]any{
		"stdout": stdout,
		"stderr": stderr,
	})
	return err
}

func (d *Shell) Rollback() error {
	return nil
}

func (d *Shell) String() string {
	return "Shell"
}
