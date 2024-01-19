/*
 *
 * Copyright (c) 2020 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 *
 */

package nebula_ng

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula"
)

var testTimezone timezoneInfo = timezoneInfo{0, []byte("UTC")}

// TODO(Aiee) To add more
func TestAsNull(t *testing.T) {
	null := nebula.XValue_{}
	valWrap := ValueWrapper{&null, testTimezone}
	assert.Equal(t, "__NULL__", valWrap.String())
	assert.Equal(t, null, *valWrap.value)
}

func TestAsBool(t *testing.T) {
	bval := new(bool)
	*bval = true
	value := nebula.XValue_{BoolVal: bval}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, true, valWrap.IsBool())
	assert.Equal(t, "true", valWrap.String())
	res, _ := valWrap.AsBool()
	assert.Equal(t, value.GetBoolVal(), res)
}

func TestAsInt(t *testing.T) {
	val := new(int64)
	*val = 100
	value := nebula.XValue_{Int64Val: val}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, "int64", valWrap.GetType())
	assert.Equal(t, "100", valWrap.String())
	res, _ := valWrap.AsInt64()
	assert.Equal(t, value.GetInt64Val(), res)
}

func TestAsFloat(t *testing.T) {
	val := new(float64)
	*val = 100.111
	value := nebula.XValue_{FloatVal: val}
	valWrap := ValueWrapper{&value, testTimezone}
	val2 := new(float64)
	*val2 = 100.00
	value2 := nebula.XValue_{FloatVal: val2}
	valWrap2 := ValueWrapper{&value2, testTimezone}
	assert.Equal(t, "100.111", valWrap.String())
	assert.Equal(t, "100.0", valWrap2.String())
	assert.Equal(t, true, valWrap.IsFloat())
	res, _ := valWrap.AsFloat()
	assert.Equal(t, value.GetFloatVal(), res)
}

func TestAsString(t *testing.T) {
	{
		val := "test_string"
		value := nebula.XValue_{StringVal: []byte(val)}
		valWrap := ValueWrapper{&value, testTimezone}
		assert.Equal(t, true, valWrap.IsString())
		assert.Equal(t, "\"test_string\"", valWrap.String())
		res, _ := valWrap.AsString()
		assert.Equal(t, string(value.GetStringVal()), res)
	}
	{
		val := "abc\rde"
		value := nebula.XValue_{StringVal: []byte(val)}
		valWrap := ValueWrapper{&value, testTimezone}
		assert.Equal(t, true, valWrap.IsString())
		assert.Equal(t, "\"dec\"", valWrap.String())
		res, _ := valWrap.AsString()
		assert.Equal(t, "abc\rde", res)
	}
}

func TestAsList(t *testing.T) {
	var valList = []*nebula.XValue_{
		{StringVal: []byte("elem1")},
		{StringVal: []byte("elem2")},
		{StringVal: []byte("elem3")},
	}
	value := nebula.XValue_{
		ListVal: &nebula.XNList_{Values: valList},
	}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, "LIST [\"elem1\", \"elem2\", \"elem3\"]", valWrap.String())
	assert.Equal(t, true, valWrap.IsList())

	res, _ := valWrap.AsList()
	for i := 0; i < len(res); i++ {
		strTemp, err := res[i].AsString()
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, string(valList[i].GetStringVal()), strTemp)
	}
}

func TestAsMap(t *testing.T) {
	valueMap := make(map[string]*nebula.XValue_)
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("key%d", i)
		val := fmt.Sprintf("val%d", i)
		valueMap[key] = &nebula.XValue_{StringVal: []byte(val)}
	}
	mval := nebula.XNRecord_{Values: valueMap}
	value := nebula.XValue_{RecordVal: &mval}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, "RECORD {key0: \"val0\", key1: \"val1\", key2: \"val2\"}", valWrap.String())
	assert.Equal(t, true, valWrap.IsMap())
	vMap := value.GetRecordVal().Values
	valWrapMap, err := valWrap.AsMap()
	if err != nil {
		t.Error(err.Error())
	}
	for i := 0; i < len(vMap); i++ {
		key := fmt.Sprintf("key%d", i)
		str, _ := valWrapMap[key].AsString()
		assert.Equal(t, string(vMap[key].GetStringVal()), str)
	}
}

// func TestAsDuration(t *testing.T) {
// 	value := nebula.Value{DurationVal: &nebula.Duration{1, 2, 3}}
// 	valWrap := ValueWrapper{&value, testTimezone}
// 	assert.Equal(t, true, valWrap.value.IsSetDurationVal())
// 	assert.Equal(t, "1 2 3", valWrap.String())
// }

func TestAsDate(t *testing.T) {
	value := nebula.XValue_{DateVal: &nebula.XDate_{2020, 12, 25}}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, true, valWrap.IsDate())
	assert.Equal(t, "DATE \"2020-12-25\"", valWrap.String())
}

func TestAsLocalTime(t *testing.T) {
	value := nebula.XValue_{LocalTimeVal: &nebula.XLocalTime_{13, 12, 25, 29}}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, true, valWrap.IsLocalTime())
	assert.Equal(t, "13:12:25.000029", valWrap.String())
}

func TestAsLocalDatetime(t *testing.T) {
	value := nebula.XValue_{LocalDatetimeVal: &nebula.XLocalDatetime_{2020, 12, 25, 22, 12, 25, 29}}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, true, valWrap.IsLocalDatetime())
	assert.Equal(t, "DATETIME \"2020-12-25T22:12:25.000029\"", valWrap.String())
}

func TestAsNode(t *testing.T) {
	value := nebula.XValue_{NodeVal: &nebula.XNode_{
		1111,
		2222,
		map[string]nebula.Value{
			"key1": {StringVal: []byte("value1")},
			"key2": {StringVal: []byte("value2")},
		},
	}}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, true, valWrap.IsNode())
	node, _ := valWrap.AsNode()
	assert.Equal(t, int64(1111), node.rawNode.NodeID)
	assert.Equal(t, int16(2222), node.rawNode.NodeTypeID)
	assert.Equal(t, 2, len(node.rawNode.Properties))
	assert.Equal(t, "value1", string(node.rawNode.Properties["key1"].GetStringVal()))
	assert.Equal(t, "value2", string(node.rawNode.Properties["key2"].GetStringVal()))
}

func TestAsEdge(t *testing.T) {
	value := nebula.XValue_{EdgeVal: &nebula.XEdge_{
		1111,
		2222,
		3333,
		4444,
		map[string]nebula.Value{
			"key1": {StringVal: []byte("value1")},
			"key2": {StringVal: []byte("value2")},
		},
	}}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, true, valWrap.IsEdge())
	edge, _ := valWrap.AsEdge()

	assert.Equal(t, int64(1111), edge.rawEdge.SrcID)
	assert.Equal(t, int64(2222), edge.rawEdge.DstID)
	assert.Equal(t, int32(3333), edge.rawEdge.EdgeTypeID)
	assert.Equal(t, int64(4444), edge.rawEdge.Rank)
	assert.Equal(t, 2, len(edge.rawEdge.Properties))
	assert.Equal(t, "value1", string(edge.rawEdge.Properties["key1"].GetStringVal()))
	assert.Equal(t, "value2", string(edge.rawEdge.Properties["key2"].GetStringVal()))
}

func TestAsPath(t *testing.T) {
	// node1 -> edge1 -> node2 -> edge2 -> node3
	var valPathElements = []*nebula.XValue_{
		{
			NodeVal: &nebula.XNode_{
				NodeID:     1,
				NodeTypeID: 1,
				Properties: map[string]nebula.Value{
					"key1": {StringVal: []byte("value1")},
					"key2": {StringVal: []byte("value2")},
				},
			},
		},
		{
			EdgeVal: &nebula.XEdge_{
				SrcID:      1,
				DstID:      2,
				EdgeTypeID: 3333,
				Rank:       4444,
				Properties: map[string]nebula.Value{},
			},
		},
		{
			NodeVal: &nebula.XNode_{
				NodeID:     2,
				NodeTypeID: 1,
				Properties: map[string]nebula.Value{
					"key1": {StringVal: []byte("value1")},
					"key2": {StringVal: []byte("value2")},
				},
			},
		},
	}

	value := nebula.XValue_{
		PathVal: &nebula.XPath_{valPathElements},
	}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, true, valWrap.IsPath())
	path, _ := valWrap.AsPath()
	assert.Equal(t, 3, len(path))
}

func TestResultSet(t *testing.T) {
	respWithNil := &nebula.XExecutionResponse_{
		&nebula.XExecutionOutcome_{
			&nebula.XGQLStatus_{
				[]byte("ERROR"),
			}, nil, nil,
		}, 1000,
	}

	resultSetWithNil, err := genResultSet(respWithNil, testTimezone)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(1000), resultSetWithNil.GetLatency())
	assert.Equal(t, "ERROR", resultSetWithNil.GetStatus())
	assert.Equal(t, false, resultSetWithNil.IsSucceed())

	// Fill a binding table
	var int64 = int64(100)
	var str = []byte("test_string")
	respWithData := &nebula.XExecutionResponse_{
		&nebula.XExecutionOutcome_{
			&nebula.XGQLStatus_{
				[]byte("SUCCESS"),
			}, nil, nil,
		}, 1000,
	}

	respWithData.ExecutionOutcome.Result_ = &nebula.XResultTable_{
		ColumnNames: [][]byte{[]byte("col1"), []byte("col2")},
		Records: []nebula.Row{
			{[]*nebula.XValue_{nebula.NewValue().SetInt64Val(&int64),
				nebula.NewValue().SetStringVal(str)}},
		},
	}
	resultSet, _ := genResultSet(respWithData, testTimezone)

	assert.Equal(t, "SUCCESS", resultSet.GetStatus())
	assert.True(t, resultSet.IsSucceed())
	assert.Equal(t, 2, len(resultSet.GetColNames()))

	expectedTableStr := [][]string([][]string{[]string{"col1", "col2"}, []string{"100", "\"test_string\""}})
	assert.Equal(t, expectedTableStr, resultSet.AsStringTable())
}

func TestString(t *testing.T) {
	testTimezone := timezoneInfo{
		0,
		[]byte("UTC"),
	}

	// Node
	{
		rawNode := &nebula.XNode_{
			1111,
			2222,
			map[string]nebula.Value{
				"key1": {StringVal: []byte("value1")},
				"key2": {StringVal: []byte("value2")},
			},
		}

		node, err := genNode(rawNode, testTimezone)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, `({key1:"value1",key2:"value2"})`, node.String())
	}

	// Edge
	{
		rawEdge := &nebula.XEdge_{
			1111,
			2222,
			3333,
			4444,
			map[string]nebula.Value{
				"key1": {StringVal: []byte("value1")},
				"key2": {StringVal: []byte("value2")},
				"key3": {StringVal: []byte("value3")},
			},
		}

		edge, err := genEdge(rawEdge, testTimezone)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, `[{key1:"value1",key2:"value2",key3:"value3"}]`, edge.String())
	}
}
