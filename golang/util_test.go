package nebula_ng

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHost(t *testing.T) {
	testcases := []struct {
		addresses string
		expected  []*hostAddress
		err       string
	}{
		{"127.0.0.1:9669,127.0.0.2:9669,127.0.0.3:9669,", []*hostAddress{
			{"127.0.0.1", 9669},
			{"127.0.0.2", 9669},
			{"127.0.0.3", 9669},
		}, ""},
		{"127.0.0.1:9669,127.0.0.2:9669,127.0.0.3:9669a,", nil, "address 127.0.0.3:9669a is not valid, port is not valid"},
		{"127.0.0.1:9669,127.0.0.2:9669,127.0.0.39669a,", nil, "address 127.0.0.39669a is not valid"},
	}
	for _, tc := range testcases {
		actual, err := parseAddresses(tc.addresses)
		if err != nil {
			assert.Equal(t, tc.err, err.Error())
		} else {
			assert.Equal(t, tc.err, "")
		}
		assert.Equal(t, tc.expected, actual)
	}
}
