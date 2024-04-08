package nebula_ng

import (
	"fmt"
	"io"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/graph"
)

type resultSet struct {
	index     int
	latency   uint64
	result    *graph.ResultTable
	summary   *graph.Summary
	extraInfo []*graph.ExtraInfoElement
}

type rowData struct {
	resultSet *resultSet
	values    []Value
}

type extraInfoType map[graph.ExtraInfoKind]Value

func (rs *resultSet) HasNext() bool {
	if rs.result == nil {
		return false
	}
	return rs.index < len(rs.result.Records)
}

func (rs *resultSet) Next() (Row, error) {
	if !rs.HasNext() {
		return nil, io.EOF
	}

	values := rs.result.GetRecords()[rs.index].Values
	row := &rowData{
		resultSet: rs,
		values:    make([]Value, 0, len(values)),
	}

	for _, v := range values {
		row.values = append(row.values, &grpcValue{data: v})
	}

	rs.index++
	return row, nil
}

func (rs *resultSet) RowSize() int {
	if rs.result == nil {
		return 0
	}
	return len(rs.result.GetRecords())
}

func (rs *resultSet) ColumnTypes() []ColumnType {
	return nil
}

func (rs *resultSet) Summary() Summary {
	if rs.summary == nil {
		return nil
	}
	return &summary{summary: rs.summary}
}

func (rs *resultSet) ExtraInfo() ExtraInfo {
	if rs.extraInfo == nil || len(rs.extraInfo) == 0 {
		return nil
	}
	extr := make(extraInfoType, len(rs.extraInfo))
	for _, e := range rs.extraInfo {
		extr[e.GetKind()] = &grpcValue{data: e.GetValue()}
	}
	return extr
}

func (rs *resultSet) Latency() int64 {
	return int64(rs.latency)
}

func (rs *resultSet) Columns() []string {
	if rs.result == nil {
		return nil
	}
	names := rs.result.ColumnNames
	var cols []string
	for _, name := range names {
		cols = append(cols, string(name))
	}
	return cols
}

func (rd *rowData) Values() []Value {
	return rd.values
}

func (rd *rowData) GetValueByName(name string) (Value, error) {
	names := rd.resultSet.Columns()
	var index int = -1
	for i, n := range names {
		if string(n) == name {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, errInternel(fmt.Sprintf("column %s not found", name))
	}
	return rd.values[index], nil
}

func (rd *rowData) GetValueByIndex(index int) (Value, error) {
	if index < 0 || index >= len(rd.values) {
		return nil, errInternel(fmt.Sprintf("index out of range"))
	}
	return rd.values[index], nil
}

func (ed extraInfoType) get(key graph.ExtraInfoKind) (Value, error) {
	v, ok := ed[key]
	if !ok {
		return nil, errInternel(fmt.Sprintf("key %d not found", key))
	}
	return v, nil
}

func (ed extraInfoType) Cursor() (string, error) {
	v, err := ed.get(graph.ExtraInfoKind_kCursor)
	if err != nil {
		return "", err
	}
	s, err := v.AsString()
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func (ed extraInfoType) AffectedNodes() (int, error) {
	v, err := ed.get(graph.ExtraInfoKind_kAffectedNodes)
	if err != nil {
		return 0, err
	}
	s, err := v.AsUInt64()
	if err != nil {
		return 0, err
	}
	return int(s), nil
}

func (ed extraInfoType) AffectedForwardEdges() (int, error) {
	v, err := ed.get(graph.ExtraInfoKind_kAffectedForwardEdges)
	if err != nil {
		return 0, err
	}
	s, err := v.AsUInt64()
	if err != nil {
		return 0, err
	}
	return int(s), nil
}
func (ed extraInfoType) AffectedReverseEdges() (int, error) {
	v, err := ed.get(graph.ExtraInfoKind_kAffectedReverseEdges)
	if err != nil {
		return 0, err
	}
	s, err := v.AsUInt64()
	if err != nil {
		return 0, err
	}
	return int(s), nil
}
