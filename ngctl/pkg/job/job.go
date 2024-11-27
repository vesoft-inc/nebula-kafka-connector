package job

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
)

type Job interface {
	UpCluster(*config.Config, string) error
	DownCluster(*config.Config, bool) error
	ServiceOperation(c *config.Config, hosts []string, serviceType string, operation string) error
	UpdateConfig(c *config.Config, serviceGroup string, hosts []string, serviceType string) error
	InstallHost(*config.Config, string) error
	UninstallHost(*config.Config, string, bool) error
}

type ngadmJobType struct{}

var NgadmJob Job = &ngadmJobType{}

const defaultMetadUser = "root"
const defaultMetadPassword = "nebula"

func newJob(name string, c *config.Config) *runner.Job {
	job := runner.NewJob(name)
	executor.SetCertConfig(executor.CertConfig{
		CAFile:  c.JobSpec.CAFile,
		CrtFile: c.JobSpec.CertFile,
		KeyFile: c.JobSpec.KeyFile,
	})
	return job
}

func (j *ngadmJobType) UpCluster(c *config.Config, metaPassword string) error {
	job := newJob("up", c)
	args := map[string]any{
		"username":     defaultMetadUser,
		"password":     defaultMetadPassword,
		"metaPassword": metaPassword,
	}
	if err := job.Run("install", args, c.JobSpec); err != nil {
		return err
	}
	return nil
}

func (j *ngadmJobType) ServiceOperation(c *config.Config, hosts []string, serviceType string, operation string) error {
	if operation != "start" && operation != "stop" && operation != "restart" && operation != "status" {
		return common.NgctlError("invalid operation on service", "")
	}
	if !isValidServiceType(serviceType) {
		return common.NgctlError("invalid service type", "")
	}
	cmd := exec.Command(filepath.Join(c.InstallPath, "scripts/nebula.service"), operation, serviceType)
	log.Printf("exec cmd: %v, on %+v", cmd, hosts)
	job := newJob(fmt.Sprintf("%s_service", operation), c)

	workflow := &types.WorkflowSpec{}
	connectTasks := &types.TaskSpec{
		Type:     "parallel",
		SubTasks: make([]*types.TaskSpec, 0),
	}
	shellTasks := &types.TaskSpec{
		Type:     "parallel",
		SubTasks: make([]*types.TaskSpec, 0),
	}
	for _, h := range hosts {
		inst, err := c.GetInstanceFromHost(h)
		if err != nil {
			return err
		}
		addr := fmt.Sprintf("%s:%d", inst.Host, inst.AgentPort)
		connectTasks.SubTasks = append(connectTasks.SubTasks, &types.TaskSpec{
			Type: "connect",
			Params: &tasks.ConnectParams{
				Host: addr,
			},
		})
		shellTasks.SubTasks = append(shellTasks.SubTasks, &types.TaskSpec{
			Type: "shell",
			Params: &tasks.ShellParams{
				Host:    addr,
				Command: cmd.String(),
				Sudo:    false,
			},
		})
	}
	workflow.Tasks = []*types.TaskSpec{connectTasks, shellTasks}

	if err := job.RunWorkflow(workflow); err != nil {
		return err
	}
	return nil
}

func (j *ngadmJobType) DownCluster(c *config.Config, drain bool) error {
	// 1. stop all services
	// 2. uninstall all services
	// 3. if drain, delete all data in install path
	var hostMap = make(map[string]struct{})
	var graphdHosts = make([]string, 0)
	var storagedHosts = make([]string, 0)
	var metadHosts = make([]string, 0)

	for _, sg := range c.Spec.ServiceGroups {
		for _, inst := range sg.Graphd.Instances {
			graphdHosts = append(graphdHosts, inst.Host)
			hostMap[inst.Host] = struct{}{}
		}
		for _, inst := range sg.Storaged.Instances {
			storagedHosts = append(storagedHosts, inst.Host)
			hostMap[inst.Host] = struct{}{}
		}
	}
	for _, inst := range c.Spec.Metad.Instances {
		metadHosts = append(metadHosts, inst.Host)
		hostMap[inst.Host] = struct{}{}
	}
	hosts := make([]string, 0)
	for h := range hostMap {
		hosts = append(hosts, h)
	}
	// stop all services
	_ = j.ServiceOperation(c, hosts, "graphd", "stop")
	_ = j.ServiceOperation(c, hosts, "storaged", "stop")
	_ = j.ServiceOperation(c, hosts, "metad", "stop")

	job := newJob("down", c)
	args := map[string]any{
		"drain":        drain,
		"uninstallAll": false,
		"selectedHost": hosts,
	}
	if err := job.Run("uninstall-host", args, c.JobSpec); err != nil {
		return err
	}
	return nil
}

func (j *ngadmJobType) UpdateConfig(c *config.Config, serviceGroup string, hosts []string, serviceType string) error {
	if !isValidServiceType(serviceType) {
		return common.NgctlError("invalid service type", "")
	}
	job := newJob("update_config", c)
	workflow := &types.WorkflowSpec{}
	connectTasks := &types.TaskSpec{
		Type:     "parallel",
		SubTasks: make([]*types.TaskSpec, 0),
	}
	updateTasks := &types.TaskSpec{
		Type:     "parallel",
		SubTasks: make([]*types.TaskSpec, 0),
	}
	for _, h := range hosts {
		inst, err := c.GetInstanceFromHost(h)
		if err != nil {
			return err
		}
		addr := fmt.Sprintf("%s:%d", inst.Host, inst.AgentPort)
		connectTasks.SubTasks = append(connectTasks.SubTasks, &types.TaskSpec{
			Type: "connect",
			Params: &tasks.ConnectParams{
				Host: addr,
			},
		})
		cm, err := getServiceConfig(c, serviceGroup, serviceType, h)
		if err != nil {
			return err
		}
		updateTasks.SubTasks = append(updateTasks.SubTasks, &types.TaskSpec{
			Type: "init_config",
			Params: &tasks.InitConfigParams{
				Host:      addr,
				ChangeMap: cm,
				Dst:       path.Join(filepath.Join(c.InstallPath, fmt.Sprintf("etc/nebula-%s.conf", serviceType))),
			},
		})
	}
	workflow.Tasks = []*types.TaskSpec{connectTasks, updateTasks}
	if err := job.RunWorkflow(workflow); err != nil {
		return err
	}
	return nil
}

func (j *ngadmJobType) InstallHost(c *config.Config, host string) error {
	job := newJob("install-host", c)
	args := map[string]any{
		"drain":        false,
		"installAll":   false,
		"selectedHost": []string{host},
	}
	if err := job.Run("install-host", args, c.JobSpec); err != nil {
		return err
	}
	return nil
}

func (j *ngadmJobType) UninstallHost(c *config.Config, host string, drain bool) error {
	// 1. stop all services in the host
	// 2. uninstall all services
	// 3. if drain, delete all data in install path
	_ = j.ServiceOperation(c, []string{host}, "graphd", "stop")
	_ = j.ServiceOperation(c, []string{host}, "storaged", "stop")
	job := newJob("uninstall-host", c)
	args := map[string]any{
		"drain":        drain,
		"uninstallAll": false,
		"selectedHost": []string{host},
	}
	if err := job.Run("uninstall-host", args, c.JobSpec); err != nil {
		return err
	}
	return nil
}
