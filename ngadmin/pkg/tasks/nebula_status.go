package tasks

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/utils"
)

type NebulaStatusParams struct {
	Host      string
	Path      string
	Component types.NebulaServiceComponent // graphd, metad, storaged
}

func NewNebulaStatus(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*NebulaStatusParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &NebulaStatus{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		host:       params.Host,
		path:       params.Path,
		component:  params.Component,
	}, nil
}

type NebulaStatus struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	host       string
	path       string
	component  types.NebulaServiceComponent
}

func (d *NebulaStatus) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	cmd := fmt.Sprintf("%s %s %s", path.Join(d.path, "scripts/nebula.service"), "status", d.component)
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}
	statusLines := strings.Split(stdout, "\n")
	for _, line := range statusLines {
		status := matchProcessStatus(line)
		name := matchProcessName(line)
		port := matchProcessPort(line)
		if name == "unknown" || status == "unknown" {
			continue
		}
		d.JobContext.SetValue("status-"+utils.GetHostIP(d.host)+"-"+name+"-"+port, types.StatusItem{
			Product: "nebulagraph",
			Service: name,
			Host:    utils.GetHostIP(d.host),
			Port:    port,
			Status:  status,
		})
	}
	return nil
}

func (d *NebulaStatus) Rollback() error {
	return nil
}

func (d *NebulaStatus) String() string {
	return "NebulaStatus-" + d.component.String()
}

func matchProcessName(line string) string {
	name := "unknown"
	r := regexp.MustCompile(`\[\S+\] (\S+)`)
	match := r.FindStringSubmatch(line)
	if len(match) > 1 {
		name = match[1]
	}
	if strings.Contains(line, "nebula-graphd") {
		name = "graphd"
	} else if strings.Contains(line, "nebula-metad") {
		name = "metad"
	} else if strings.Contains(line, "nebula-storaged") {
		name = "storaged"
	}
	return name
}

func matchProcessStatus(line string) string {
	status := "unknown"
	if strings.Contains(strings.ToLower(line), "running") {
		status = "running"
	} else if strings.Contains(strings.ToLower(line), "exited") {
		status = "exited"
	}
	return status
}

// [INFO] nebula-graphd(302572cd): Running as 16312, Listening on 9669
func matchProcessPort(line string) string {
	strs := strings.Split(line, "Listening on ")
	if len(strs) < 2 {
		return ""
	}
	port := strs[1]
	return port
}
