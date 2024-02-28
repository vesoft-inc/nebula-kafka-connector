package nebula_ng

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestErrorCode(t *testing.T) {
	assert.Equal(t, ErrorLeaderChange.codeInt(), uint64(228441736270))
}

func TestErrorFromInt(t *testing.T) {
	c := uint64(228441736270)
	assert.Equal(t, ErrorLeaderChange, ErrorFromInt(c))
}
