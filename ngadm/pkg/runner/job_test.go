package runner

import (
	"log"
	"testing"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

// !!! test need start a agent server
func TestRunWorkflow(t *testing.T) {
	job := NewJob("test")

	workflow := &types.WorkflowSpec{
		Tasks: []*types.TaskSpec{
			{Type: "debug", Params: make(map[string]any)},
			{Type: "task2", Params: make(map[string]any)},
		},
	}

	err := job.RunWorkflow(workflow)
	if err == nil {
		t.Errorf("RunTask failed: %v", err)
	}
}

func TestRunTask(t *testing.T) {
	job := NewJob("test")

	task := &types.TaskSpec{Type: "debug", Params: make(map[string]any)}

	taskInstance, err := job.RunTask(task)
	if err != nil {
		t.Errorf("RunTask failed: %v", err)
	}
	err = taskInstance.Rollback()
	if err != nil {
		t.Errorf("Rollback failed: %v", err)
	}
}

func TestAgent(t *testing.T) {
	executor.SetCertConfig(executor.CertConfig{
		CAFile:  "../../certs/ca.crt",
		KeyFile: "../../certs/client.key",
		CrtFile: "../../certs/client.crt",
	})
	job := NewJob("test")

	workflow := &types.WorkflowSpec{
		Tasks: []*types.TaskSpec{
			{Type: "connect", Params: &tasks.ConnectParams{Host: "https://127.0.0.1:6688"}},
			{Type: "shell", Params: &tasks.ShellParams{
				Host:    "https://127.0.0.1:6688",
				Command: "ls",
				Sudo:    false,
				CmdID:   "lscmd",
			}},
		},
	}
	err := job.RunWorkflow(workflow)
	if err != nil {
		t.Errorf("RunTask failed: %v", err)
	}
	log.Print(job.Context.ValueMap["lscmd"])
}

func TestSSH(t *testing.T) {
	job := NewJob("test")
	workflow := &types.WorkflowSpec{
		Tasks: []*types.TaskSpec{
			{Type: "connect", Params: &tasks.ConnectParams{
				Host: "192.168.8.240",
				SSHConfig: &types.SSHConfig{
					Host:     "192.168.8.240",
					Port:     22,
					User:     "vesoft",
					Password: "nebula",
				},
			}},
			{Type: "shell", Params: &tasks.ShellParams{
				Host:    "192.168.8.240",
				Command: "ls",
				Sudo:    false,
				CmdID:   "lscmd",
			}},
		},
	}
	err := job.RunWorkflow(workflow)
	if err != nil {
		t.Errorf("RunTask failed: %v", err)
	}
	log.Print(job.Context.ValueMap["lscmd"])
}

func TestUpload(t *testing.T) {
	executor.SetCertConfig(executor.CertConfig{
		CAFile:  "../../certs/ca.crt",
		KeyFile: "../../certs/client.key",
		CrtFile: "../../certs/client.crt",
	})
	job := NewJob("test")
	workflow := &types.WorkflowSpec{
		Tasks: []*types.TaskSpec{
			{Type: "connect", Params: &tasks.ConnectParams{Host: "https://127.0.0.1:6688"}},
			{Type: "upload", Params: &tasks.UploadParams{
				SrcPath: "../../certs/client.crt",
				DstPath: "/Users/mizy/projects/nebula-ng-tools/agent/api/agent/certs/",
				Host:    "https://127.0.0.1:6688",
			}},
		},
	}
	err := job.RunWorkflow(workflow)
	if err != nil {
		t.Errorf("RunTask failed: %v", err)
	}
	log.Print(job.Context.ValueMap)
}
