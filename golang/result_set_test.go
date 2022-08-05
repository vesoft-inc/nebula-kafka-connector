/*
 *
 * Copyright (c) 2020 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 *
 */

package nebula_ng_go

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula/graph"
)

var testTimezone timezoneInfo = timezoneInfo{0, []byte("UTC")}

// TODO(Aiee) To add more
func TestAsNull(t *testing.T) {
	null := nebula.Value{}
	valWrap := ValueWrapper{&null, testTimezone}
	assert.Equal(t, "__NULL__", valWrap.String())
	assert.Equal(t, null, *valWrap.value)
}

func TestAsBool(t *testing.T) {
	bval := new(bool)
	*bval = true
	value := nebula.Value{BoolVal: bval}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, true, valWrap.IsBool())
	assert.Equal(t, "true", valWrap.String())
	res, _ := valWrap.AsBool()
	assert.Equal(t, value.GetBoolVal(), res)
}

func TestAsInt(t *testing.T) {
	val := new(int64)
	*val = 100
	value := nebula.Value{Int64Val: val}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, "int64", valWrap.GetType())
	assert.Equal(t, "100", valWrap.String())
	res, _ := valWrap.AsInt64()
	assert.Equal(t, value.GetInt64Val(), res)
}

func TestAsFloat(t *testing.T) {
	val := new(float64)
	*val = 100.111
	value := nebula.Value{FloatVal: val}
	valWrap := ValueWrapper{&value, testTimezone}
	val2 := new(float64)
	*val2 = 100.00
	value2 := nebula.Value{FloatVal: val2}
	valWrap2 := ValueWrapper{&value2, testTimezone}
	assert.Equal(t, "100.111", valWrap.String())
	assert.Equal(t, "100.0", valWrap2.String())
	assert.Equal(t, true, valWrap.IsFloat())
	res, _ := valWrap.AsFloat()
	assert.Equal(t, value.GetFloatVal(), res)
}

func TestAsString(t *testing.T) {
	val := "test_string"
	value := nebula.Value{StringVal: []byte(val)}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, true, valWrap.IsString())
	assert.Equal(t, "\"test_string\"", valWrap.String())
	res, _ := valWrap.AsString()
	assert.Equal(t, string(value.GetStringVal()), res)
}

func TestAsList(t *testing.T) {
	var valList = []*nebula.Value{
		{StringVal: []byte("elem1")},
		{StringVal: []byte("elem2")},
		{StringVal: []byte("elem3")},
	}
	value := nebula.Value{
		ListVal: &nebula.NList{Values: valList},
	}
	valWrap := ValueWrapper{&value, testTimezone}
	assert.Equal(t, "[\"elem1\", \"elem2\", \"elem3\"]", valWrap.String())
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

func TestResultSet(t *testing.T) {
	respWithNil := &graph.ExecutionResponse{
		&graph.ExecutionOutcome{
			&graph.GQLStatus{
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
	assert.Equal(t, false, resultSetWithNil.IsSuccess())

	// Fill a binding table
	var int64 = int64(100)
	var str = []byte("test_string")
	respWithData := &graph.ExecutionResponse{
		&graph.ExecutionOutcome{
			&graph.GQLStatus{
				[]byte("SUCCESS"),
			}, nil, nil,
		}, 1000,
	}

	respWithData.ExecutionOutcome.Result_ = &nebula.BindingTable{
		ColumnNames: [][]byte{[]byte("col1"), []byte("col2")},
		Records: []*nebula.RawRecord{
			{[]*nebula.Value{nebula.NewValue().SetInt64Val(&int64),
				nebula.NewValue().SetStringVal(str)}},
		},
	}
	resultSet, _ := genResultSet(respWithData, testTimezone)

	assert.Equal(t, "SUCCESS", resultSet.GetStatus())
	assert.True(t, resultSet.IsSuccess())
	assert.Equal(t, 2, len(resultSet.GetColNames()))

	expectedTableStr := [][]string([][]string{[]string{"col1", "col2"}, []string{"100", "\"test_string\""}})
	assert.Equal(t, expectedTableStr, resultSet.AsStringTable())
}
