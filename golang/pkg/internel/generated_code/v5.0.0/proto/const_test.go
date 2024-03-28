package proto

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestProtocalVersion(t *testing.T) {
	assert.Equal(t, PROTOCOL_VERSION, []byte("5.0.0"))
}
