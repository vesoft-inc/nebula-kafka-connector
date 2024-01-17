package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCode(t *testing.T) {
	leaderChangeCode := 0x343030444e58
	assert.Equal(t, uint64(leaderChangeCode), ErrorLeaderChange.Code())
}
