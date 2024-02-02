package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

func TestNewParallel(t *testing.T) {
	taskSpec := &types.TaskSpec{
		SubTasks: []*types.TaskSpec{
			// Add subtask specifications here
		},
	}

	taskContext := &JobContext{
		// Add task context initialization here
	}

	task, err := NewParallel(taskSpec, taskContext)
	assert.NoError(t, err)
	assert.NotNil(t, task)

	// Add assertions for the task properties here
}
