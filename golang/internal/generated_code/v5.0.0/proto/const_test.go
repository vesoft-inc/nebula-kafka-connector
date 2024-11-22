package proto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocalVersion(t *testing.T) {
	assert.Equal(t, PROTOCOL_VERSION, []byte("5.0.0"))
}
