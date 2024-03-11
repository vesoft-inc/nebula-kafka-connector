package tasks

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/utils"
)

type StatusParams struct {
	Host          string
	ExecShellPath string
	Name          string
	Port          string
	Component     types.NebulaServiceComponent // graphd, metad, storaged
}

func NewStatus(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*StatusParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &Status{
		JobContext:    taskContext,
		taskSpec:      taskSpec,
		host:          params.Host,
		execShellPath: params.ExecShellPath,
		component:     params.Component,
		name:          params.Name,
		port:          params.Port,
	}, nil
}

type Status struct {
	JobContext    *JobContext
	taskSpec      *types.TaskSpec
	host          string
	execShellPath string
	component     types.NebulaServiceComponent
	name          string
	port          string
}

func (d *Status) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	cmd := fmt.Sprintf("%s %s %s", d.execShellPath, "status", d.component)
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil {
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}
	str := stdout
	if len(stdout) == 0 {
		str = stderr
	}
	statusLines := strings.Split(str, "\n")
	for _, line := range statusLines {
		if line == "" {
			continue
		}
		status := matchProcessStatus(line)
		name := matchProcessName(line)
		port := matchProcessPort(line)
		if port == "" {
			port = d.port
		}

		if name == "unknown" || status == "unknown" {
			continue
		}
		d.JobContext.SetValue("status-"+utils.GetHostIP(d.host)+"-"+name+"-"+port, types.StatusItem{
			Product: d.name,
			Service: name,
			Host:    utils.GetHostIP(d.host),
			Port:    port,
			Status:  status,
		})
	}
	return nil
}

func (d *Status) Rollback() error {
	return nil
}

func (d *Status) String() string {
	return "Status-" + d.component.String()
}
