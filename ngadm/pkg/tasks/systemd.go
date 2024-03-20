package tasks

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

type SystemdParams struct {
	Host             string
	ExecStartPath    string //needed for install
	WorkingDirectory string //needed for install
	Name             string
	Operate          string // start, stop, restart, install, uninstall, status
	Port             string
}

func NewSystemd(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*SystemdParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &Systemd{
		JobContext:       taskContext,
		taskSpec:         taskSpec,
		host:             params.Host,
		execStartPath:    params.ExecStartPath,
		workingDirectory: params.WorkingDirectory,
		name:             params.Name,
		operate:          params.Operate,
		port:             params.Port,
	}, nil
}

type Systemd struct {
	JobContext       *JobContext
	taskSpec         *types.TaskSpec
	host             string
	execStartPath    string
	workingDirectory string
	name             string
	operate          string
	port             string
}

func (d *Systemd) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	if d.operate == "install" {
		err := d.Install()
		if err != nil {
			return err
		}
		return d.Operate("start")
	}
	if d.operate == "uninstall" {
		return d.Uninstall()
	}
	if d.operate == "status" {
		return d.Status()
	}

	return d.Operate(d.operate)
}

func (d *Systemd) Operate(operate string) error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	cmd := fmt.Sprintf("systemctl %s %s", operate, d.name)
	stdout, stderr, err := executor.Shell(cmd, true)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}
	return nil
}

func (d *Systemd) Install() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	serviceFileContent := fmt.Sprintf(`
[Unit]
Description=%s

[Service]
Type=simple
ExecStart=%s
WorkingDirectory=%s
Restart=always

[Install]
WantedBy=multi-user.target`, d.name, d.execStartPath, d.workingDirectory)
	cmd := "mkdir -p ~/.config/systemd/user/"
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}

	cmd = fmt.Sprintf("cat <<EOF > ~/.config/systemd/user/%s.service \n %s\nEOF", d.name, serviceFileContent)
	stdout, stderr, err = executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}
	cmd = fmt.Sprintf("systemctl --user daemon-reload && systemctl --user enable %s ", d.name)
	stdout, stderr, err = executor.Shell(cmd, false)
	if err != nil || (len(stderr) > 0 && !strings.Contains("Created symlink", stderr)) {
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}
	return nil
}

func (d *Systemd) Uninstall() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	cmd := fmt.Sprintf("systemctl --user stop %s", d.name)
	executor.Shell(cmd, true)
	cmd = fmt.Sprintf("systemctl --user disable %s", d.name)
	executor.Shell(cmd, true)
	cmd = fmt.Sprintf("rm -f ~/.config/systemd/user/%s.service", d.name)
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}
	return nil
}

func (d *Systemd) Status() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	cmd := fmt.Sprintf("systemctl --user status %s", d.name)
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil {
		return fmt.Errorf("failed to execute cmd: %s, err: %s, stdout: %s, stderr: %s", cmd, err, stdout, stderr)
	}
	status := "unknown"
	if !strings.Contains(stdout, "failed") {
		status = "running"
	} else if len(stderr) > 0 {
		status = "exited"
	}
	d.JobContext.SetValue("status-"+utils.GetHostIP(d.host)+"-"+d.name+"-"+d.port, types.StatusItem{
		Product: d.name,
		Service: d.name,
		Host:    utils.GetHostIP(d.host),
		Port:    d.port,
		Status:  status,
	})
	return nil
}

func (d *Systemd) Rollback() error {
	d.Uninstall()

	return nil
}

func (d *Systemd) String() string {
	return "Systemd-install-" + d.name
}
