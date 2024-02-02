package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
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
