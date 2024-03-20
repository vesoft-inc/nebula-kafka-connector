package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

func TestNewSerial(t *testing.T) {
	Init()
	taskSpec := &types.TaskSpec{
		Type: "serial",
		SubTasks: []*types.TaskSpec{
			{
				Type: "debug",
				Params: &DebugParams{
					Message: "task1",
				},
			},
			{
				Type: "debug",
				Params: &DebugParams{
					Message: "task2",
				},
			},
		},
	}
	jobContext := NewJobContext()
	task, err := NewSerial(taskSpec, jobContext)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	err = task.Execute()
	assert.NoError(t, err)
	err = task.Rollback()
	assert.NoError(t, err)
}

func TestShell(t *testing.T) {
	Init()
	taskSpec := &types.TaskSpec{
		Type: "serial",
		SubTasks: []*types.TaskSpec{
			{
				Type: "connect",
				Params: &ConnectParams{
					Host: "192.168.8.240:6688",
				},
			},
			{
				Type: "shell",
				Params: &ShellParams{
					Host:    "192.168.8.240:6688",
					Sudo:    true,
					Command: "ls",
				},
			},
		},
	}
	executor.SetCertConfig(executor.CertConfig{
		CAFile:  "../../certs/ca.crt",
		KeyFile: "../../certs/ngadmin.key",
		CrtFile: "../../certs/ngadmin.crt",
	})
	jobContext := NewJobContext()
	task, err := NewSerial(taskSpec, jobContext)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	err = task.Execute()
	assert.NoError(t, err)
	err = task.Rollback()
	assert.NoError(t, err)
}
