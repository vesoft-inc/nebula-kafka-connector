package nebula_ng

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/common"
)

type nullableValue interface {
	scan(value Value) error
	getData() any
	isValid() bool
}

type testcase struct {
	value    common.Value
	expected interface{}
	scanner  nullableValue
}

func TestScan(t *testing.T) {
	testcases := []testcase{
		{
			value: common.Value{
				Data: &common.Value_StringValue{
					StringValue: []byte("test"),
				},
			},
			expected: String("test"),
			scanner:  &NullString{},
		},
		{
			value: common.Value{
				Data: &common.Value_StringValue{
					StringValue: []byte("test"),
				},
			},
			expected: String("test"),
			scanner:  &NullString{},
		},
		{
			value: common.Value{
				Data: &common.Value_BoolValue{
					BoolValue: true,
				},
			},
			expected: Bool(true),
			scanner:  &NullBool{},
		},
		{
			value: common.Value{
				Data: &common.Value_Int8Value{
					Int8Value: 1,
				},
			},
			expected: Int64(1),
			scanner:  &NullInt{},
		},
	}

	for _, tc := range testcases {
		v := grpcValue{data: &tc.value}
		if err := tc.scanner.scan(&v); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, tc.expected, tc.scanner.getData())
		assert.Equal(t, true, tc.scanner.isValid())
	}
}

func TestNode(t *testing.T) {
	value := common.Value{
		Data: &common.Value_NodeValue{
			NodeValue: &common.Node{
				Labels: []string{"test"},
				Properties: map[string]*common.Value{
					"address": {
						Data: &common.Value_StringValue{
							StringValue: []byte("this is an address"),
						},
					},
				},
				NodeId: 999,
				Graph:  "graph_name",
				Type:   "hhh",
			},
		},
	}
	scanner := &NullNode{}
	v := grpcValue{data: &value}
	if err := scanner.scan(&v); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, scanner.isValid())
	n, ok := scanner.getData().(Node)
	assert.Equal(t, true, ok)
	assert.Equal(t, "test", n.GetLabels()[0])
	addr, err := n.GetProperties()["address"].AsString()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, String("this is an address"), addr)
	assert.Equal(t, int64(999), n.GetId())
	assert.Equal(t, "graph_name", n.GetGraph())
	assert.Equal(t, "hhh", n.GetType())
}

func TestNull(t *testing.T) {
	testcases := []testcase{
		{
			value: common.Value{
				Data: nil,
			},
			expected: String("test"),
			scanner:  &NullString{},
		},
		{
			value: common.Value{
				Data: nil,
			},
			expected: String("test"),
			scanner:  &NullNode{},
		},
	}
	for _, tc := range testcases {
		v := grpcValue{data: &tc.value}
		if err := tc.scanner.scan(&v); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, false, tc.scanner.isValid())
	}
}

func TestTypeError(t *testing.T) {
	testcases := []testcase{
		{
			value: common.Value{
				Data: &common.Value_StringValue{
					StringValue: []byte("test"),
				},
			},
			expected: String("test"),
			scanner:  &NullNode{},
		},
	}
	for _, tc := range testcases {
		v := grpcValue{data: &tc.value}
		err := tc.scanner.scan(&v)
		assert.NotNil(t, err)
	}
}
