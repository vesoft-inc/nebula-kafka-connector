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
		{"127.0.0.1:9669,127.0.0.2:9669,127.0.0.3:9669a,", nil, `[99000]: address 127.0.0.3:9669a is not valid, strconv.Atoi: parsing "9669a": invalid syntax`},
		{"127.0.0.1:9669,127.0.0.2:9669,127.0.0.39669a,", nil, "[99000]: address 127.0.0.39669a is not valid, address 127.0.0.39669a: missing port in address"},
		{"harris:9669,", []*hostAddress{
			{"harris", 9669},
		}, ""},
		{"[2001:0db8:85a3::8a2e:0370:7334]:9669", []*hostAddress{
			{"2001:0db8:85a3::8a2e:0370:7334", 9669},
		}, ""},
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
