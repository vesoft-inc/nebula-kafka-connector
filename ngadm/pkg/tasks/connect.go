package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type ConnectParams struct {
	Host      string
	SSHConfig *types.SSHConfig
	Typ       string
}
type Connect struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	host       string
	typ        string
	sshConfig  *types.SSHConfig
}

func NewConnect(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*ConnectParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	var typ = "agent"
	if params.SSHConfig != nil {
		typ = "ssh"
	}
	return &Connect{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		host:       params.Host,
		typ:        typ,
		sshConfig:  params.SSHConfig,
	}, nil // Change Debug{} to &Debug{}
}

func (d *Connect) Execute() error { // Change (d *Debug) to (d *Debug)
	var instance executor.Executor
	var err error
	if d.typ == "ssh" {
		if d.sshConfig.Host == "" {
			d.sshConfig.Host = d.host
		}
		if d.sshConfig.Port == 0 {
			d.sshConfig.Port = 22
		}
		instance, err = executor.NewSSHExecuter(d.sshConfig)
		if err != nil {
			return err
		}
	} else {
		instance, err = executor.NewAgentExecuter(d.host, 60)
		if err != nil {
			return err
		}
	}
	d.JobContext.SetExecuter(d.host, instance)
	return nil
}

func (d *Connect) Rollback() error {
	return nil
}

func (d *Connect) String() string {
	return "Connect " + d.host
}
